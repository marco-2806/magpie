###############  Stage 1 – Angular build  ###############
FROM node:20-alpine AS frontend-build

WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm ci --silent         # faster + reproducible
COPY frontend/ .
RUN npm run build --prod    # creates dist/<app-name>

###############  Stage 2 – Go build (+ embed) ############
FROM golang:1.24-alpine AS backend-build

WORKDIR /backend
COPY --from=frontend-build /frontend/dist ./web/dist
COPY backend/ .

ENV CGO_ENABLED=0
RUN go build -o server .   # point to the ONE main package

############ Stage 3 – grab Chromium’s shared libraries ############
FROM debian:bookworm-slim AS chromium-deps

RUN apt-get update && apt-get install -y --no-install-recommends \
      libglib2.0-0  libgtk-3-0  libnss3  libasound2 \
      libatk-bridge2.0-0  libatk1.0-0  libcups2  libdrm2  libgbm1 \
      libx11-xcb1  libxcomposite1  libxdamage1  libxrandr2  libxkbcommon0 \
      fonts-liberation  ca-certificates  xdg-utils \
  && rm -rf /var/lib/apt/lists/*

############ Stage 4 – tiny runtime image ##########################
FROM gcr.io/distroless/base-debian12

# copy just the libs/fonts we installed above
COPY --from=chromium-deps /lib/x86_64-linux-gnu /lib/x86_64-linux-gnu
COPY --from=chromium-deps /usr/lib/x86_64-linux-gnu /usr/lib/x86_64-linux-gnu
COPY --from=chromium-deps /usr/share/fonts /usr/share/fonts

WORKDIR /app
COPY --from=backend-build /backend/server .
COPY --from=backend-build /backend/web/dist ./static

EXPOSE 8082
ENTRYPOINT ["./server"]
