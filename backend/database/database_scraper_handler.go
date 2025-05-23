package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"magpie/helper"
	"magpie/models/routeModels"

	"magpie/models"
)

const scrapeSitesPerPage = 20

// GetScrapingSourcesOfUsers returns all URLs associated with the given user.
func GetScrapingSourcesOfUsers(userID uint) []string {
	var user models.User
	if err := DB.Preload("ScrapeSites").First(&user, userID).Error; err != nil {
		return nil
	}

	out := make([]string, 0, len(user.ScrapeSites))
	for _, s := range user.ScrapeSites {
		out = append(out, s.URL)
	}
	return out
}

// SaveScrapingSourcesOfUsers replaces the user’s current list with `sources`.
func SaveScrapingSourcesOfUsers(userID uint, sources []string) ([]models.ScrapeSite, error) {
	var sites []models.ScrapeSite
	err := DB.Transaction(func(tx *gorm.DB) error {
		// Prepare slices to collect sites and their IDs
		sites = make([]models.ScrapeSite, 0, len(sources))
		siteIDs := make([]uint64, 0, len(sources))

		// Create or find each ScrapeSite
		for _, raw := range sources {
			if raw == "" || !helper.IsValidURL(raw) {
				continue
			}

			var site models.ScrapeSite
			if err := tx.Where("url = ?", raw).
				FirstOrCreate(&site, &models.ScrapeSite{URL: raw}).Error; err != nil {
				return err
			}
			sites = append(sites, site)
			siteIDs = append(siteIDs, site.ID)
		}

		// Load the user
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		// Replace association to the new list of ScrapeSites
		if err := tx.Model(&user).
			Association("ScrapeSites").
			Replace(&sites); err != nil {
			return err
		}

		// Reload all sites with Users preloaded
		var loaded []models.ScrapeSite
		if err := tx.Preload("Users").
			Where("id IN ?", siteIDs).
			Find(&loaded).Error; err != nil {
			return err
		}
		// Overwrite the sites slice with the fully loaded records
		sites = loaded

		return nil
	})

	return sites, err
}

func GetAllScrapeSites() ([]models.ScrapeSite, error) {
	var allProxies []models.ScrapeSite
	const batchSize = maxParamsPerBatch

	collectedProxies := make([]models.ScrapeSite, 0)

	err := DB.Preload("Users").Order("id").FindInBatches(&allProxies, batchSize, func(tx *gorm.DB, batch int) error {
		collectedProxies = append(collectedProxies, allProxies...)
		return nil
	})

	if err.Error != nil {
		return nil, err.Error
	}

	return collectedProxies, nil
}

func AssociateProxiesToScrapeSite(siteID uint64, proxies []models.Proxy) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if len(proxies) == 0 {
			return nil
		}

		assoc := make([]models.ProxyScrapeSite, len(proxies))
		for i, p := range proxies {
			assoc[i] = models.ProxyScrapeSite{
				ProxyID:      p.ID,
				ScrapeSiteID: siteID,
			}
		}

		return tx.
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "proxy_id"}, {Name: "scrape_site_id"}},
				DoNothing: true,
			}).
			Create(&assoc).Error
	})
}

func GetAllScrapeSiteCountOfUser(userId uint) int64 {
	var count int64
	DB.Model(&models.ScrapeSite{}).
		Joins(
			"JOIN user_scrape_site uss ON uss.scrape_site_id = scrape_sites.id AND uss.user_id = ?",
			userId,
		).
		Count(&count)
	return count
}

func GetScrapeSiteInfoPage(userId uint, page int) []routeModels.ScrapeSiteInfo {
	offset := (page - 1) * scrapeSitesPerPage

	// subquery: for each scrape_site_id, count only the proxies that this user has
	subQuery := DB.
		Model(&models.ProxyScrapeSite{}).
		Select("scrape_site_id, COUNT(*) AS proxy_count").
		Joins("JOIN user_proxies up ON up.proxy_id = proxy_scrape_site.proxy_id AND up.user_id = ?", userId).
		Group("scrape_site_id")

	var results []routeModels.ScrapeSiteInfo

	DB.
		Model(&models.ScrapeSite{}).
		Select(
			"scrape_sites.id         AS id, "+
				"scrape_sites.url        AS url, "+
				"COALESCE(pc.proxy_count, 0) AS proxy_count, "+
				"uss.created_at          AS added_at",
		).
		// only the sites this user has added
		Joins("JOIN user_scrape_site uss ON uss.scrape_site_id = scrape_sites.id AND uss.user_id = ?", userId).
		// attach the per-site, per-user proxy counts
		Joins("LEFT JOIN (?) AS pc ON pc.scrape_site_id = scrape_sites.id", subQuery).
		Order("uss.created_at DESC").
		Offset(offset).
		Limit(scrapeSitesPerPage).
		Scan(&results)

	return results
}

func DeleteScrapeSiteRelation(userId uint, scrapeSite []int) {
	DB.Where("scrape_site_id IN (?)", scrapeSite).Where("user_id = (?)", userId).Delete(&models.UserScrapeSite{})
}
