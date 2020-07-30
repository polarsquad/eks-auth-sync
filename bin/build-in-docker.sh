#!/bin/sh
#
# Builds the binary for eks-auth-sync in Docker and the Docker image for it.
# Parameters:
#   $1 = Image tag for eks-auth-sync Docker image (default: eks-auth-sync)
#
set -eu

# Image tag for the image to build
IMAGE_TAG="${1:-"eks-auth-sync"}"

# Hack to get build to work when Docker is not available but Podman is.
# In those cases, Podman can be used as a drop-in replacement for Docker.
DOCKER_CMD="docker"
if command -v podman >/dev/null && ! docker version >/dev/null; then
    DOCKER_CMD="podman"
fi

# Cache for Go mod dependencies
GOMODCACHEPATH="$(pwd)/.gomodcache"
if [ ! -d "${GOMODCACHEPATH}" ]; then
    mkdir -p "${GOMODCACHEPATH}"
fi

"$DOCKER_CMD" run \
    --rm \
    -v "${GOMODCACHEPATH}:/go/pkg/mod:z" \
    -v "$(pwd):/project:z" \
    -w /project \
    -e APP_VERSION="$(bin/version.sh)" \
    -e GIT_HASH="$(bin/git-hash.sh)" \
    golang:1.14 \
    ./bin/build.sh

"$DOCKER_CMD" build -t "${IMAGE_TAG}" .
