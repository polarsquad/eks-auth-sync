#!/usr/bin/env bash
set -euo pipefail
cd "${0%/*}/.."

check_git_sha() {
    if [ -n "${1:-}" ]; then
        echo "$1"
        exit
    fi
}

check_git_sha "${CI_COMMIT_SHA:-}"
check_git_sha "${GIT_HASH:-}"
git rev-parse HEAD