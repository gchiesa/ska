#!/usr/bin/env bash
set -euo pipefail

VERSION="${VERSION:-$(git describe --abbrev=0 --tags 2>/dev/null || echo "devel")}"

go run -ldflags "-X main.version=${VERSION}" main.go "$@"
