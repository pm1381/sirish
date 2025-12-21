#!/usr/bin/env bash

set -euo pipefail

echo "▶ Running all Go tests (verbose, race)..."

go test -v -race ./...

echo "✅ All tests passed"