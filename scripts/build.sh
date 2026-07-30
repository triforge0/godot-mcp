#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT="${ROOT}/bin/godot-mcp"

mkdir -p "${ROOT}/bin"

echo "Building godot-mcp..."
(cd "${ROOT}" && go build -o "${OUT}" ./cmd/godot-mcp)

echo "Done: ${OUT}"
"${OUT}" version
