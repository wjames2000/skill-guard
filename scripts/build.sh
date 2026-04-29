#!/bin/bash
set -e
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS="-X 'github.com/wjames2000/skill-guard/cmd.Version=${VERSION}' -X 'github.com/wjames2000/skill-guard/cmd.Commit=${COMMIT}' -X 'github.com/wjames2000/skill-guard/cmd.Date=${DATE}'"
mkdir -p dist
PLATFORMS=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")
for PLATFORM in "${PLATFORMS[@]}"; do
	GOOS=${PLATFORM%/*}
	GOARCH=${PLATFORM#*/}
	OUTPUT="dist/skill-guard-${GOOS}-${GOARCH}"
	[ "$GOOS" = "windows" ] && OUTPUT+=".exe"
	echo "Building ${OUTPUT}..."
	GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -ldflags "${LDFLAGS}" -o "${OUTPUT}" .
done
echo "Done: $(ls dist/ | wc -l) files"
ls -lh dist/
