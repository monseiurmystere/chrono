#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SRC_DIR="${ROOT_DIR}/src"
BIN_DIR="${ROOT_DIR}/bin"
OUTPUT_BIN="${BIN_DIR}/chrono"

VERSION="${CHRONO_VERSION:-0.1.0-dev}"
GIT_COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")"
BUILD_DATE="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"

LDFLAGS="-s -w -X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT} -X main.BuildDate=${BUILD_DATE}"

print_usage() {
    cat <<EOF
Chrono Native Build Script

Usage:
  ./build/build.sh [command]

Commands:
  build     Build the chrono binary to bin/chrono (default)
  test      Run all unit tests in src/
  clean     Remove built binaries
  run       Build and run chrono in test mode (-t)
  help      Show this help message
EOF
}

cmd_build() {
    echo "==> Building Chrono (${VERSION}, commit: ${GIT_COMMIT})..."
    mkdir -p "${BIN_DIR}"
    (
        cd "${SRC_DIR}"
        go build -ldflags "${LDFLAGS}" -o "${OUTPUT_BIN}" ./cmd/chrono
    )
    echo "==> Binary created successfully at: ${OUTPUT_BIN}"
}

cmd_test() {
    echo "==> Running Chrono unit tests..."
    (
        cd "${SRC_DIR}"
        go test -v ./...
    )
}

cmd_clean() {
    echo "==> Cleaning build artifacts..."
    rm -rf "${BIN_DIR}"
    echo "==> Clean complete."
}

cmd_run() {
    cmd_build
    echo "==> Launching Chrono in test mode..."
    "${OUTPUT_BIN}" -t
}

ACTION="${1:-build}"

case "${ACTION}" in
    build)
        cmd_build
        ;;
    test)
        cmd_test
        ;;
    clean)
        cmd_clean
        ;;
    run)
        cmd_run
        ;;
    help|--help|-h)
        print_usage
        ;;
    *)
        echo "Error: Unknown command '${ACTION}'"
        print_usage
        exit 1
        ;;
esac
