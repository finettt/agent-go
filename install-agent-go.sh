#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Function to print error messages and exit
error_exit() {
    echo "Error: $1" >&2
    exit 1
}

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
sudo mv "./agent-go${EXT}" "/usr/local/bin/agent-go${EXT}"

# Verify installation
if command -v agent-go >/dev/null 2>&1; then
    echo "Agent-go installed successfully!"
    echo "You can now run 'agent-go' to start the application."
else
    error_exit "Installation failed: agent-go not found in PATH."
fi
