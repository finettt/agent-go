# Development Guide

This guide provides instructions for setting up the development environment and contributing to Agent-Go.

## Prerequisites

- Go 1.18 or higher
- A configured text editor (e.g., VS Code with the Go extension)

## Getting Started

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/finettt/agent-go.git
    cd agent-go
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Build the application:**
    Use the provided `Makefile` for convenience.
    ```bash
    make build
    ```
    This will create the `agent-go` binary in the root directory.

4.  **Run the application:**
    ```bash
    make run
    ```
    or
    ```bash
    ./agent-go
    ```

## Running Tests

To run the test suite:
```bash
go test ./...
```

## Code Style

Please run `go fmt` to format your code before committing.
```bash
go fmt ./...