<div align="center">
    <img width="300" height="300" alt="logo" src="https://github.com/user-attachments/assets/b8914faa-998d-487e-9173-7008d75b36df" />
    <h1>Agent-Go</h1>
</div>

**Standard Install (Pre-built Binary)**
```bash
curl -fsSL https://raw.githubusercontent.com/finettt/agent-go/main/install-agent-go.sh | bash
```

**Rolling Install (Build from Source)**
```bash
curl -fsSL https://raw.githubusercontent.com/finettt/agent-go/main/install-agent-go.sh | bash -s -- --rolling
```

<img width="1110" height="600" alt="image" src="https://github.com/user-attachments/assets/4775958a-5c5b-4184-8cc1-2f32aa693e86" />

Agent-Go is a powerful, command-line AI agent written in Go. It integrates with OpenAI-compatible APIs to provide intelligent command execution, context-aware responses, and a rich set of features for developers and power users. This project is a modern rewrite of the original [Agent-C](https://github.com/finettt/agent-c), redesigned with an improved architecture and enhanced capabilities.

## How It Works

Agent-Go works by interpreting your natural language requests, converting them into actionable commands, and executing them in a secure shell environment. It maintains a continuous conversation with an AI model, using the model's reasoning capabilities to break down complex tasks into smaller, executable steps.

The agent's workflow is as follows:

1. **Parse Input**: Your input is processed and, if RAG is enabled, augmented with relevant context from your local files.
2. **API Request**: The conversation history is sent to an OpenAI-compatible API.
3. **Tool-Assisted Response**: The AI model can use a set of built-in tools, including `execute_command`, `spawn_agent`, and `create_todo`, to perform actions.
4. **Command Execution**: If a command is returned, Agent-Go executes it and sends the output back to the model for the next step.
5. **Loop or Respond**: The agent continues this loop until the task is complete, then provides a final response.

## Key Features

| Feature | Description |
| :--- | :--- |
| **Agent Studio** | Complete agent management system for creating, managing, and using task-specific agents with the `/agent` command family. |
| **Session Management** | Save and restore conversation sessions with `/session` commands for seamless context switching and continuity. |
| **Background Command Execution** | Run shell commands in the background with support for monitoring, logging, and management. |
| **Intelligent Command Execution** | Natively executes shell commands from AI responses, with support for multi-step command chaining (`&&`). |
| **MCP (Model Context Protocol)** | Integrates with MCP servers to extend capabilities with external tools and resources. Includes default context7 integration for up-to-date library documentation. |
| **Sub-agent Delegation** | Spawns autonomous sub-agents to handle complex, multi-step tasks in the background, freeing up the main agent for other work. |
| **Unlimited Conversation Memory** | Overcomes token limits with automatic context compression, ensuring long-running conversations are never lost. |
| **Retrieval-Augmented Generation (RAG)** | Enriches AI responses by searching a local knowledge base, providing highly relevant and context-aware answers. |
| **Todo List Management** | Create, update, and manage a persistent todo list to track tasks and project status across sessions. |
| **Project Memory with Notes** | Create, update, and delete persistent notes that are automatically injected into the system prompt for context across sessions. |
| **Advanced User Experience** | Features a dedicated shell mode, dynamic command auto-completion, and a first-time interactive setup wizard. |
| **Streaming Mode** | Delivers real-time response generation for an improved, interactive user experience. |
| **Cross-Platform** | Compiles and runs seamlessly on macOS, Linux, and Windows. |
| **Highly Configurable** | Manage behavior via environment variables, a central JSON config, or command-line arguments. |
| **Custom Agent Behavior** | Define custom agent instructions and behavior using a simple `AGENTS.md` file. |
| **Real-time Token Tracking** | Monitor token usage for each conversation to manage API costs effectively. |

## Quick Start

### Prerequisites

- Go 1.25 or later
- An API key for an OpenAI-compatible service
- Git (optional, for cloning the source)
- Node.js/npm (optional, for MCP servers like context7)

### Installation

You can install Agent-Go either by building from the source, using Docker, or directly with `go install`.

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

**Using Docker**

```bash
# Clone the repository
git clone https://github.com/finettt/agent-go.git
cd agent-go

# Build the Docker image
docker build -t agent-go .

# Run with your current directory mounted as /workspace
docker run -it -v $(pwd):/workspace agent-go
```

**Using `go install`**

```bash
go install github.com/finettt/agent-go@latest
```

### Initial Setup

1. Run the application for the first time:

    ```bash
    agent-go
    ```

2. If your API key is not configured, the interactive setup wizard will prompt you:

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

### Agent Studio

Agent-Go can create and manage reusable, task-specific agents via the `/agent` command family.

Built-in agent:
- `default` (always available, cannot be deleted)

- `/agent studio [spec]`
  Starts an interactive "Agent Studio" chat. Describe the agent you want (goal, constraints, workflow, tooling rules).
  Studio will save it by calling [`create_agent_definition()`](src/agents.go:235).

- `/agent list`
  Lists saved agent definitions.

- `/agent view <name>`
  Prints the saved agent definition (including its system prompt).

- `/agent use <name>`
  Activates the agent for the current chat by rebuilding the system prompt (context is cleared).
  Optional overrides from the agent definition (model / temperature / max_tokens) are applied in-memory.

Subagents can also use task-specific agents: the `spawn_agent` tool supports an optional `agent` argument (e.g. `"default"` or a saved agent name).

- `/agent clear`
  Deactivates the active agent and restores the previous model settings snapshot (also clears context).

- `/agent rm <name>`
  Deletes a saved agent definition.

Saved agent JSON files are stored in `~/.config/agent-go/agents/*.json`.

Studio is intentionally restricted: it only permits agent creation and rejects all other tools. See [`processAgentStudioToolCalls()`](src/agent_studio.go:97).

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
    /subagents on|off  - Toggle sub-agent spawning
    /session new       - Create new session and save current context
    /session list      - View saved sessions
    /session restore <name> - Restore previous session
    /session rm <name> - Delete saved session
    /todo              - Display the current todo list
    /notes list        - List all notes
    /notes view <name> - View a specific note
    /agent studio      - Start Agent Studio for creating custom agents
    /agent list        - List saved agent definitions
    /agent view <name> - View a specific agent definition
    /agent use <name>  - Activate a specific agent
    /agent clear       - Clear agent-specific context
    /agent rm <name>   - Delete a saved agent definition
    /mcp add <name> <command> - Add an MCP server
    /mcp remove <name> - Remove an MCP server
    /mcp list          - List MCP servers
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

### MCP (Model Context Protocol) Integration

Agent-Go supports MCP servers to extend functionality with external tools and resources:

```
> /mcp list
Configured MCP servers:
- context7: npx -y @upstash/context7-mcp

> /mcp add time uvx mcp-server-time
MCP server 'time' added.

> Ask the agent to use MCP tools
> What's the current time in New York?
[Using MCP tool: get_current_time from server 'time']
Current time in New York: 10:30 AM EST
```

**Default MCP Server:**

- **context7**: Automatically configured for accessing up-to-date library documentation
- Provides tools: `resolve-library-id`, `get-library-docs`

**MCP Features:**

- Add custom MCP servers with `/mcp add`
- Remove servers with `/mcp remove`
- List configured servers with `/mcp list`
- AI automatically discovers and uses available tools from connected servers
- Server information is included in the system prompt

### Todo List Management

Agent-Go includes built-in todo list management for tracking tasks:

```
> Create a todo list for this project
Created todo list:
- [ID: 1] Set up development environment (pending)
- [ID: 2] Write documentation (pending)
- [ID: 3] Run tests (pending)

> Update todo 1 to in-progress
Updated todo list:
- [ID: 1] Set up development environment (in-progress)
- [ID: 2] Write documentation (pending)
- [ID: 3] Run tests (pending)

> /todo
Current Todo List:
- [ID: 1] Set up development environment (in-progress)
- [ID: 2] Write documentation (pending)
- [ID: 3] Run tests (pending)
```

**Todo Features:**

- Persistent storage in `~/.config/agent-go/todos/`
- Per-agent todo lists (main agent and sub-agents have separate lists)
- Status tracking: pending, in-progress, completed
- Create, update, and view todos via AI commands or `/todo` slash command

### Notes (Project Memory)

Agent-Go provides a notes feature for persistent project memory that carries across sessions:

```
> Ask the agent to remember something
> Remember that our API uses JWT authentication with RS256

Created note 'api_authentication':
The API uses JWT authentication with RS256 algorithm.

> /notes list
Notes:
  - api_authentication (updated: 2025-11-25 12:30)

> /notes view api_authentication
=== api_authentication ===
Created: 2025-11-25 12:30:00
Updated: 2025-11-25 12:30:00

The API uses JWT authentication with RS256 algorithm.
```

**Notes Features:**

- Persistent storage in `.agent-go/notes/` directory (per-project)
- Automatic injection into system prompt for cross-session context
- Agent tools: `create_note`, `update_note`, `delete_note`
- User commands: `/notes list`, `/notes view <name>`
- Tab autocomplete support for note names
- See [Notes Documentation](docs/notes.md) for complete details

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

1. **Environment Variables** (highest priority)
2. **Configuration File** (`~/.config/agent-go/config.json`)
3. **Default Values** (lowest priority)

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
  "model_context_length": 131072,
  "stream": false,
  "subagents_enabled": true,
  "mcp_servers": {
    "context7": {
      "name": "context7",
      "command": "npx -y @upstash/context7-mcp"
    }
  }
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
- **[api.go](src/api.go)**: Handles all communication with OpenAI-compatible APIs, including tool calling, streaming, and context compression.
- **[executor.go](src/executor.go)**: Manages secure, platform-aware shell command execution.
- **[processor.go](src/processor.go)**: Processes tool calls from API responses and coordinates tool execution.
- **[tools.go](src/tools.go)**: Defines available tools and their schemas for the AI agent.
- **[rag.go](src/rag.go)**: Implements local document search and context retrieval.
- **[mcp.go](src/mcp.go)**: Manages MCP (Model Context Protocol) client connections and tool calls.
- **[subagent.go](src/subagent.go)**: Manages the lifecycle and execution of sub-agents (max 50 iterations per sub-agent).
- **[todo.go](src/todo.go)**: Handles todo list creation, updates, and persistence.
- **[notes.go](src/notes.go)**: Manages persistent notes for project memory across sessions.
- **[completion.go](src/completion.go)**: Provides dynamic auto-completion for models and commands.
- **[system.go](src/system.go)**: Gathers system information for context.
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
- [Notes Feature](docs/notes.md) - Persistent project memory documentation

## Contributing

We welcome contributions! Please see our [Development Guide](docs/development.md) for details on setting up your development environment, running tests, and our pull request process.

### Development Workflow

1. Fork the repository.
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Make and test your changes: `go test ./...`
4. Commit your changes with a descriptive message: `git commit -m "feat: add new feature description"`
5. Push to your fork: `git push origin feature/your-feature-name`
6. Create a Pull Request.

## License

This project is licensed under the Apache 2.0 with Commons Clause License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [readline](https://github.com/chzyer/readline) for an enhanced CLI experience
- MCP integration via [go-sdk](https://github.com/modelcontextprotocol/go-sdk)
- Powered by OpenAI-compatible APIs
- Default MCP server: [context7](https://github.com/upstash/context7-mcp) by Upstash
