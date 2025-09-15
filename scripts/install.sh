#!/usr/bin/env bash

set -e

# Determine OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Set binary name based on OS and ARCH
if [[ "$OS" == "linux" ]]; then
  BINARY="cortex-linux-amd64"
elif [[ "$OS" == "darwin" ]]; then
  BINARY="cortex-darwin-arm64"
elif [[ "$OS" == "msys" || "$OS" == "cygwin" ]]; then
  BINARY="cortex-windows-amd64.exe"
else
  echo "Unsupported OS: $OS"
  exit 1
fi

VERSION=$(curl -s https://api.github.com/repos/tejastn10/cortex/releases/latest | jq -r '.tag_name')
BASE_URL="https://github.com/tejastn10/cortex/releases/download/$VERSION"
FILE="$BINARY.zip"

# Download the binary
echo "Downloading $BINARY from $BASE_URL/$FILE ..."
curl -LO "$BASE_URL/$FILE"

# Extract if on Unix-like system
if [[ "$OS" != "msys" && "$OS" != "cygwin" ]]; then
  unzip "$FILE" -d cortex_temp
  BINARY_PATH=$(find cortex_temp -type f -name "$BINARY")
else
  # Windows users need a different extraction method
  BINARY_PATH="$BINARY"
fi

# Install: move binary to /usr/local/bin (requires sudo)
echo "Installing cortex to /usr/local/bin ..."
sudo mv "$BINARY_PATH" /usr/local/bin/cortex
sudo chmod +x /usr/local/bin/cortex

# Cleanup
rm -rf "$FILE" cortex_temp

echo "cortex installed successfully!"
