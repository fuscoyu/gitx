#!/bin/bash

# GitX Installation Script
# This script downloads and installs the latest GitX binary for your platform

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="goeoeo/gitx"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="gitx"

# Detect platform
detect_platform() {
    local os arch
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)
    
    case $arch in
        x86_64)
            arch="amd64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        *)
            echo -e "${RED}Error: Unsupported architecture: $arch${NC}"
            exit 1
            ;;
    esac
    
    case $os in
        linux|darwin)
            echo "${os}_${arch}"
            ;;
        *)
            echo -e "${RED}Error: Unsupported OS: $os${NC}"
            exit 1
            ;;
    esac
}

# Get latest release
get_latest_release() {
    local api_url="https://api.github.com/repos/$REPO/releases/latest"
    local version
    
    if command -v curl >/dev/null 2>&1; then
        version=$(curl -s "$api_url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        version=$(wget -qO- "$api_url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        echo -e "${RED}Error: curl or wget is required${NC}"
        exit 1
    fi
    
    if [ -z "$version" ]; then
        echo -e "${RED}Error: Could not get latest version${NC}"
        exit 1
    fi
    
    echo "$version"
}

# Get version info from binary
get_binary_version() {
    local binary_path="$1"
    if [ -f "$binary_path" ] && [ -x "$binary_path" ]; then
        "$binary_path" version 2>/dev/null | head -1 | grep -o '[0-9]\+\.[0-9]\+\.[0-9]\+[a-zA-Z0-9.-]*' || echo "unknown"
    else
        echo "unknown"
    fi
}

# Download and install
install_gitx() {
    local platform version download_url temp_file
    
    platform=$(detect_platform)
    version=$(get_latest_release)
    
    echo -e "${BLUE}Installing GitX $version for $platform...${NC}"
    
    # Construct download URL
    if [[ "$platform" == *"windows"* ]]; then
        download_url="https://github.com/$REPO/releases/download/$version/gitx-$version-$platform.tar.gz"
    else
        download_url="https://github.com/$REPO/releases/download/$version/gitx-$version-$platform.tar.gz"
    fi
    
    # Create temporary file
    temp_file=$(mktemp)
    
    echo -e "${YELLOW}Downloading from: $download_url${NC}"
    
    # Download
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$temp_file" "$download_url"
    else
        wget -O "$temp_file" "$download_url"
    fi
    
    # Extract and install
    echo -e "${YELLOW}Extracting and installing...${NC}"
    
    # Create temp directory for extraction
    local temp_dir
    temp_dir=$(mktemp -d)
    
    # Extract
    tar -xzf "$temp_file" -C "$temp_dir"
    
    # Find the binary
    local binary_path
    binary_path=$(find "$temp_dir" -name "gitx*" -type f | head -1)
    
    if [ -z "$binary_path" ]; then
        echo -e "${RED}Error: Could not find binary in archive${NC}"
        rm -rf "$temp_dir"
        rm -f "$temp_file"
        exit 1
    fi
    
    # Make executable
    chmod +x "$binary_path"
    
    # Install
    if [ -w "$INSTALL_DIR" ]; then
        cp "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
    else
        echo -e "${YELLOW}Requiring sudo to install to $INSTALL_DIR${NC}"
        sudo cp "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Cleanup
    rm -rf "$temp_dir"
    rm -f "$temp_file"
    
    echo -e "${GREEN}✅ GitX $version installed successfully!${NC}"
    echo -e "${BLUE}Run 'gitx --help' to get started${NC}"
}

# Check if already installed
check_existing() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local current_version
        current_version=$(get_binary_version "$(which $BINARY_NAME)")
        echo -e "${YELLOW}GitX is already installed (version: $current_version)${NC}"
        read -p "Do you want to update? [y/N]: " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${BLUE}Installation cancelled${NC}"
            exit 0
        fi
    fi
}

# Main
main() {
    echo -e "${BLUE}GitX Installation Script${NC}"
    echo -e "${BLUE}========================${NC}"
    
    check_existing
    install_gitx
}

# Run main function
main "$@"
