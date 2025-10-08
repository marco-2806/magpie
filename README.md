<div align="center">
  <img src="frontend/src/assets/logo/magpie-light.png" style="height: 150px" alt="MAGPIE logo"/>

## M A G P I E
**Multi‑user AIO Proxy Manager**
</div>

---

> [!WARNING]
> This project is still in alpha development. Some features may not work correctly.


Magpie is an **open‑source, Docker‑first** proxy management suite written in **Go (back‑end)** and **Angular (front‑end)**. It continuously scrapes, de‑duplicates, and health‑checks HTTP/SOCKS proxies so you don’t have to.

---

## ✨ Feature Highlights

| Category | Details                                                                                               |
|----------|-------------------------------------------------------------------------------------------------------|
| **Automatic Scraping** | Pulls from open‑web APIs, plaintext lists, RSS feeds, or your own endpoints on an adjustable schedule. |
| **Smart Health Checks** | Concurrency‑limited pingers with configurable `timeout`, `retries`, and `interval`.                   |
| **De‑duplication** | A proxy is checked **once**, even if multiple users request it—saving bandwidth & IP reputation.      |
| **Thread Pool** | Can dynamically scale Go workers (threads) based on the proxy count and settings.                     |
| **Tagging & Groups** | Organize proxies by geo, anonymity, speed, or if the proxy is alive.                                  |

[//]: # (| **Live Dashboard** | Angular UI with filtering, charts, and real‑time WebSocket updates.                                   |)

---

## 🚀 Quick Start (Docker‑Compose)

1. **Install Prerequisites:**
    - [Docker Desktop](https://www.docker.com/)
    - [Git](https://git-scm.com/downloads)

2. **(OPTIONAL) Set Up GeoLite2 Database:**

   If you want to determine the country and type (ISP, Datacenter, or Residential) of the proxies, you'll need to download the [GeoLite2 Country and GeoLite2 ASN Database](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) from MaxMind.

    - After downloading, replace the existing databases in the `backend/database` directory.
    - Ensure they have the same names: `GeoLite2-ASN.mmdb` and `GeoLite2-Country.mmdb`.<br><br>

   > **Note**: These databases are not included in the repository due to their licensing restrictions.

3. **Clone and start the Project**:

   Open your terminal and enter the following commands
   ```bash
   $ git https://github.com/Kuucheen/magpie.git
   ```
   
   ```bash
   $ cd magpie
   ```
   Now start the program with this command
   ```bash
   $ docker compose up -d --build
    ```

   After everything started up (this can take a few seconds) you can open your browser and enter the following URL:
   http://localhost:8082
   
   Now register an account with an email (does not need to be a real email) and a password and your good to go.

### Required Secrets

Set the following environment variables before starting the back-end service:

- `PROXY_ENCRYPTION_KEY` –  A 32-character (or longer) secret used to encrypt proxy credentials stored in the database. Changing this key after proxies have been saved makes them unreadable, so keep it safe and stable across deployments.

---

For early support, join our <a href="https://discord.gg/7FWAGXzhkC">discord server</a>


[//]: # (## ⚙️ Configuration)

[//]: # ()
[//]: # (| Variable | Default | Description |)

[//]: # (|----------|---------|-------------|)

[//]: # (| `MAGPIE_DB_DSN` | `postgres://magpie:magpie@db:5432/magpie` | PostgreSQL DSN |)

[//]: # (| `MAGPIE_API_PORT` | `8080` | HTTP port exposed by the Go service |)

[//]: # (| `MAGPIE_SCRAPE_INTERVAL` | `15m` | How often to trigger the global scraper loop |)

[//]: # (| `MAGPIE_CHECK_TIMEOUT` | `5s` | Per‑proxy health‑check timeout |)

[//]: # (| `MAGPIE_CHECK_RETRIES` | `2` | Retries before marking a proxy “dead” |)

[//]: # (| `MAGPIE_MAX_WORKERS` | `250` | Hard cap for concurrent workers |)

[//]: # (| `MAGPIE_JWT_SECRET` | `change‑me` | Auth token signing key |)

[//]: # (| `MAGPIE_ADMIN_EMAIL` | `admin@example.com` | First admin user &#40;auto‑created&#41; |)

[//]: # ()
[//]: # (Put these in a `.env` or pass `-e KEY=value` to `docker compose`.)

[//]: # (---)

[//]: # ()
[//]: # (## 🖥️ Using Magpie)

[//]: # ()
[//]: # (### Add Proxies via UI)

[//]: # (1. Navigate to **Proxies → Import**.)

[//]: # (2. Paste raw list or upload a `.csv` file &#40;format: `ip,port[,username,password]`&#41;.)

[//]: # (3. Click **Import** and watch them validate in real time.)


[//]: # (## 🛠️ Development)

[//]: # ()
[//]: # (| Service | Command |)

[//]: # (|---------|---------|)

[//]: # (| **Back‑end** | `go run ./cmd/server` &#40;auto‑reload via `air`&#41; |)

[//]: # (| **Front‑end** | `npm i && ng serve --open` |)

[//]: # (| **Lint / Tests** | `make lint test` |)

[//]: # (---)

[//]: # ()
[//]: # (## ♻️ Updating)

[//]: # ()
[//]: # (```bash)

[//]: # ($ git pull)

[//]: # ($ docker compose pull && docker compose up -d --build)

[//]: # (```)

[//]: # (*&#40;Zero‑downtime migrations are applied automatically.&#41;*)

[//]: #---

[//]: # (## 📜 License)

[//]: # ()
[//]: # (Magpie is released under the **MIT License**—see [`LICENSE`]&#40;LICENSE&#41; for details.)

[//]: # ()
[//]: # (---)

[//]: # ()
[//]: # (## 🙏 Contributing)

[//]: # ()
[//]: # (Pull requests are welcome! Please open an issue first to discuss major changes. Make sure to run `make test` and abide by the [code of conduct]&#40;CODE_OF_CONDUCT.md&#41;.)

[//]: # ()
