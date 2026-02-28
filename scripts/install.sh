#!/bin/bash

# AI AgentGuard Installation Script
# Usage: curl -sSL https://raw.githubusercontent.com/imdlan/AIAgentGuard/main/scripts/install.sh | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

REPO="imdlan/AIAgentGuard"
BINARY_NAME="agent-guard"
INSTALL_DIR="/usr/local/bin"

echo "🛡️  AI AgentGuard Installer"
echo "=============================="
echo

# Fetch latest release version from GitHub API
echo -e "${YELLOW}Fetching latest version...${NC}"
VERSION=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo -e "${RED}Error: Could not determine latest version${NC}"
    exit 1
fi

echo -e "${GREEN}Latest version:${NC} $VERSION"

# Detect OS
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
    Darwin)
        PLATFORM="darwin"
        ;;
    Linux)
        PLATFORM="linux"
        ;;
    *)
        echo -e "${RED}Unsupported OS: $OS${NC}"
        exit 1
        ;;
esac

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}Detected platform:${NC} $PLATFORM-$ARCH"
echo

# Check if binary already exists
if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
    echo -e "${YELLOW}Warning: $BINARY_NAME already exists in $INSTALL_DIR${NC}"
    read -p "Overwrite? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Installation cancelled"
        exit 0
    fi
fi

# Download URL (tar.gz archive)
# Version format: v1.3.0 -> archive name: agent-guard_1.3.0_darwin_arm64.tar.gz
VERSION_NUM="${VERSION#v}"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/${BINARY_NAME}_${VERSION_NUM}_${PLATFORM}_${ARCH}.tar.gz"

echo -e "${YELLOW}Downloading from:${NC} $DOWNLOAD_URL"
echo

# Create temp directory
temp_dir=$(mktemp -d)
trap "rm -rf $temp_dir" EXIT

# Download archive
if command -v curl &> /dev/null; then
    curl -L -o "$temp_dir/archive.tar.gz" "$DOWNLOAD_URL"
elif command -v wget &> /dev/null; then
    wget -O "$temp_dir/archive.tar.gz" "$DOWNLOAD_URL"
else
    echo -e "${RED}Error: Neither curl nor wget found${NC}"
    exit 1
fi

# Extract binary
tar -xzf "$temp_dir/archive.tar.gz" -C "$temp_dir"

# Make executable
chmod +x "$temp_dir/$BINARY_NAME"

# Install
echo -e "${YELLOW}Installing to $INSTALL_DIR...${NC}"
if [ -w "$INSTALL_DIR" ]; then
    mv "$temp_dir/$BINARY_NAME" "$INSTALL_DIR/"
else
    echo -e "${YELLOW}Requires sudo privileges${NC}"
    sudo mv "$temp_dir/$BINARY_NAME" "$INSTALL_DIR/"
fi

# Verify installation
if command -v $BINARY_NAME &> /dev/null; then
    echo
    echo -e "${GREEN}✅ Installation successful!${NC}"
    echo
    echo "Installed version: $VERSION"
    echo
    echo "Run '$BINARY_NAME --help' to get started"
    echo
    echo "Quick start:"
    echo "  $BINARY_NAME scan          # Scan for security risks"
    echo "  $BINARY_NAME run 'echo hi' # Run command in sandbox"
    echo "  $BINARY_NAME report        # Generate security report"
else
    echo -e "${RED}❌ Installation failed${NC}"
    exit 1
fi