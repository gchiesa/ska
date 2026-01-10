#!/usr/bin/env bash
set -euo pipefail

echo "Formatting Go files..."
gofumpt -w .
gci write .
echo "Formatting complete"
