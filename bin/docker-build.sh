#!/bin/sh
set -eu

# Hack to get build to work when Docker is not available but Podman is.
# In those cases, Podman can be used as a drop-in replacement for Docker.
DOCKER_CMD="docker"
if command -v podman >/dev/null && ! docker version >/dev/null; then
    DOCKER_CMD="podman"
fi

if [ -z "${1:-}" ]; then
    echo "No image tag specified!" >&2
    exit 1
fi

"$DOCKER_CMD" build \
    -t "$1" \
    --build-arg APP_VERSION="$(bin/version.sh)" \
    --build-arg GIT_HASH="$(bin/git-hash.sh)" \
    .