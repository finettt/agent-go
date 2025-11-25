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

# Check for required tools
echo "Checking prerequisites..."
command -v curl >/dev/null 2>&1 || error_exit "curl is required but not installed."
command -v git >/dev/null 2>&1 || error_exit "git is required but not installed."
command -v make >/dev/null 2>&1 || error_exit "make is required but not installed."
command -v go >/dev/null 2>&1 || error_exit "go is required but not installed."

# Check if sudo is available
command -v sudo >/dev/null 2>&1 || error_exit "sudo is required but not installed."

# Create a temporary directory
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT

echo "Downloading agent-go..."
cd "$TEMP_DIR"

# Clone the repository
git clone https://github.com/finettt/agent-go.git
cd agent-go

echo "Building agent-go..."
# Build the binary (avoid compress step which requires upx)
make build

# Check if binary was created
if [[ ! -f "./agent-go" ]]; then
    error_exit "Build failed: agent-go binary not found."
fi

echo "Installing agent-go..."
# Set ownership and move to /usr/local/bin
sudo chown "$USER:$USER" ./agent-go
sudo mv ./agent-go /usr/local/bin/

# Verify installation
if command -v agent-go >/dev/null 2>&1; then
    echo "Agent-go installed successfully!"
    echo "You can now run 'agent-go' to start the application."
else
    error_exit "Installation failed: agent-go not found in PATH."
fi
