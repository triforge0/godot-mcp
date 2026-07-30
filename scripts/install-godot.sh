#!/usr/bin/env bash
# Install Godot 4 editor binary for CI or local integration tests.
set -euo pipefail

VERSION="${GODOT_VERSION:-4.3-stable}"
OS="$(uname -s)"
ARCH="$(uname -m)"
INSTALL_DIR="${GODOT_INSTALL_DIR:-/usr/local/bin}"

case "${OS}" in
  Linux)
    case "${ARCH}" in
      x86_64) ARCHIVE="Godot_v${VERSION}_linux.x86_64.zip" ;;
      aarch64|arm64) ARCHIVE="Godot_v${VERSION}_linux.arm64.zip" ;;
      *) echo "unsupported Linux arch: ${ARCH}" >&2; exit 1 ;;
    esac
    ;;
  Darwin)
    ARCHIVE="Godot_v${VERSION}_macos.universal.zip"
    ;;
  *)
    echo "unsupported OS: ${OS}" >&2
    exit 1
    ;;
esac

URL="https://github.com/godotengine/godot/releases/download/${VERSION}/${ARCHIVE}"
TMP="$(mktemp -d)"
trap 'rm -rf "${TMP}"' EXIT

echo "Downloading ${URL}"
curl -fsSL "${URL}" -o "${TMP}/godot.zip"
unzip -q "${TMP}/godot.zip" -d "${TMP}"

DEST="${INSTALL_DIR}/godot"
mkdir -p "${INSTALL_DIR}"

if [[ "${OS}" == "Darwin" ]]; then
  APP="$(find "${TMP}" -maxdepth 1 -name 'Godot*.app' | head -1)"
  if [[ -z "${APP}" ]]; then
    echo "Godot.app not found in archive" >&2
    exit 1
  fi
  cp "${APP}/Contents/MacOS/Godot" "${DEST}"
else
  BIN="$(find "${TMP}" -maxdepth 1 -type f -name 'Godot*' ! -name '*.zip' | head -1)"
  if [[ -z "${BIN}" ]]; then
    echo "Godot binary not found in archive" >&2
    exit 1
  fi
  cp "${BIN}" "${DEST}"
fi

chmod +x "${DEST}"
echo "Installed ${DEST}"
"${DEST}" --version || true
