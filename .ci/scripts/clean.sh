#!/usr/bin/env bash
set -euo pipefail

PROJECT_NAME="${PROJECT_NAME:-ska}"

echo "Cleaning up..."
rm -rf coverage.out dist/ "${PROJECT_NAME}" test-output.json
echo "Cleanup complete"
