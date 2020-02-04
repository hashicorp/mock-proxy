#! /usr/bin/env bash
set -euo pipefail

# Always work from the root of the repo.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR"/.. && pwd)"
cd "$ROOT_DIR"

docker-compose --file deployments/docker-compose.yml build
docker-compose --file deployments/docker-compose.yml run client \
  /bin/sh -c '
    dockerize \
      -wait tcp://squid.proxy:8888 \
      -timeout 60s \
      -wait-retry-interval 5s \
    && /bin/bash
  '
