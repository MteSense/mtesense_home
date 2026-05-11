#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RELEASE_DIR="${ROOT_DIR}/release/mtesense-home"

cd "${ROOT_DIR}/web/app"
npm ci
npm run build

cd "${ROOT_DIR}"
mkdir -p "${RELEASE_DIR}/web/app" "${RELEASE_DIR}/data" "${RELEASE_DIR}/public_uploads"

GOOS=linux GOARCH="${GOARCH:-amd64}" CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "${RELEASE_DIR}/mtesense-home" ./cmd/server

cp -R web/app/dist "${RELEASE_DIR}/web/app/"
cp .env.example "${RELEASE_DIR}/.env.example"
cp README.md "${RELEASE_DIR}/README.md"

echo "Linux release written to ${RELEASE_DIR}"
