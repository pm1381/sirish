#!/usr/bin/env bash

set -x
set -euo pipefail

echo "▶ Running all Go tests (verbose)..."

go test -v  ./... -coverprofile=coverage.txt -covermode count || exit 1
go tool cover -func=coverage.txt
go tool cover -html coverage.txt -o coverage.html
go get github.com/boumenot/gocover-cobertura
bash -c "go run github.com/boumenot/gocover-cobertura < coverage.txt > coverage.xml"
