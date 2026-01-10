#!/usr/bin/env bash
set -euo pipefail

COVERAGE_FILE="${COVERAGE_FILE:-coverage.out}"

echo "Running tests with coverage..."
go test -v -race -coverprofile="${COVERAGE_FILE}" -covermode=atomic ./...

echo ""
echo "=== Coverage Report ==="
go tool cover -func="${COVERAGE_FILE}"
