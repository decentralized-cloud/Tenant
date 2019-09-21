#!/usr/bin/env bash
#
# bootstrap.sh check for and install any dependencies we have for building and using Tenant service
#

set -e

# making sure we change the directory to where the current script resides
current_directory=$(dirname "$0")
cd "$current_directory"

# importing versions
. ./versions.sh

# Pulling down the required docker images
docker pull decentralizedcloud/golang-build-base:$DOCKER_GOLANG_BUILD_BASE_IMAGE
docker pull decentralizedcloud/golang-test-and-coverage-base:$DOCKER_GOLANG_TEST_AND_COVERAGE_BASE_IMAGE

