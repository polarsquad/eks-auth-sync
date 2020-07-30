#!/usr/bin/env bash
set -euo pipefail
cd "${0%/*}/.."

# Hack to get build to work when Docker is not available but Podman is.
# In those cases, Podman can be used as a drop-in replacement for Docker.
if command -v podman >/dev/null && ! docker version >/dev/null; then
    docker() {
        podman "$@"
    }
    export -f docker
fi

if [ -z "${1:-}" ]; then
    echo "No image tag specified!" >&2
    exit 1
fi

docker build \
    -t "$1" \
    --build-arg APP_VERSION="$(bin/version.sh)" \
    --build-arg GIT_HASH="$(bin/git-hash.sh)" \
    .