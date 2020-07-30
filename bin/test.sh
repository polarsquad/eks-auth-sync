#!/bin/sh
set -eu

go test \
    -race \
    -coverpkg ./internal/... \
    -coverprofile=coverage.txt \
    -covermode=atomic \
    ./...