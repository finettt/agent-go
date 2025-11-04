# Agent-Go
<img width="730" height="216" alt="image" src="https://github.com/user-attachments/assets/c2586899-0dcf-4544-9ed3-bb8a3c9a2bcc" />

A powerful AI agent written in Go that integrates with OpenAI-compatible APIs and provides intelligent command execution capabilities. This is a modern rewrite of the original [Agent-C](https://github.com/finettt/agent-c) project with enhanced features and improved architecture.

## Features

- **Tool Calling**: Execute shell commands directly through AI responses with intelligent error handling
- **Conversation Memory**: Maintains a sliding window of the last 20 messages for context-aware interactions
- **Cross-Platform**: Works seamlessly on macOS, Linux, and Windows
- **RAG (Retrieval-Augmented Generation)**: Searches local files to provide context-aware responses
- **Auto-completion**: Intelligent command-line autocompletion for models and commands
- **Slash Commands**: Built-in commands for configuration and feature management
- **Custom Instructions**: Support for AGENTS.md file for custom agent behavior
- **Interactive Setup**: First-time configuration wizard

## Quick Start

### Prerequisites

- Go 1.22+ (recommended)
- An OpenAI API key or access to a compatible API service

### Installation

#### Build from Source

```bash
# Clone the repository
git clone https://github.com/finettt/agent-go.git
cd agent-go

# Build the application
make build

# Or run directly
make run
```

#### Using Go

```bash
go install github.com/finettt/agent-go@latest
```

### Initial Setup

1. **Run the application**:
   ```bash
   ./agent-go
   ```

2. **First-time setup** (if no API key is configured):
   ```
   Agent-Go is ready. Type your requests, or /help for a list of commands.
   OpenAI API key is not set. Let's set it up.
   Enter your OpenAI API key: your_api_key_here
   ```

### Environment Variables

Configure Agent-Go using environment variables:

```bash
# Required: OpenAI API key
export OPENAI_KEY="your_openai_api_key_here"

# Optional: API provider URL (defaults to https://api.openai.com)
export OPENAI_BASE="https://api.openai.com"

# Optional: Default model (defaults to gpt-3.5-turbo)
export OPENAI_MODEL="gpt-4-turbo"

# Optional: RAG configuration
export RAG_PATH="/path/to/your/documents"  # Path to your documents
export RAG_ENABLED=1                       # Enable RAG (1=enabled, 0=disabled)
export RAG_SNIPPETS=5                      # Number of snippets to retrieve
```

## Usage Examples

### Basic Interaction

```
> Create a new directory called "test-project" and navigate into it
Created directory "test-project" and navigated into it.
```

### Multi-step Commands

```
> Create a Python script that prints "Hello, World!" and run it
Created script.py and executed it successfully.
Hello, World!
```

### Using RAG

```
> What does our documentation say about authentication?
Based on the provided documentation, authentication requires:
1. API key configuration
2. Proper endpoint setup
...
```

### Slash Commands

```
> /help
Available commands:
  /help              - Show this help message
  /model <name>      - Set the AI model (e.g., gpt-4)
  /provider <url>    - Set the API provider URL
  /config            - Display current configuration
  /rag on|off        - Toggle RAG feature
  /rag path <path>   - Set the RAG documents path
  /quit              - Exit the application

> /model gpt-4-turbo
Model set to: gpt-4-turbo
```

## Advanced Configuration

### Configuration File

Agent-Go automatically creates and manages a configuration file at `~/.config/agent-go/config.json`:

```json
{
  "api_url": "https://api.openai.com",
  "model": "gpt-4-turbo",
  "api_key": "your_secret_key_here",
  "rag_path": "/path/to/your/documents",
  "temp": 0.1,
  "max_tokens": 1000,
  "rag_enabled": true,
  "rag_snippets": 5
}
```

### Custom Agent Instructions

Create an `AGENTS.md` file in your project directory to provide custom instructions:

```markdown
# Agent Instructions

You are a helpful programming assistant. When users ask you to create files:
1. Always include proper error handling
2. Add comments explaining the code
3. Follow best practices for the language
```

## Architecture

Agent-Go is built with a clean, modular architecture:

- **main.go**: Application entry point and main CLI loop
- **config.go**: Configuration management with hierarchical settings
- **api.go**: OpenAI API communication and response handling
- **executor.go**: Secure shell command execution
- **rag.go**: Local file search and context retrieval
- **completion.go**: Auto-completion functionality
- **types.go**: Data structures and type definitions

## Documentation

For detailed documentation, see the `/docs` directory:

- [Architecture](docs/architecture.md) - System architecture and flow
- [Commands](docs/commands.md) - Complete command reference
- [Configuration](docs/configuration.md) - Configuration options
- [Development Guide](docs/development.md) - Contributing guidelines

## Contributing

We welcome contributions! Please see our [Development Guide](docs/development.md) for:

- Setting up your development environment
- Running tests
- Code style guidelines
- Pull request process

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [readline](https://github.com/chzyer/readline) for enhanced CLI experience
- Powered by OpenAI-compatible APIs
