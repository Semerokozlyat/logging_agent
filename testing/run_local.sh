#!/usr/bin/env bash

set -e
shopt -s expand_aliases

CUR_SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
cd "${CUR_SCRIPT_DIR}"

alias docker_compose='docker compose -f docker-compose-local.yml'

function docker_cleanup() {
  echo "Cleaning up Docker containers..."
  docker_compose down --remove-orphans --volumes
}

trap docker_cleanup EXIT

echo "Starting local services..."
docker_cleanup  # Clean any previous state
docker_compose up --remove-orphans