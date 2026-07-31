#!/bin/bash
set -e

REPO="niranjangowdaj/odu"
INSTALL_DIR="/usr/local/bin"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

if [ "$OS" != "darwin" ] && [ "$OS" != "linux" ]; then
  echo "Unsupported OS: $OS"
  echo "Windows users: download odu-windows-amd64.exe from https://github.com/$REPO/releases"
  exit 1
fi

BINARY="odu-${OS}-${ARCH}"

echo "Fetching latest release..."
LATEST=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed 's/.*"tag_name": *"\(.*\)".*/\1/')

if [ -z "$LATEST" ]; then
  echo "Could not determine latest release. Check https://github.com/$REPO/releases"
  exit 1
fi

URL="https://github.com/$REPO/releases/download/$LATEST/$BINARY"

echo "Downloading odu $LATEST ($OS/$ARCH)..."
TMP=$(mktemp)
curl -fsSL "$URL" -o "$TMP"
chmod +x "$TMP"

echo "Installing to $INSTALL_DIR/odu (may require sudo)..."
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP" "$INSTALL_DIR/odu"
else
  sudo mv "$TMP" "$INSTALL_DIR/odu"
fi

echo ""
echo "✓ odu $LATEST installed successfully!"
echo "  Run 'odu --help' to get started."
