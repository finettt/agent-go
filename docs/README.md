# Agent-Go Developer Documentation

Welcome to the developer documentation for Agent-Go. This documentation provides a comprehensive overview of the project's architecture, features, and development process.

## Getting Started

For quick start information, see the main [README.md](../README.md) file.

## Documentation Structure

### Core Documentation

- [Architecture](architecture.md) - System architecture, data flow, and component overview
- [Commands](commands.md) - Complete reference for all slash commands and CLI features
- [Configuration](configuration.md) - Configuration options, environment variables, and best practices

### Development Resources

- [Development Guide](development.md) - Contributing guidelines, testing, development workflow, and Docker deployment
- [Examples and Best Practices](examples.md) - Practical use cases, integration scenarios, and performance optimization

## Quick Navigation

### For Users
- **Basic Usage**: Start with the main [README.md](../README.md)
- **Commands**: Reference all available commands in [commands.md](commands.md)
- **Configuration**: Learn how to configure Agent-Go in [configuration.md](configuration.md)

### For Developers
- **Architecture**: Understand the system design in [architecture.md](architecture.md)
- **Development**: Follow the contributing guide in [development.md](development.md)
- **Docker Deployment**: Learn how to run Agent-Go in Docker containers in [development.md](development.md#using-docker)
- **Examples**: Explore practical implementations in [examples.md](examples.md)

## Key Features

### Core Functionality
- **AI Agent**: OpenAI-compatible API integration with tool calling
- **Command Execution**: Secure shell command execution with platform awareness
- **Context Management**: Unlimited conversation history with intelligent compression
- **RAG Integration**: Local document search and context retrieval
- **Auto-completion**: Dynamic model and command completion

### Advanced Features
- **Multi-step Commands**: Support for complex task chaining
- **Custom Instructions**: AGENTS.md file support for personalized behavior
- **Shell Mode**: Direct command execution interface
- **Token Tracking**: Real-time usage monitoring
- **Streaming Mode**: Real-time response generation
- **Cross-platform**: Works on macOS, Linux, and Windows

## Contributing

We welcome contributions! Please see our [Development Guide](development.md) for:

- Setting up your development environment
- Running tests
- Code style guidelines
- Pull request process

### Development Workflow

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Make your changes following the guidelines
4. Test your changes: `go test ./...`
5. Commit your changes: `git commit -m "feat: add new feature description"`
6. Push to your fork: `git push origin feature/your-feature-name`
7. Create a Pull Request with a clear description

## Support

If you encounter issues or have questions:

1. Check the main [README.md](../README.md) for basic usage
2. Review the [Commands](commands.md) for available features
3. Consult the [Configuration](configuration.md) for setup issues
4. Explore [Examples](examples.md) for practical use cases
5. Check the [Development Guide](development.md) for contributing

For bugs or feature requests, please open an issue on the GitHub repository.