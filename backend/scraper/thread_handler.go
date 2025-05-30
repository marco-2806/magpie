package scraper

import (
	"fmt"
	redis_queue2 "magpie/checker/redis_queue"
	"magpie/database"
	"magpie/helper"
	"magpie/models"
	"math"
	"strings"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
	"magpie/scraper/redis_queue"
	"magpie/settings"
)

var (
	/* ─────────────────────────────  thread control  ─────────────────────────── */
	currentThreads atomic.Uint32
	stopThread     = make(chan struct{}) // signals a worker to exit

	/* ─────────────────────────────  page pool  ─────────────────────────────── */
	browser      *rod.Browser
	pagePool     chan *rod.Page
	currentPages atomic.Int32
	stopPage     = make(chan struct{}) // signals that a page should be closed
)

/*──────────────────────────────────────────────────────────────────────────────*/

func init() {
	browser = rod.New().MustConnect()
	pagePool = make(chan *rod.Page, 2000) // hard cap – raise if you really need more
}

/*─────────────────────────────  dynamic dispatcher  ──────────────────────────*/

func ThreadDispatcher() {
	for {
		cfg := settings.GetConfig()

		var target uint32
		if cfg.Scraper.DynamicThreads {
			target = autoThreadCount(cfg)
		} else {
			target = cfg.Scraper.Threads
		}

		/* spawn */
		for currentThreads.Load() < target {
			go scrapeWorker()
			currentThreads.Add(1)
		}

		/* retire */
		for currentThreads.Load() > target {
			stopThread <- struct{}{}
			currentThreads.Add(^uint32(0)) // decrement
		}

		log.Debug("scraper threads", "active", currentThreads.Load())
		time.Sleep(15 * time.Second)
	}
}

/*──────────────────────────────  worker goroutine  ───────────────────────────*/

func scrapeWorker() {
	for {
		select {
		case <-stopThread:
			return
		default:
		}

		site, due, err := redis_queue.PublicScrapeSiteQueue.GetNextScrapeSite()
		if err != nil {
			log.Error("pop scrape site", "err", err)
			time.Sleep(2 * time.Second)
			continue
		}

		cfg := settings.GetConfig()
		timeout := time.Duration(cfg.Scraper.Timeout) * time.Millisecond

		var html string
		var scrapeErr error

		for attempts := 0; attempts < 3; attempts++ {
			html, scrapeErr = ScraperRequest(site.URL, timeout)
			if scrapeErr == nil || !strings.Contains(scrapeErr.Error(), "timeout waiting for available page") {
				break
			}
			log.Debug("retrying after page timeout", "url", site.URL, "attempt", attempts+1)
			time.Sleep(1 * time.Second)
		}

		if scrapeErr != nil {
			log.Warn("scrape failed", "url", site.URL, "err", scrapeErr)
		} else {
			go handleScrapedHTML(site, html)
		}

		if err := redis_queue.PublicScrapeSiteQueue.RequeueScrapeSite(site, due); err != nil {
			log.Error("requeue site", "err", err)
		}
	}
}

/*───────────────────────────────  auto‑sizing  ───────────────────────────────*/

func autoThreadCount(cfg settings.Config) uint32 {
	totalSites, err := redis_queue.PublicScrapeSiteQueue.GetScrapeSiteCount()
	if err != nil {
		log.Error("count sites", "err", err)
		return 1
	}

	instances, err := redis_queue.PublicScrapeSiteQueue.GetActiveInstances()
	if err != nil || instances == 0 {
		instances = 1
	}

	perInstance := (totalSites + int64(instances) - 1) / int64(instances)

	period := settings.CalculateMillisecondsOfCheckingPeriod(cfg.Scraper.ScraperTimer)
	if period == 0 {
		log.Warn("scraper period 0 → forcing 1 day")
		period = 86_400_000
	}

	numerator := uint64(perInstance) * uint64(cfg.Scraper.Timeout) * uint64(cfg.Scraper.Retries+1)
	threads := (numerator + period - 1) / period

	switch {
	case threads == 0 && perInstance > 0:
		threads = 1
	case threads > math.MaxUint32:
		threads = math.MaxUint32
	}
	return uint32(threads)
}

/*──────────────────────────────  page‑pool logic  ────────────────────────────*/

func ManagePagePool() {
	for {
		cfg := settings.GetConfig()
		targetPages := calcRequiredPages(cfg)

		/* add pages */
		for currentPages.Load() < targetPages {
			if err := addPage(); err != nil {
				log.Error("add page", "err", err)
				time.Sleep(1 * time.Second)
				continue
			}
		}

		/* shed pages */
		for currentPages.Load() > targetPages {
			select {
			case p := <-pagePool:
				p.MustClose()
				currentPages.Add(-1)
			default:
				stopPage <- struct{}{}
			}
		}

		time.Sleep(15 * time.Second)
	}
}

func calcRequiredPages(cfg settings.Config) int32 {
	count := uint64(1) // default

	if n, err := redis_queue.PublicScrapeSiteQueue.GetScrapeSiteCount(); err == nil {
		count = uint64(n)
	}

	interval := settings.CalculateMillisecondsOfCheckingPeriod(cfg.Scraper.ScraperTimer)
	if interval == 0 {
		interval = 86_400_000
	}
	avg := uint64(cfg.Scraper.Timeout * (cfg.Scraper.Retries + 1)) // ms

	required := (count * avg) / uint64(interval)
	if required < 1 && count > 0 {
		required = 1
	}
	if required > 2000 {
		required = 2000
	}
	return int32(required)
}

func addPage() error {
	p, err := stealth.Page(browser)
	if err != nil {
		return fmt.Errorf("stealth page: %w", err)
	}
	select {
	case pagePool <- p:
		currentPages.Add(1)
	default:
		p.MustClose()
		return fmt.Errorf("pool full")
	}
	return nil
}

/*───────────────────────────────  helpers  ───────────────────────────────────*/

func recyclePage(p *rod.Page) {
	select {
	case <-stopPage:
		p.MustClose()
		currentPages.Add(-1)
	default:
		if err := resetPage(p); err != nil {
			log.Debug("page reset failed, replacing", "err", err)
			p.MustClose()
			currentPages.Add(-1)
			// Always add a replacement when a page is closed
			go func() {
				if err := addPage(); err != nil {
					log.Error("add replacement page", "err", err)
				}
			}()
		} else {
			select {
			case pagePool <- p:
				// Successfully recycled
			default:
				// Pool is full, close the page
				p.MustClose()
				currentPages.Add(-1)
			}
		}
	}
}

func handleScrapedHTML(site models.ScrapeSite, rawHTML string) {
	proxyList := helper.GetProxiesOfHTML(rawHTML)

	parsedProxies := helper.ParseTextToProxies(strings.Join(proxyList, "\n"))

	proxies, err := database.InsertAndGetProxies(parsedProxies, helper.GetUserIdsFromList(site.Users)...)
	if err != nil {
		log.Error("insert proxies from scraping failed", "err", err)
	}

	err = database.AssociateProxiesToScrapeSite(site.ID, proxies)
	if err != nil {
		log.Warn("associate proxies to ScrapeSite failed", "err", err)
	}

	err = redis_queue2.PublicProxyQueue.AddToQueue(proxies)
	if err != nil {
		log.Error("adding scraped proxies to queue failed", "err", err)
	}

	log.Info(fmt.Sprintf("Found %d unique proxies that users don't have", len(proxies)), "url", site.URL)
}
