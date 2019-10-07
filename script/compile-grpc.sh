#!/usr/bin/env sh

set -e
set -x

cleanup() {
	docker rm extract-tenant-contract-grpc-builder
}

trap 'cleanup' EXIT

if [ $# -eq 0 ]; then
	current_directory=$(dirname "$0")
else
	current_directory="$1"
fi

cd "$current_directory"/..

docker build -f docker/Dockerfile.buildGrpcContract -t tenant-contract-grpc-builder .
docker create --name extract-tenant-contract-grpc-builder tenant-contract-grpc-builder
docker cp extract-tenant-contract-grpc-builder:/src/contract/grpc/go/tenant.pb.go ./contract/grpc/go

