#!/usr/bin/env bash
set -euo pipefail

MISE_GO_BIN="$(mise where go)/bin"

echo "Installing golangci-lint with the current Go toolchain ($(go version))..."
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

echo "Running golangci-lint..."
"${MISE_GO_BIN}/golangci-lint" run -c .golang-ci.yml
