#!/usr/bin/env bash
# Release packaging script for mllt-cli.
# Usage: ./release.sh v1.0.0
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "Usage: $0 <version>" >&2
  exit 1
fi

VERSION="$1"
APP_NAME="mllt-cli"
CMD_PATH="./cmd/mllt-cli"
DIST_ROOT="dist/${VERSION}"
mkdir -p "$DIST_ROOT"

echo "==> Building release artifacts for ${APP_NAME} ${VERSION}"

checksum_tool=""
if command -v sha256sum >/dev/null 2>&1; then
  checksum_tool="sha256sum"
elif command -v shasum >/dev/null 2>&1; then
  checksum_tool="shasum -a 256"
else
  echo "Warning: sha256sum or shasum not found; checksums will not be generated." >&2
fi

TMP_DIR=$(mktemp -d)
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

build_binary() {
  local goos=$1
  local goarch=$2
  local output_path=$3
  local label=$4
  echo "-> Building ${label} (GOOS=${goos}, GOARCH=${goarch})"
  GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 go build -ldflags "-s -w" -o "$output_path" "$CMD_PATH"
}

package_tar() {
  local dir=$1
  local archive=$2
  tar -C "$dir" -czf "$archive" "$APP_NAME"
}

package_zip() {
  local dir=$1
  local archive=$2
  zip -j -q "$archive" "$dir/${APP_NAME}.exe"
}

record_checksum() {
  local file=$1
  if [[ -n "$checksum_tool" ]]; then
    (cd "$DIST_ROOT" && $checksum_tool "$(basename "$file")") >> "$DIST_ROOT/SHA256SUMS"
  fi
}

# macOS builds (amd64 & arm64)
for arch in amd64 arm64; do
  build_dir="$TMP_DIR/darwin_${arch}"
  mkdir -p "$build_dir"
  build_binary "darwin" "$arch" "$build_dir/$APP_NAME" "macOS ${arch}"
  archive="${DIST_ROOT}/${APP_NAME}_${VERSION}_darwin_${arch}.tar.gz"
  package_tar "$build_dir" "$archive"
  record_checksum "$archive"
  rm -rf "$build_dir"
done

# Optional universal binary if lipo is available
if command -v lipo >/dev/null 2>&1; then
  echo "-> Creating macOS universal binary"
  build_binary "darwin" "amd64" "$TMP_DIR/${APP_NAME}_amd64" "macOS amd64"
  build_binary "darwin" "arm64" "$TMP_DIR/${APP_NAME}_arm64" "macOS arm64"
  universal_dir="$TMP_DIR/darwin_universal"
  mkdir -p "$universal_dir"
  lipo -create -output "$universal_dir/${APP_NAME}" "$TMP_DIR/${APP_NAME}_amd64" "$TMP_DIR/${APP_NAME}_arm64"
  archive="${DIST_ROOT}/${APP_NAME}_${VERSION}_darwin_universal.tar.gz"
  package_tar "$universal_dir" "$archive"
  record_checksum "$archive"
  rm -rf "$universal_dir" "$TMP_DIR/${APP_NAME}_amd64" "$TMP_DIR/${APP_NAME}_arm64"
fi

# Linux amd64
linux_dir="$TMP_DIR/linux_amd64"
mkdir -p "$linux_dir"
build_binary "linux" "amd64" "$linux_dir/$APP_NAME" "Linux amd64"
archive="${DIST_ROOT}/${APP_NAME}_${VERSION}_linux_amd64.tar.gz"
package_tar "$linux_dir" "$archive"
record_checksum "$archive"
rm -rf "$linux_dir"

# Windows amd64
windows_dir="$TMP_DIR/windows_amd64"
mkdir -p "$windows_dir"
build_binary "windows" "amd64" "$windows_dir/${APP_NAME}.exe" "Windows amd64"
archive="${DIST_ROOT}/${APP_NAME}_${VERSION}_windows_amd64.zip"
package_zip "$windows_dir" "$archive"
record_checksum "$archive"
rm -rf "$windows_dir"

if [[ -f "$DIST_ROOT/SHA256SUMS" ]]; then
  echo "==> Checksums written to $DIST_ROOT/SHA256SUMS"
fi

echo "==> Artifacts created in $DIST_ROOT"
