# Development Guide

This comprehensive guide provides instructions for setting up the development environment, contributing to Agent-Go, and maintaining code quality.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Development Environment Setup](#development-environment-setup)
- [Project Structure](#project-structure)
- [Building and Running](#building-and-running)
- [Testing](#testing)
- [Code Quality and Style](#code-quality-and-style)
- [Contributing Guidelines](#contributing-guidelines)
- [Debugging](#debugging)
- [Performance Considerations](#performance-considerations)
- [Security Guidelines](#security-guidelines)
- [Release Process](#release-process)

## Prerequisites

### System Requirements

- **Go**: 1.25.3 (recommended)
- **Git**: For version control
- **Make**: For build automation (optional, but recommended)
- **Text Editor**: VS Code, GoLand, or any editor with Go support

### Recommended Tools

- **Go Extension**: For VS Code (provides IntelliSense, debugging, and testing)
- **gopls**: Go language server for better IDE support
- **golangci-lint**: Static analysis tool for code quality
- **Air**: Live reload tool for development (optional)

### Development Environment Setup

#### 1. Clone the Repository

```bash
# Clone with HTTPS
git clone https://github.com/finettt/agent-go.git
cd agent-go

# Or clone with SSH
git clone git@github.com:finettt/agent-go.git
cd agent-go
```

#### 2. Install Dependencies

```bash
# Download and verify Go modules
go mod tidy

# Verify dependencies
go mod verify
```

#### 3. Set Up Development Tools

```bash
# Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Install Air for live reload (optional)
go install github.com/cosmtrek/air@latest
```

#### 4. IDE Setup (VS Code Example)

```bash
# Install VS Code extensions
code --install-extension golang.go
code --install-extension ms-vscode.vscode-json
code --install-extension bradlc.vscode-tailwindcss
```

Create `.vscode/settings.json`:

```json
{
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "workspace",
    "go.lintFlags": ["--fast"],
    "go.formatTool": "goimports",
    "go.useLanguageServer": true,
    "files.associations": {
        "*.go": "go"
    }
}
```

## Project Structure

```
agent-go/
├── src/                    # Source code directory
│   ├── main.go            # Application entry point
│   ├── api.go             # API communication logic
│   ├── config.go          # Configuration management
│   ├── executor.go        # Command execution
│   ├── rag.go             # RAG functionality
│   ├── completion.go      # Auto-completion
│   └── types.go           # Data structures
├── docs/                  # Documentation
│   ├── README.md          # Developer documentation
│   ├── architecture.md    # Architecture overview
│   ├── commands.md        # Command reference
│   ├── configuration.md   # Configuration guide
│   └── development.md     # This development guide
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
├── Makefile               # Build automation
├── .gitignore             # Git ignore rules
└── LICENSE                # MIT License
```

### Key Components

- **`src/main.go`**: Application entry point, CLI loop, and signal handling
- **`src/api.go`**: OpenAI API integration and response processing
- **`src/config.go`**: Configuration loading and validation
- **`src/executor.go`**: Secure command execution
- **`src/rag.go`**: Document search and context retrieval
- **`src/completion.go`**: CLI auto-completion functionality
- **`src/types.go`**: Type definitions and data structures

## Building and Running

### Using Make (Recommended)

```bash
# Build the application
make build

# Run the application (builds first)
make run

# Clean build artifacts
make clean

# Install to system (optional)
make install
```

**Build Output**: The binary is created as `agent-go` in the project root directory.

### Using Go Directly

```bash
# Build for current platform
go build -o agent-go ./src

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o agent-go-linux ./src

# Run directly
go run ./src
```

### Development with Live Reload

```bash
# Using Air (if installed)
air

# Or using Go's built-in live reload
go run -mod=mod ./src
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run tests and generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -run TestFunctionName ./src
```

### Writing Tests

Follow these guidelines for writing tests:

```go
// Example test structure
package main

import (
    "testing"
    "os"
)

func TestConfigLoading(t *testing.T) {
    // Setup environment for test
    os.Setenv("OPENAI_KEY", "test-key")
    os.Setenv("OPENAI_MODEL", "gpt-4")
    
    // Test the function
    result := loadConfig()
    
    // Assert the results
    if result.APIKey != "test-key" {
        t.Errorf("Expected APIKey 'test-key', got %s", result.APIKey)
    }
    if result.Model != "gpt-4" {
        t.Errorf("Expected Model 'gpt-4', got %s", result.Model)
    }
    
    // Cleanup
    os.Unsetenv("OPENAI_KEY")
    os.Unsetenv("OPENAI_MODEL")
}

// Mock test for sendAPIRequest would require mocking HTTP client
func BenchmarkAPIRequest(b *testing.B) {
    // Note: Actual benchmark would require mock setup
    for i := 0; i < b.N; i++ {
        // Benchmark code would go here
    }
}
```

### Test Coverage Requirements

- **Unit Tests**: Minimum 80% coverage for individual packages
- **Integration Tests**: Test all major workflows
- **Edge Cases**: Test error conditions and boundary values

## Code Quality and Style

### Formatting

```bash
# Format code
go fmt ./...

# Format and simplify imports
goimports -w .

# Check formatting without modifying files
gofmt -d .
```

### Linting

```bash
# Run golangci-lint
golangci-lint run

# Run specific linters
golangci-lint run --enable=goimports
golangci-lint run --enable=gosimple
golangci-lint run --enable=staticcheck
```

### Code Style Guidelines

1. **Naming Conventions**:
   - Use `camelCase` for exported functions and variables
   - Use `PascalCase` for exported types and structs
   - Use `snake_case` for private variables
   - Use descriptive names: `calculateTotal` instead of `calcTot`

2. **Error Handling**:
   - Always check for errors
   - Provide meaningful error messages
   - Use `fmt.Errorf` for error wrapping
   - Handle errors at the appropriate level

3. **Comments**:
   - Package comments for all exported packages
   - Function comments for exported functions
   - Complex logic should be well-commented
   - Use `//` for single-line comments
   - Use `/* */` for multi-line comments

4. **Error Handling Example**:

```go
func sendAPIRequest(agent *Agent, config *Config) (*ChatCompletionResponse, error) {
    if config.APIKey == "" {
        return nil, fmt.Errorf("API key is required")
    }
    
    resp, err := http.Post(config.APIURL, "application/json", requestBody)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
    }
    
    var result ChatCompletionResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &result, nil
}
```

## Contributing Guidelines

### Development Workflow

1. **Fork the repository** on GitHub
2. **Create a feature branch**:

   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes** following the guidelines above
4. **Test your changes**:

   ```bash
   go test ./...
   go vet ./...
   golangci-lint run
   ```

5. **Commit your changes**:

   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

6. **Push to your fork**:

   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request** with a clear description

### Commit Message Guidelines

Use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build process or auxiliary tool changes

**Examples:**

```
feat: add RAG configuration support
fix: resolve command execution timeout issue
docs: update API documentation
style: format code according to guidelines
test: add unit tests for API client
```

### Pull Request Process

1. **PR Title**: Use clear, descriptive titles
2. **PR Description**: Include:
   - Problem description
   - Solution overview
   - Testing performed
   - Breaking changes (if any)
   - Related issues (if any)
3. **PR Checklist**:
   - All tests pass
   - Code follows style guidelines
   - Documentation is updated
   - Changes are tested
   - Breaking changes are documented

## Debugging

### Debug Mode

Enable debug logging for troubleshooting:

```bash
# Set debug environment variable
export DEBUG=1

# Run with debug output
./agent-go
```

### Debugging Tools

```bash
# Debug with Delve (Go debugger)
dlv debug ./src

# Debug specific issues
go run -race ./src  # Data race detection
go run -v ./src     # Verbose output
```

### Common Debug Scenarios

**API Connection Issues**:

```bash
# Test API connectivity
curl -H "Authorization: Bearer $OPENAI_KEY" $OPENAI_BASE/v1/models

# Check network configuration
netstat -an | grep $OPENAI_BASE
```

**Configuration Issues**:

```bash
# Validate configuration file
cat ~/.config/agent-go/config.json | jq .

# Test configuration loading
go run -v ./src 2>&1 | grep config
```

## Performance Considerations

### Optimization Techniques

1. **Memory Management**:
   - Use object pooling for frequently allocated objects
   - Avoid unnecessary allocations in hot paths
   - Use `sync.Pool` for temporary objects

2. **Concurrency**:
   - Use goroutines for I/O-bound operations
   - Use channels for communication between goroutines
   - Avoid shared memory when possible

3. **Network Optimization**:
   - Use connection pooling for HTTP clients
   - Implement request batching where appropriate
   - Cache responses when possible

### Performance Testing

```bash
# Run benchmarks
go test -bench=. ./src

# Profile CPU usage
go test -cpuprofile=cpu.prof ./src
go tool pprof cpu.prof

# Profile memory usage
go test -memprofile=mem.prof ./src
go tool pprof mem.prof
```

## Security Guidelines

### Security Best Practices

1. **Input Validation**:
   - Validate all user inputs
   - Sanitize file paths to prevent directory traversal
   - Validate API responses

2. **Secret Management**:
   - Never hardcode API keys
   - Use environment variables for sensitive data
   - Set proper file permissions for configuration files

3. **Command Execution**:
   - Validate commands before execution
   - Use safe execution methods
   - Implement proper error handling

### Security Testing

```bash
# Run security checks
go vet -vettool=$(which gosec) ./src

# Check for known vulnerabilities
go list -json -m all | nancy sleuth
```

## Release Process

### Version Management

1. **Update Version Numbers**:
   - Update `go.mod` with new version
   - Update documentation with version changes
   - Update changelog

2. **Create Release Tag**:

   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

3. **Build Release Artifacts**:

   ```bash
   # Build for multiple platforms
   GOOS=linux GOARCH=amd64 go build -o agent-go-linux ./src
   GOOS=darwin GOARCH=amd64 go build -o agent-go-darwin ./src
   GOOS=windows GOARCH=amd64 go build -o agent-go.exe ./src
   ```

### Release Checklist

- [ ] All tests pass
- [ ] Documentation is updated
- [ ] Changelog is updated
- [ ] Version numbers are updated
- [ ] Release artifacts are built
- [ ] GitHub release is created
- [ ] Documentation is published

## Additional Resources

### Documentation Links

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Conventional Commits](https://www.conventionalcommits.org/)

### Community and Support

- **GitHub Issues**: Report bugs and request features
- **Discussions**: Join community discussions
- **Discord**: Join our community server (if available)

### Contributing to Documentation

When updating documentation:

1. Keep examples up-to-date with code changes
2. Use clear, concise language
3. Include code examples where helpful
4. Update related documentation when making changes
5. Test all examples to ensure they work

---

Happy coding! 🚀
