# Agent-Go
<img width="730" height="216" alt="image" src="https://github.com/user-attachments/assets/c2586899-0dcf-4544-9ed3-bb8a3c9a2bcc" />

Agent-Go is a powerful, command-line AI agent written in Go. It integrates with OpenAI-compatible APIs to provide intelligent command execution, context-aware responses, and a rich set of features for developers and power users. This project is a modern rewrite of the original [Agent-C](https://github.com/finettt/agent-c), redesigned with an improved architecture and enhanced capabilities.

## Key Features

- **Intelligent Command Execution**: Natively executes shell commands from AI responses, with support for multi-step command chaining (`&&`) and platform-aware security.
- **Unlimited Conversation Memory**: Overcomes message limits through intelligent, automatic context compression, ensuring long-running conversations stay within token limits.
- **Retrieval-Augmented Generation (RAG)**: Enhances AI responses by searching a local knowledge base of your files, providing highly relevant and context-aware answers.
- **Advanced User Experience**: Features a dedicated shell mode, dynamic command-line auto-completion, built-in slash commands, and a first-time interactive setup wizard.
- **Streaming Mode**: Real-time response generation for improved user experience.
- **Cross-Platform**: Compiles and runs seamlessly on macOS, Linux, and Windows.
- **Highly Configurable**: Manage behavior via environment variables, a central JSON configuration file, or command-line arguments.
- **Custom Agent Behavior**: Define custom agent instructions and behavior using a simple `AGENTS.md` file.
- **Real-time Token Tracking**: Monitor token usage for each conversation to manage API costs effectively.

## Quick Start

### Prerequisites

- Go 1.25 or later (tested with 1.25.3)
- An API key for an OpenAI-compatible service
- Git (optional, for cloning the source)

### Installation

You can install Agent-Go either by building from the source or directly with `go install`.

**Build from Source**
```bash
# Clone the repository
git clone https://github.com/finettt/agent-go.git
cd agent-go

# Build the binary (creates ./agent-go)
make build

# Or, build and run the application
make run
```

**Using `go install`**
```bash
go install github.com/finettt/agent-go@latest
```

### Initial Setup

1.  Run the application for the first time:
    ```bash
    agent-go
    ```

2.  If your API key is not configured, the interactive setup wizard will prompt you:
    ```
    Agent-Go is ready. Type your requests, or /help for a list of commands.
    OpenAI API key is not set. Let's set it up.
    Enter your OpenAI API key: your_api_key_here
    ```

You are now ready to interact with the agent.

## Usage

Agent-Go can be used in an interactive chat mode or to execute single tasks directly from the command line.

### Interactive Mode
Start the agent without any arguments to enter the interactive chat loop.

**Basic Interaction**
```
> Create a new directory called "test-project" and navigate into it
$ mkdir test-project && cd test-project
Created directory "test-project" and navigated into it.
```

**Multi-step Commands**
```
> Create a Python script that prints "Hello, World!" and then run it
$ echo 'print("Hello, World!")' > hello.py && python hello.py
Created hello.py and executed successfully.
Hello, World!
```

### Single Task Execution
Pass a task as a command-line argument for non-interactive execution.
```bash
agent-go "Create a new directory called 'test-project' and navigate into it"
$ mkdir test-project && cd test-project
Created directory "test-project" and navigated into it.
```

### Slash Commands
Manage Agent-Go's configuration and features from within the application using slash commands. Type `/help` to see the full list.

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
  /stream on|off     - Toggle streaming mode
  /quit              - Exit the application

> /model gpt-4-turbo
Model set to: gpt-4-turbo
```

### Streaming Mode
Enable streaming mode for real-time response generation:

```
> /stream on
Streaming enabled.

> Write a Python script that calculates Fibonacci numbers
[Streaming] Writing Python script...
[Streaming] Script created successfully...
```

**Benefits:**
- Reduces perceived latency for long responses
- Provides immediate feedback during generation
- Can be toggled at any time during a session

### Shell Mode
For a more traditional shell experience, use the `/shell` command. The agent will not interpret your input, and commands are executed directly.

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

### Using RAG (Retrieval-Augmented Generation)
Enable RAG to allow the agent to search local files for context before answering.

```
# Enable RAG and set the path to your documents
> /rag on
RAG enabled.
> /rag path ./docs
RAG path set to: ./docs

# Ask a question related to your documents
> What does our documentation say about authentication?
Based on the provided documentation, authentication requires:
1. API key configuration
2. Proper endpoint setup
...
```

### Token Usage Tracking
Agent-Go tracks and displays token usage in real-time:

```
> Create a Python script that calculates factorial
[Tokens: 156]

> Run the factorial script with input 10
[Tokens: 342]

> /compress
Context compressed. Starting new chat with compressed summary as system message.
[Tokens: 0]  # Token counter resets after compression
```

**Understanding Token Usage:**
- Tokens are cumulative throughout the session
- Each API call adds to the total token count
- Auto-compression triggers at 75% of `model_context_length`
- Manual compression with `/compress` resets the token counter

### Command-Line Task Execution
Execute tasks directly without interactive mode:

```bash
# Single task execution
./agent-go "Create a new directory called 'test-project' and initialize git"

# The agent will execute the task and exit automatically
```

**Benefits:**
- Ideal for scripting and automation
- No interactive prompts
- Returns exit code 0 on success, non-zero on failure
- Perfect for CI/CD pipelines

## Configuration

Agent-Go uses a hierarchical configuration system. Settings are applied in the following order of precedence:
1.  **Environment Variables** (highest priority)
2.  **Configuration File** (`~/.config/agent-go/config.json`)
3.  **Default Values** (lowest priority)

### Environment Variables
Configure the agent by setting the following environment variables.

```bash
# Required: Your API key
export OPENAI_KEY="your_openai_api_key_here"

# Optional: API provider URL
# Defaults to https://api.openai.com
export OPENAI_BASE="https://api.openai.com"

# Optional: Default model
# Defaults to gpt-3.5-turbo
export OPENAI_MODEL="gpt-4-turbo"

# Optional: RAG configuration
export RAG_PATH="/path/to/your/documents"
export RAG_ENABLED=1                       # (1=enabled, 0=disabled)
export RAG_SNIPPETS=5                      # Number of snippets to retrieve

# Optional: Context management
export AUTO_COMPRESS=1                     # (1=enabled, 0=disabled)
export AUTO_COMPRESS_THRESHOLD=20          # Threshold (percentage) to trigger compression
export MODEL_CONTEXT_LENGTH=131072         # Max context tokens for the model
```

### Configuration File
On its first run, Agent-Go creates a configuration file at `~/.config/agent-go/config.json`. You can edit this file directly for persistent settings.

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
To provide the agent with custom instructions, create an `AGENTS.md` file in the directory where you run Agent-Go. The content of this file will be included as a system message in every request, guiding the agent's behavior and responses.

**Example `AGENTS.md`:**
```markdown
# Agent Instructions

You are a helpful programming assistant. When users ask you to create files:
1. Always include proper error handling.
2. Add comments explaining the code.
3. Follow best practices for the language.
```

## Architecture

Agent-Go is built with a clean, modular architecture for maintainability and extensibility.

- **[main.go](src/main.go)**: Application entry point, CLI loop, and command orchestration.
- **[config.go](src/config.go)**: Hierarchical configuration management.
- **[api.go](src/api.go)**: Handles all communication with OpenAI-compatible APIs, including tool calling and context compression.
- **[executor.go](src/executor.go)**: Manages secure, platform-aware shell command execution.
- **[rag.go](src/rag.go)**: Implements local document search and context retrieval.
- **[completion.go](src/completion.go)**: Provides dynamic auto-completion for models and commands.
- **[types.go](src/types.go)**: Defines shared data structures and type definitions.

### Key Architectural Components

- **Agent Loop**: Manages multi-turn conversations and orchestrates tool execution.
- **Context Management**: Maintains an unlimited conversation history by using AI-powered summarization to compress context when a configurable token threshold is reached.
- **Security**: Executes commands in a platform-aware manner (`cmd.exe` on Windows, `sh` on Unix-like systems) with proper validation.
- **Error Handling**: Provides comprehensive error handling with user-friendly messages.

### Documentation

For a deeper dive into the project's design and features, please see the `/docs` directory:

- [Architecture](docs/architecture.md) - System architecture and flow
- [Commands](docs/commands.md) - Complete command reference
- [Configuration](docs/configuration.md) - Configuration options
- [Development Guide](docs/development.md) - Contributing guidelines
- [Examples and Best Practices](docs/examples.md) - Practical use cases and examples

## Contributing

We welcome contributions! Please see our [Development Guide](docs/development.md) for details on setting up your development environment, running tests, and our pull request process.

### Development Workflow

1.  Fork the repository.
2.  Create a feature branch: `git checkout -b feature/your-feature-name`
3.  Make and test your changes: `go test ./...`
4.  Commit your changes with a descriptive message: `git commit -m "feat: add new feature description"`
5.  Push to your fork: `git push origin feature/your-feature-name`
6.  Create a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [readline](https://github.com/chzyer/readline) for an enhanced CLI experience
- Powered by OpenAI-compatible APIs
