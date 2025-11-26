#!/usr/bin/env bash
set -euo pipefail

IMAGE_FRONTEND="kuuchen/magpie-frontend"
IMAGE_BACKEND="kuuchen/magpie-backend"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

if ! command -v docker >/dev/null 2>&1; then
  echo "Docker is required but not found in PATH." >&2
  exit 1
fi

tag="${1:-}"
if [[ -z "${tag}" ]]; then
  if git -C "${REPO_ROOT}" rev-parse --short HEAD >/dev/null 2>&1; then
    tag="$(git -C "${REPO_ROOT}" rev-parse --short HEAD)"
  else
    tag="$(date -u +%Y%m%d%H%M%S)"
  fi
fi

build_time="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
push_latest="${PUSH_LATEST:-1}"

echo "Using tag: ${tag}"
echo "Building backend image ${IMAGE_BACKEND}:${tag}"
docker build -f "${REPO_ROOT}/Dockerfile" \
  --build-arg BUILD_VERSION="${tag}" \
  --build-arg BUILD_TIME="${build_time}" \
  -t "${IMAGE_BACKEND}:${tag}" \
  "${REPO_ROOT}"

echo "Building frontend image ${IMAGE_FRONTEND}:${tag}"
docker build -f "${REPO_ROOT}/frontend/Dockerfile" \
  --build-arg BUILD_COMMIT="${tag}" \
  -t "${IMAGE_FRONTEND}:${tag}" \
  "${REPO_ROOT}/frontend"

echo "Pushing ${IMAGE_BACKEND}:${tag}"
docker push "${IMAGE_BACKEND}:${tag}"

echo "Pushing ${IMAGE_FRONTEND}:${tag}"
docker push "${IMAGE_FRONTEND}:${tag}"

if [[ "${push_latest}" == "1" ]]; then
  for image in "${IMAGE_BACKEND}" "${IMAGE_FRONTEND}"; do
    echo "Tagging and pushing ${image}:latest"
    docker tag "${image}:${tag}" "${image}:latest"
    docker push "${image}:latest"
  done
fi

echo "Done. Pushed images with tag ${tag}${push_latest:+ and latest}."
