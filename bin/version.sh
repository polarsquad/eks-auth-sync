#!/bin/sh
#
# Reads the eks-auth-sync version from the environment.
#
set -eu

VERSION_MATCH_REGEX='^v[0-9]+\.[0-9]+\.[0-9]+$'
VERSION_FILTER_REGEX='s/^v//'

check_version() {
    if echo "${1:-}" | grep -Eq "$VERSION_MATCH_REGEX"; then
        echo "$1" | sed "$VERSION_FILTER_REGEX"
        exit
    fi
}

check_version "${CI_COMMIT_TAG:-}"
check_version "${VERSION_TAG:-}"
git describe --match 'v[0-9]*' --dirty | sed "$VERSION_FILTER_REGEX"