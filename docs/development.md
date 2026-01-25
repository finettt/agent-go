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
- [Documentation Guidelines](#documentation-guidelines)
- [Community and Support](#community-and-support)

## Prerequisites

### System Requirements

- **Go**: 1.25 or later
- **Git**: For version control
- **Make**: For build automation (optional, but recommended)
- **Node.js/npm**: For MCP servers (optional, but recommended for full functionality)
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
│   ├── constants.go       # Constants and default values
│   ├── executor.go        # Command execution
│   ├── processor.go       # Tool call processing
│   ├── tools.go           # Tool definitions
│   ├── rag.go             # RAG functionality
│   ├── mcp.go             # MCP server integration
│   ├── subagent.go        # Sub-agent management
│   ├── todo.go            # Todo list management
│   ├── notes.go           # Notes management
│   ├── agents.go          # Agent Studio and agent management
│   ├── agent_studio.go    # Agent Studio implementation
│   ├── session.go         # Session management
│   ├── completion.go      # Auto-completion
│   ├── system.go          # System information
│   ├── input.go           # Input handling
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
- **`src/api.go`**: OpenAI API integration, streaming, and response processing
- **`src/config.go`**: Configuration loading and validation
- **`src/constants.go`**: Application constants and defaults (includes `DefaultMaxTokens = -1` for unlimited tokens)
- **`src/executor.go`**: Secure command execution
- **`src/processor.go`**: Tool call processing and coordination
- **`src/tools.go`**: Tool definitions and schemas
- **`src/rag.go`**: Document search and context retrieval
- **`src/mcp.go`**: MCP server connection and tool management
- **`src/subagent.go`**: Sub-agent spawning and lifecycle management
- **`src/todo.go`**: Todo list CRUD operations
- **`src/notes.go`**: Persistent notes management
- **`src/agents.go`**: Agent Studio interface and agent management
- **`src/agent_studio.go`**: Agent Studio implementation
- **`src/session.go`**: Session save/restore functionality
- **`src/completion.go`**: CLI auto-completion functionality
- **`src/system.go`**: System information gathering
- **`src/input.go`**: Enhanced input handling
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

### Using Docker

```bash
# Build the Docker image
docker build -t agent-go .

# Run with current directory mounted as /workspace
docker run -it -v $(pwd):/workspace agent-go

# Run with environment variables
docker run -it \
  -v $(pwd):/workspace \
  -e OPENAI_KEY="your-api-key" \
  -e OPENAI_MODEL="gpt-4-turbo" \
  agent-go

# Run with RAG documents mounted
docker run -it \
  -v $(pwd):/workspace \
  -v /path/to/documents:/documents \
  -e OPENAI_KEY="your-api-key" \
  -e RAG_ENABLED=1 \
  -e RAG_PATH="/documents" \
  agent-go
```

**Docker Configuration:**
- The `/workspace` directory is where your host's current directory will be mounted
- The container runs as a non-root user (`appuser`) for security
- Configuration is stored in `/home/appuser/.config/agent-go` inside the container
- Use volume mounts to persist configuration across container restarts

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

### Installing Test Dependencies

```bash
# Install testing tools
go install gotest.tools/gotestsum@latest

# Install MCP testing tools (for MCP integration tests)
npm install -g @modelcontextprotocol/inspector
```

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

// Test MCP integration
func TestMCPConnection(t *testing.T) {
    // Setup mock MCP server
    config := &Config{
        MCPs: map[string]MCPServer{
            "test": {
                Name: "test",
                Command: "npx -y @test/mcp-server",
            },
        },
    }
    
    // Test connection
    mgr := newMCPManager()
    session, err := mgr.ensureMCP("test")
    
    if err != nil {
        t.Errorf("Expected successful connection, got error: %v", err)
    }
    
    if session == nil {
        t.Error("Expected valid session, got nil")
    }
}

// Test todo list operations
func TestTodoOperations(t *testing.T) {
    agentID := "test-agent"
    
    // Create todo
    result, err := createTodo(agentID, `{"task":"Test task"}`)
    if err != nil {
        t.Errorf("Failed to create todo: %v", err)
    }
    
    // Verify todo was created
    if !strings.Contains(result, "Test task") {
        t.Error("Todo not found in result")
    }
    
    // Update todo
    result, err = updateTodo(agentID, `{"id":1,"status":"completed"}`)
    if err != nil {
        t.Errorf("Failed to update todo: %v", err)
    }
    
    // Verify update
    if !strings.Contains(result, "completed") {
        t.Error("Todo status not updated")
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

## New Features Development

### Agent Studio Development

**Key Files:**
- `src/agents.go` - Agent Studio interface and management
- `src/agent_studio.go` - Agent Studio implementation

**Development Guidelines:**
1. **Agent Definition Structure**: Follow the `Agent` struct in `types.go`
2. **Validation**: Implement proper validation for agent specifications
3. **Storage**: Agent definitions stored in `~/.config/agent-go/agents/*.json`
4. **Studio Restrictions**: Agent Studio should only permit agent creation
5. **Built-in Protection**: Ensure `default` agent cannot be deleted or overwritten

**Testing Agent Studio:**
```bash
# Test agent creation
> /agent studio

# Test agent management
> /agent list
> /agent view <name>
> /agent use <name>
> /agent rm <name>
```

### Session Management Development

**Key Files:**
- `src/session.go` - Session save/restore functionality

**Development Guidelines:**
1. **Session Structure**: Include metadata (timestamps, token counts, message history)
2. **Storage Format**: Use JSON format for session data
3. **Compression**: Automatically compress context before saving sessions
4. **Timestamp Tracking**: Record creation and last access times
5. **Cleanup**: Implement proper cleanup for unused sessions

**Testing Session Management:**
```bash
# Test session creation and restoration
> /session new
> /session list
> /session restore <name>
> /session rm <name>
```

### Background Command Execution

**Key Files:**
- `src/executor.go` - Extended for background process management

**Development Guidelines:**
1. **Process Management**: Track PIDs and manage lifecycle properly
2. **Output Streaming**: Implement real-time output capture
3. **Resource Cleanup**: Automatically clean up completed processes
4. **Safety**: Prevent application exit while background tasks run
5. **Error Handling**: Gracefully handle process failures

**Testing Background Commands:**
```bash
# Test background execution
$ long-running-command --background

# Monitor and manage
> /list_background_commands
> /get_background_logs <pid>
> /kill_background_command <pid>
```

### Notes Management Development

**Key Files:**
- `src/notes.go` - Persistent notes management

**Development Guidelines:**
1. **Storage Format**: Use JSON for note storage with metadata
2. **Filename Sanitization**: Prevent path traversal attacks
3. **System Integration**: Notes should be automatically injected into system prompt
4. **Autocomplete**: Implement tab completion for note names
5. **Cross-Session Persistence**: Notes should persist across different agent sessions

**Testing Notes Management:**
```bash
# Test note operations
> /notes list
> /notes view <name>

# Test AI-driven note management
> Create a note called "api_endpoint" with content "https://api.example.com"
> Update the api_endpoint note
> Delete the old_note
```

### Usage Tracking Development

**Key Files:**
- Extended functionality in `src/main.go` and `src/api.go`

**Development Guidelines:**
1. **Token Tracking**: Monitor prompt, completion, and reasoning tokens
2. **Cost Calculation**: Implement accurate cost tracking based on model pricing
3. **Historical Data**: Track usage across sessions and time periods
4. **Reset Behavior**: Token counter should reset after compression
5. **Performance**: Ensure tracking doesn't impact application performance

**Testing Usage Tracking:**
```bash
# Test usage monitoring
> /usage
> /cost

# Verify reset after compression
> /compress
> /usage  # Should show 0 tokens
```

### Enhanced CLI Features

**Key Files:**
- `src/completion.go` - Extended for new command completion
- `src/main.go` - Enhanced CLI handling

**Development Guidelines:**
1. **Autocomplete**: Add completion for new commands and parameters
2. **Color Scheme**: Maintain consistent color coding
3. **User Feedback**: Provide clear, actionable error messages
4. **Performance**: Ensure autocomplete doesn't slow down CLI
5. **Accessibility**: Consider accessibility requirements

### Testing New Features

**Unit Tests:**
```bash
# Test individual components
go test ./src/...

# Test specific new features
go test -v src/agents.go
go test -v src/session.go
go test -v src/notes.go
```

**Integration Tests:**
```bash
# Test complete workflows
./agent-go "/agent studio" < test-spec.txt
./agent-go "/session new" < test-session.txt
./agent-go "/notes list" < test-notes.txt
```

**Manual Testing:**
1. Test all new slash commands
2. Verify feature interactions
3. Test error scenarios
4. Validate edge cases

### Debugging New Features

**Agent Studio Debugging:**
- Check agent definition JSON files
- Verify validation logic
- Test studio restrictions

**Session Management Debugging:**
- Check session JSON files
- Verify compression/decompression
- Test context restoration

**Background Command Debugging:**
- Monitor process PIDs
- Check output buffers
- Test cleanup mechanisms

**Notes Management Debugging:**
- Check note JSON files
- Verify system prompt injection
- Test autocomplete functionality

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

## Documentation Guidelines

### Documentation Structure

- **README.md**: Main project overview and quick start guide
- **docs/README.md**: Developer documentation overview
- **docs/architecture.md**: System architecture and design decisions
- **docs/commands.md**: Complete command reference
- **docs/configuration.md**: Configuration options and examples
- **docs/development.md**: This development guide
- **docs/examples.md**: Practical examples and best practices

### Writing Good Documentation

1. **Audience Awareness**: Write for both users and developers
2. **Clarity**: Use simple, direct language
3. **Examples**: Include practical, working examples
4. **Consistency**: Follow established formatting and style
5. **Completeness**: Cover all important aspects of the topic
6. **Maintenance**: Keep documentation up-to-date with code changes

### Documentation Format

- Use Markdown for all documentation
- Follow consistent heading structure
- Use code blocks for examples and configuration
- Include tables for parameter references
- Use proper link formatting for cross-references

### Code Examples

```go
// Good example with context
func loadConfig() (*Config, error) {
    // Load configuration from multiple sources
    // with proper error handling
    return &config, nil
}
```

## Community and Support

### Getting Help

- **GitHub Issues**: Report bugs and request features
- **GitHub Discussions**: Join community discussions
- **Discord**: Join our community server (if available)
- **Email**: Contact maintainers for private matters

### Contributing

We welcome contributions of all types:

1. **Code Contributions**: New features, bug fixes, optimizations
2. **Documentation**: Improvements to existing docs or new guides
3. **Bug Reports**: Detailed issue reports with reproduction steps
4. **Feature Requests**: Well-thought-out suggestions for new functionality
5. **Testing**: Additional test coverage and edge case testing

### Development Community

- **Code of Conduct**: Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md)
- **License**: This project is licensed under the MIT License
- **Attribution**: Please give appropriate credit when using this project

### Stay Updated

- **Watch the repository** for releases and announcements
- **Subscribe to releases** for version notifications
- **Follow the project** on GitHub for updates

## Performance Optimization

### Memory Management

1. **Use object pooling** for frequently allocated objects
2. **Avoid unnecessary allocations** in hot paths
3. **Use `sync.Pool`** for temporary objects
4. **Profile memory usage** regularly with `go tool pprof`

### Concurrency Patterns

1. **Use goroutines** for I/O-bound operations
2. **Use channels** for communication between goroutines
3. **Avoid shared memory** when possible
4. **Use `sync.WaitGroup`** for goroutine synchronization
5. **Implement proper cancellation** with `context.Context`

### Network Optimization

1. **Use connection pooling** for HTTP clients
2. **Implement request batching** where appropriate
3. **Cache responses** when possible
4. **Use timeouts** to prevent hanging operations
5. **Implement retry logic** for transient failures

## Security Guidelines

### Input Validation

1. **Validate all user inputs** before processing
2. **Sanitize file paths** to prevent directory traversal
3. **Validate API responses** for expected structure
4. **Use parameterized queries** for database operations
5. **Implement rate limiting** for API endpoints

### Secret Management

1. **Never hardcode API keys** or secrets
2. **Use environment variables** for sensitive data
3. **Set proper file permissions** for configuration files
4. **Use secret management tools** in production
5. **Implement audit logging** for sensitive operations

### Command Execution Security

1. **Validate commands** before execution
2. **Use safe execution methods** with proper sandboxing
3. **Implement proper error handling** to avoid information leakage
4. **Use platform-specific shell handling** (cmd.exe on Windows, sh on Unix-like systems)
5. **Implement timeout mechanisms** for long-running commands

## Release Process

### Version Management

1. **Semantic Versioning**: Follow semantic versioning (SemVer) principles
2. **Changelog**: Maintain a CHANGELOG.md file
3. **Tagging**: Use meaningful git tags for releases
4. **Branching**: Use feature branches for development

### Release Checklist

- [ ] All tests pass
- [ ] Documentation is updated
- [ ] Changelog is updated
- [ ] Version numbers are updated
- [ ] Release artifacts are built
- [ ] GitHub release is created
- [ ] Documentation is published

### Post-Release Activities

1. **Monitor feedback** from early adopters
2. **Address critical issues** promptly
3. **Plan next release** based on feedback and roadmap
4. **Update roadmap** based on completed features

---

Happy coding! 🚀
