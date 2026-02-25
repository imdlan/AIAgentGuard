#!/bin/bash

# AI AgentGuard Installation Script
# Usage: curl -sSL https://raw.githubusercontent.com/imdlan/AIAgentGuard/main/scripts/install.sh | bash

set -e

VERSION="v1.0.0"
REPO="imdlan/AIAgentGuard"
BINARY_NAME="agent-guard"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🛡️  AI AgentGuard Installer"
echo "=============================="
echo

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

# Download URL
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/${BINARY_NAME}-${PLATFORM}-${ARCH}"

echo -e "${YELLOW}Downloading from:${NC} $DOWNLOAD_URL"
echo

# Download binary
if command -v curl &> /dev/null; then
    curl -L -o /tmp/$BINARY_NAME "$DOWNLOAD_URL"
elif command -v wget &> /dev/null; then
    wget -O /tmp/$BINARY_NAME "$DOWNLOAD_URL"
else
    echo -e "${RED}Error: Neither curl nor wget found${NC}"
    exit 1
fi

# Make executable
chmod +x /tmp/$BINARY_NAME

# Install
echo -e "${YELLOW}Installing to $INSTALL_DIR...${NC}"
if [ -w "$INSTALL_DIR" ]; then
    mv /tmp/$BINARY_NAME "$INSTALL_DIR/"
else
    echo -e "${YELLOW}Requires sudo privileges${NC}"
    sudo mv /tmp/$BINARY_NAME "$INSTALL_DIR/"
fi

# Verify installation
if command -v $BINARY_NAME &> /dev/null; then
    echo
    echo -e "${GREEN}✅ Installation successful!${NC}"
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
