<div align="center">
  <img src="frontend/src/assets/logo/magpie-light.png" alt="Magpie logo" height="150">
  <h1>MAGPIE</h1>
  <p><strong>Multi-user AIO Proxy Manager</strong></p>
</div>

<div align="center">
  <img src="https://img.shields.io/github/license/Kuucheen/magpie.svg" alt="license">
  <img src="https://img.shields.io/github/issues/Kuucheen/magpie.svg" alt="issues">

[//]: # (  <img src="https://img.shields.io/github/stars/Kuucheen/magpie.svg?style=social" alt="stars">)
</div>

---

> [!NOTE]
> Magpie is in active development. Features may shift, but the core promise stays the same: less proxy chaos, more time for your real work.

Magpie takes the grind out of running shared proxy infrastructure. It hunts for fresh HTTP / SOCKS proxies, checks their health, enriches them with geo data, and serves them back to you through a dashboard and rotating proxy listeners powered by Go.

## Why Magpie
- **Always fresh lists** – Scheduled scrapers pull from APIs, text dumps, RSS feeds, and dynamic pages (Rod + headless Chromium).
- **Reliable quality** – Configurable judges, retries, and timeouts keep noisy proxies out of your pool.
- **Team friendly** – Multiple accounts share one brain. Magpie de-duplicates work automatically and tracks who owns what.
- **Instant rotation** – Launch rotating proxy listeners with a couple of clicks; Magpie picks free ports for you.
- **Actionable insights** – Charts, breakdowns, and per-proxy history help you decide what to keep or drop.

## Feature Highlights
- **Scraping & Discovery**: Build personal or global scrape lists; Magpie queues them in Redis so nothing gets double-checked across instances.
- **Health Checks**: Smart worker pool in Go keeps throughput high without melting your network.
- **Geo & Reputation**: Optional MaxMind GeoLite2 databases label proxies with country and ISP type.
- **Export & Sharing**: Filter, search, and export directly from the UI or tap into the REST / GraphQL endpoints.
- **Security**: Proxy credentials stay encrypted at rest (bring your own `PROXY_ENCRYPTION_KEY`). JWT auth, admin roles, and user-specific defaults included.

## Quick Start

1. **Install Prerequisites:**
    - [Docker Desktop](https://www.docker.com/)
    - [Git](https://git-scm.com/downloads)

2. **Clone the project**
   ```bash
   git clone https://github.com/Kuucheen/magpie.git
   cd magpie
   ```
3. **Change secrets** – Change `backend/.env` and set a proxy encryption key:
   ```env
   PROXY_ENCRYPTION_KEY=<my-secure-random-key>
   ```
   (Keep that key safe; it secures passwords, proxy auth, and other secrets.)
4. **Bring everything up**
   ```bash
   docker compose up -d --build
   ```
5. **Dive in**
    - UI: http://localhost:8080
    - API: http://localhost:8082/api  
      Register the first account to become the admin.

For geo lookups, create a [MaxMind GeoLite2 account](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) and generate a License Key. Enter it in the dashboard (Admin → Other) to enable automatic database downloads and updates.

## Local Development
- **Services**: `docker compose up -d postgres redis`
- **Backend**: `cd backend && go run ./cmd/magpie`
- **Frontend**: `cd frontend && npm install && npm run start`

Magpie uses Go 1.24.x, Angular 20, PostgreSQL for storage, and Redis for all queueing magic.

## Community
- Chat with us on Discord: https://discord.gg/7FWAGXzhkC
- Bug reports & feature requests: open an issue on GitHub.

## License
Magpie ships under the **GNU Affero General Public License v3.0**. See `LICENSE` for the full text. Contributions are more than welcome.
