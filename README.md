# Agent-Go
<img width="730" height="216" alt="image" src="https://github.com/user-attachments/assets/c2586899-0dcf-4544-9ed3-bb8a3c9a2bcc" />

A powerful AI agent written in Go that integrates with OpenAI-compatible APIs and provides intelligent command execution capabilities. This is a modern rewrite of the original [Agent-C](https://github.com/finettt/agent-c) project with enhanced features and improved architecture.

## Features

- **Tool Calling**: Execute shell commands directly through AI responses with intelligent error handling
- **Unlimited Conversation Memory**: No message limits with intelligent context compression
- **Cross-Platform**: Works seamlessly on macOS, Linux, and Windows
- **RAG (Retrieval-Augmented Generation)**: Searches local files to provide context-aware responses
- **Auto-completion**: Intelligent command-line autocompletion for models and commands
- **Slash Commands**: Built-in commands for configuration and feature management
- **Custom Instructions**: Support for AGENTS.md file for custom agent behavior
- **Interactive Setup**: First-time configuration wizard
- **Context Compression**: Automatically compresses long conversations to avoid token limits
- **Shell Mode**: Direct command execution mode for interactive shell sessions
- **Auto-compress Threshold**: Configurable context compression when approaching token limits
- **Model Context Length**: Configurable context length for different AI models
- **Command Line Task Execution**: Execute single tasks directly from command line
- **Token Tracking**: Real-time token usage monitoring
- **Multi-step Command Chaining**: Support for && operators in complex tasks
- **Secure Command Execution**: Platform-aware command execution with proper validation

## Quick Start

### Prerequisites

- Go 1.25 (recommended, currently tested with 1.25.3)
- An OpenAI-compatible API key or access to a compatible API service
- Git (for version control, optional)

### Installation

#### Build from Source

```bash
# Clone the repository
git clone https://github.com/finettt/agent-go.git
cd agent-go

# Build the application
make build

# Run the application (this builds first)
make run

# The binary is created as ./agent-go in the project root
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

3. **Command line task execution** (alternative to interactive mode):
   ```bash
   ./agent-go "Create a new directory called 'test-project' and navigate into it"
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

# Optional: Context management
export AUTO_COMPRESS=1                     # Enable auto context compression (1=enabled, 0=disabled)
export AUTO_COMPRESS_THRESHOLD=20          # Threshold for auto compression
export MODEL_CONTEXT_LENGTH=131072         # Model context length (e.g., 131072 for GPT-4)
```

## Examples

### Basic Interaction

```
> Create a new directory called "test-project" and navigate into it
$ mkdir test-project && cd test-project
Created directory "test-project" and navigated into it.
```

### Multi-step Commands

```
> Create a Python script that prints "Hello, World!" and run it
$ echo 'print("Hello, World!")' > hello.py && python hello.py
Created hello.py and executed successfully.
Hello, World!
```

### Using RAG (Retrieval-Augmented Generation)

```
> What does our documentation say about authentication?
Based on the provided documentation, authentication requires:
1. API key configuration
2. Proper endpoint setup
...

# Enable RAG with custom document path
> /rag on
RAG enabled.

> /rag path ./docs
RAG path set to: ./docs
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
  /shell             - Enter shell mode for direct command execution
  /compress          - Compress context and start new chat thread
  /contextlength <value> - Set the model context length
  /quit              - Exit the application

> /model gpt-4-turbo
Model set to: gpt-4-turbo
```

### Command Line Task Execution

```
$ ./agent-go "Create a new directory called 'test-project' and navigate into it"
$ mkdir test-project && cd test-project
Created directory "test-project" and navigated into it.
```

### Shell Mode

```
> /shell
Entered shell mode. Type 'exit' to return.
shell> ls -la
total 8
drwxr-xr-x 2 user user 4096 Oct 27 10:00 .
drwxr-xr-x 5 user user 4096 Oct 27 10:00 ..
shell> exit
Exited shell mode.
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
  "rag_snippets": 5,
  "auto_compress": true,
  "auto_compress_threshold": 20,
  "model_context_length": 131072
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

# Optional: Context management
export AUTO_COMPRESS=1                     # Enable auto context compression (1=enabled, 0=disabled)
export AUTO_COMPRESS_THRESHOLD=20          # Threshold for auto compression (percentage)
export MODEL_CONTEXT_LENGTH=131072         # Model context length (e.g., 131072 for GPT-4)
```

## Architecture

Agent-Go is built with a clean, modular architecture:

- **[main.go](src/main.go)**: Application entry point, CLI loop, and command orchestration
- **[config.go](src/config.go)**: Hierarchical configuration management (env vars → config file → defaults)
- **[api.go](src/api.go)**: OpenAI-compatible API communication, tool calling, and context compression
- **[executor.go](src/executor.go)**: Secure, platform-aware shell command execution
- **[rag.go](src/rag.go)**: Local document search and context retrieval functionality
- **[completion.go](src/completion.go)**: Dynamic auto-completion for models and commands
- **[types.go](src/types.go)**: Shared data structures and type definitions

### Key Components

- **Agent Loop**: Multi-turn conversation handling with intelligent tool execution
- **Context Management**: Unlimited history with AI-powered compression at configurable thresholds
- **Security**: Platform-aware command execution (cmd.exe on Windows, sh on Unix-like systems)
- **Auto-completion**: Real-time model fetching from API with graceful fallbacks
- **Error Handling**: Comprehensive error handling with user-friendly messages
- **Token Tracking**: Real-time cumulative token usage monitoring and display

## Documentation

For detailed documentation, see the `/docs` directory:

- [Architecture](docs/architecture.md) - System architecture and flow
- [Commands](docs/commands.md) - Complete command reference
- [Configuration](docs/configuration.md) - Configuration options
- [Development Guide](docs/development.md) - Contributing guidelines
- [Examples and Best Practices](docs/examples.md) - Practical use cases and examples

## Contributing

We welcome contributions! Please see our [Development Guide](docs/development.md) for:

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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [readline](https://github.com/chzyer/readline) for enhanced CLI experience
- Powered by OpenAI-compatible APIs

### Token Usage

Agent-Go tracks and displays token usage in real-time:
- Shows total tokens consumed in the conversation
- Automatically compresses context when approaching token limits
- Configurable context length for different AI models
