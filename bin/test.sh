#!/bin/sh
#
# Run tests and collect a coverage report.
#
set -eu

go test \
    -v \
    -race \
    -coverpkg ./internal/... \
    -coverprofile=coverage.txt \
    -covermode=atomic \
    ./...
