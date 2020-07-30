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

docker run --rm -v "$(pwd):/app:z" -w /app golangci/golangci-lint:v1.27.0 golangci-lint run -v