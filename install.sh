#!/bin/bash -eu

REPO="DavidRR-F/shm"
VERSION="${1:-latest}"

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [[ "$ARCH" == "x86_64" ]]; then
  ARCH="amd64"
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
  ARCH="arm64"
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

if [[ "$VERSION" == "latest" ]]; then
  VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep tag_name | cut -d '"' -f4)
fi

BINARY_NAME="shm"
ASSET_NAME="${BINARY_NAME}-${OS}-${ARCH}"

URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET_NAME}.tar.gz"

echo "Downloading ${URL}..."

curl -sSfL -o "${ASSET_NAME}.tar.gz" "${URL}" 

echo "Extracting ${ASSET_NAME}.tar.gz"
tar -xzf "${ASSET_NAME}.tar.gz"
rm "${ASSET_NAME}.tar.gz"

chmod +x "$ASSET_NAME"

INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"
echo "Installing to $INSTALL_DIR"
mv "${ASSET_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"


# Warn if not in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
  echo "⚠️  Warning: $INSTALL_DIR is not in your PATH."
  echo "   You can add it with:"
  echo "   echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.bashrc && source ~/.bashrc"
fi

echo "$BINARY_NAME installed to $INSTALL_DIR"
