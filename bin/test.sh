#!/usr/bin/env bash
set -euo pipefail
cd "${0%/*}/.."

go test -race -coverpkg ./internal/... -coverprofile=coverage.txt -covermode=atomic ./...
