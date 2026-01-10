#!/usr/bin/env bash
set -euo pipefail

echo "Installing build dependencies..."
go generate -tags tools tools/tools.go
echo "Bootstrap complete"
