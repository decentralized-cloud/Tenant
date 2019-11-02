#!/usr/bin/env sh

set -e
set -x

if [ $# -eq 0 ]; then
	current_directory=$(dirname "$0")
else
	current_directory="$1"
fi

cd "$current_directory"/..

docker-compose  -f docker/docker-compose-integration.yml build
docker-compose  -f docker/docker-compose-integration.yml run --rm tenant
docker-compose  -f docker/docker-compose-integration.yml down

