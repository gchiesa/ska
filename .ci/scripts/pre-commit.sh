#!/usr/bin/env bash
set -euo pipefail

echo "Running pre-commit hooks..."
pre-commit run --all-files
