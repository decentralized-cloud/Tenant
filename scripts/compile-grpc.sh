#!/usr/bin/env sh

set -e
set -x

cleanup() {
	docker rm extract-contract-grpc-builder
}

trap 'cleanup' EXIT

if [ $# -eq 0 ]; then
	current_directory=$(dirname "$0")
else
	current_directory="$1"
fi

cd "$current_directory"/..

docker build -f docker/Dockerfile.buildGrpcContract -t contract-grpc-builder .
docker create --name extract-contract-grpc-builder contract-grpc-builder
docker cp extract-contract-grpc-builder:/src/contract/grpc/go ./contract/grpc/
