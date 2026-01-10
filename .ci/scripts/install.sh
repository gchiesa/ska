#!/usr/bin/env bash
set -euo pipefail

VERSION="${VERSION:-$(git describe --abbrev=0 --tags 2>/dev/null || echo "devel")}"

echo "Installing ska version ${VERSION}..."
go install -ldflags "-X main.version=${VERSION}"
echo "Installation complete"
