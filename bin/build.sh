#!/usr/bin/env bash
set -euo pipefail
cd "${0%/*}/.."

export CGO_ENABLED=0
ROOT_PKG="gitlab.com/polarsquad/eks-auth-sync"
BUILDINFO_PKG="$ROOT_PKG/internal/buildinfo"

if [ -z "${APP_VERSION:-}" ]; then
    APP_VERSION=$(bin/version.sh)
fi
if [ -z "${GIT_HASH:-}" ]; then
    GIT_HASH=$(bin/git-hash.sh)
fi

LDFLAGS="-s -w"
LDFLAGS+=" -X '$BUILDINFO_PKG.Version=$APP_VERSION'"
LDFLAGS+=" -X '$BUILDINFO_PKG.GitHash=$GIT_HASH'"

go build \
    -ldflags="$LDFLAGS" \
    -o eks-auth-sync \
    ./cmd/eksauthsync
