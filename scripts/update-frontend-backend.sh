#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

cd "${REPO_ROOT}"

if docker compose version >/dev/null 2>&1; then
  compose_cmd=(docker compose)
elif command -v docker-compose >/dev/null 2>&1; then
  compose_cmd=(docker-compose)
else
  echo "Docker Compose is required but was not found. Install Docker Desktop or docker-compose." >&2
  exit 1
fi

echo "Pulling latest changes..."
git pull --ff-only

echo "Rebuilding frontend and backend containers..."
"${compose_cmd[@]}" up -d --build frontend backend

echo "Done. Frontend is available at http://localhost:8080 and backend API at http://localhost:8082/api"
