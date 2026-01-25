#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Function to print error messages and exit
error_exit() {
    echo "Error: $1" >&2
    exit 1
}

# Check for ROLLING environment variable first
if [ "$ROLLING" = "1" ] || [ "$ROLLING" = "true" ]; then
    ROLLING=true
else
    ROLLING=false
fi

# Check if running as root
if [[ $EUID -eq 0 ]]; then
    error_exit "This script should not be run as root. Please run as a regular user."
fi
# Detect OS and Architecture
OS_RAW=$(uname -s)
ARCH=$(uname -m)
EXT=""

case "$OS_RAW" in
    Linux*)     OS="linux";;
    Darwin*)    OS="darwin";;
    CYGWIN*|MINGW*|MSYS*) OS="windows"; EXT=".exe";;
    *)          error_exit "Unsupported OS: $OS_RAW";;
esac

if [ "$ARCH" == "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" == "aarch64" ] || [ "$ARCH" == "arm64" ]; then
    ARCH="arm64"
else
    error_exit "Unsupported architecture: $ARCH"
fi

if [ "$ROLLING" = true ]; then
    echo "Rolling update selected. Building from source..."

    # Check dependencies
    if ! command -v go >/dev/null 2>&1; then
        error_exit "Go is required for rolling installation but not found."
    fi
    if ! command -v git >/dev/null 2>&1; then
        error_exit "Git is required for rolling installation but not found."
    fi

    # Create temp directory
    TEMP_DIR=$(mktemp -d)
    cleanup() {
        rm -rf "$TEMP_DIR"
    }
    trap cleanup EXIT

    echo "Cloning repository..."
    if ! git clone https://github.com/finettt/agent-go.git "$TEMP_DIR/agent-go"; then
        error_exit "Failed to clone repository"
    fi

    echo "Building..."
    cd "$TEMP_DIR/agent-go/src" || error_exit "Failed to change to source directory"
    
    if ! go build -ldflags="-s -w" -o "../agent-go${EXT}" .; then
        error_exit "Build failed"
    fi
    cd ..

    echo "Installing agent-go (rolling)..."
    chmod +x "agent-go${EXT}"
    if [ "$OS" == "windows" ]; then
        # On Windows (Git Bash/MinGW), sudo might not be available or needed
        # Try to install to a common location or just advise user
        INSTALL_DIR="$HOME/bin"
        mkdir -p "$INSTALL_DIR"
        mv "./agent-go${EXT}" "$INSTALL_DIR/agent-go${EXT}"
        echo "Installed to $INSTALL_DIR"
        echo "Please ensure $INSTALL_DIR is in your PATH."
    else
        sudo mv "./agent-go${EXT}" "/usr/local/bin/agent-go${EXT}"
    fi

else
    echo "Fetching latest version..."
    LATEST_VERSION=$(curl -fsSL https://raw.githubusercontent.com/finettt/agent-go/main/latest | tr -d '[:space:]')

    if [ -z "$LATEST_VERSION" ]; then
        error_exit "Failed to fetch latest version."
    fi

    DOWNLOAD_URL="https://github.com/finettt/agent-go/releases/download/${LATEST_VERSION}/agent-go-${OS}-${ARCH}${EXT}"

    echo "Downloading agent-go ${LATEST_VERSION} for ${OS}/${ARCH}..."
    if ! curl -fsSL -o "agent-go${EXT}" "$DOWNLOAD_URL"; then
        error_exit "Failed to download agent-go from $DOWNLOAD_URL"
    fi

    echo "Installing agent-go..."
    chmod +x "agent-go${EXT}"
    if [ "$OS" == "windows" ]; then
        INSTALL_DIR="$HOME/bin"
        mkdir -p "$INSTALL_DIR"
        mv "./agent-go${EXT}" "$INSTALL_DIR/agent-go${EXT}"
        echo "Installed to $INSTALL_DIR"
        echo "Please ensure $INSTALL_DIR is in your PATH."
    else
        sudo mv "./agent-go${EXT}" "/usr/local/bin/agent-go${EXT}"
    fi
fi

# Verify installation
if command -v agent-go >/dev/null 2>&1; then
    echo "Agent-go installed successfully!"
    echo "You can now run 'agent-go' to start the application."
else
    error_exit "Installation failed: agent-go not found in PATH."
fi
