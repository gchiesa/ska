#!/usr/bin/env bash
set -euo pipefail

COVERAGE_FILE="${COVERAGE_FILE:-coverage.out}"

echo "Running tests with coverage..."

# Run tests with gotestsum and pipe JSON output to tparse for summary
gotestsum --format=standard-verbose --jsonfile=test-output.json -- \
    -cover \
    -coverprofile="${COVERAGE_FILE}" \
    -covermode=atomic \
    -race \
    ./...

echo ""
echo "=== Test Summary ==="
tparse -file=test-output.json -all

echo ""
echo "=== Coverage Summary ==="
go tool cover -func="${COVERAGE_FILE}" | tail -1

# Cleanup temp file
rm -f test-output.json
