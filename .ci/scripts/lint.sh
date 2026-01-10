#!/usr/bin/env bash
set -euo pipefail

echo "Running golangci-lint..."
golangci-lint run -c .golang-ci.yml
