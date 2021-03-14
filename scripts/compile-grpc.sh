#!/usr/bin/env sh

set -e
set -x

cleanup() {
	docker rm extract-project-contract-grpc-builder
}

trap 'cleanup' EXIT

if [ $# -eq 0 ]; then
	current_directory=$(dirname "$0")
else
	current_directory="$1"
fi

cd "$current_directory"/..

docker build -f docker/Dockerfile.buildGrpcContract -t project-contract-grpc-builder .
docker create --name extract-project-contract-grpc-builder project-contract-grpc-builder
docker cp extract-project-contract-grpc-builder:/src/contract/grpc/go/project.pb.go ./contract/grpc/go

