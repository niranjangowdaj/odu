#!/bin/bash
set -e

REPO="niranjangowdaj/odu"
INSTALL_DIR="$HOME/.local/bin"

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

# spinner runs while curl downloads in background
_spinner() {
  local frames='⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏'
  local i=0
  while kill -0 "$1" 2>/dev/null; do
    printf "\r  %s" "${frames:$((i % ${#frames})):1}"
    i=$((i + 1))
    sleep 0.08
  done
  printf "\r  ✓\n"
}

curl -fsSL "$URL" -o "$TMP" &
CURL_PID=$!
_spinner "$CURL_PID"
wait "$CURL_PID" || { echo "Download failed. Check the URL: $URL"; exit 1; }
chmod +x "$TMP"

mkdir -p "$INSTALL_DIR"
mv "$TMP" "$INSTALL_DIR/odu"

echo ""
echo "✓ odu $LATEST installed to $INSTALL_DIR/odu"

# Warn if INSTALL_DIR is not in PATH
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo ""
    echo "  ⚠️  $INSTALL_DIR is not in your PATH."
    echo "  Add this to your ~/.zshrc or ~/.bashrc:"
    echo ""
    echo "    export PATH=\"\$HOME/.local/bin:\$PATH\""
    echo ""
    echo "  Then restart your shell or run: source ~/.zshrc"
    ;;
esac
