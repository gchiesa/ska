#!/usr/bin/env bash
set -euo pipefail

PROJECT_NAME="${PROJECT_NAME:-ska}"
VERSION="${VERSION:-$(git describe --abbrev=0 --tags 2>/dev/null || echo "devel")}"

echo "Building ${PROJECT_NAME} version ${VERSION}..."
go build -ldflags "-X cmd.version=${VERSION}" -o "${PROJECT_NAME}"
echo "Build complete: ./${PROJECT_NAME}"
