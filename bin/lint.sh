#!/bin/sh
set -eu

# Hack to get build to work when Docker is not available but Podman is.
# In those cases, Podman can be used as a drop-in replacement for Docker.
DOCKER_CMD="docker"
if command -v podman >/dev/null && ! docker version >/dev/null; then
    DOCKER_CMD="podman"
fi

"$DOCKER_CMD" run \
    --rm \
    -v "$(pwd):/app:z" \
    -w /app \
    golangci/golangci-lint:v1.27.0 \
    golangci-lint run -v --timeout 2m --skip-dirs '(^|/).go($|/)'
