# Slash Commands

Agent-Go supports a comprehensive set of slash commands for managing its configuration and features directly from the CLI. These commands provide quick access to common operations without interrupting your workflow.

## Command Overview

Slash commands are invoked by typing a forward slash (`/`) followed by the command name and any required arguments. They are executed immediately when you press Enter.

## General Commands

### `/help`

Displays a comprehensive list of all available commands with brief descriptions.

**Usage:**

```
/help
```

**Example:**

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
  /quit              - Exit the application
```

### `/config`

Shows the current configuration, including the model, provider URL, and RAG status. Useful for verifying your settings.

**Usage:**

```
/config
```

**Example:**

```
> /config
Model: gpt-4-turbo
Provider: https://api.openai.com
RAG Enabled: true
RAG Path: /home/user/documents
Auto Compress Enabled: true
Auto Compress Threshold: 20
Model Context Length: 131072
```

### `/shell`

Enters shell mode, allowing you to execute shell commands directly.

**Usage:**

```
/shell
```

**Example:**

```
> /shell
Entered shell mode. Type 'exit' to return.
shell> ls -l
-rw-r--r-- 1 user user 1024 Oct 27 10:00 file.txt
shell> exit
Exited shell mode.
```

**Notes:**

- In shell mode, all input is treated as a shell command.
- To exit shell mode, type `exit` and press Enter.
- Slash commands are not available while in shell mode.
- Commands are executed with platform-specific shell handling (cmd.exe on Windows, sh on Unix-like systems)

### `/quit`

Exits the Agent-Go application gracefully, saving any unsaved changes and cleaning up resources.

**Usage:**

```
/quit
```

**Example:**

```
> /quit
Bye!
```

## Model and Provider Commands

### `/model <model_name>`

Sets the AI model to be used for generating responses. The change is immediately reflected in subsequent conversations.

**Usage:**

```
/model <model_name>
```

**Parameters:**

- `model_name`: The name of the AI model (e.g., `gpt-4-turbo`, `gpt-3.5-turbo`, `claude-3-sonnet`)

**Examples:**

```
> /model gpt-4-turbo
Model set to: gpt-4-turbo

> /model claude-3-sonnet
Model set to: claude-3-sonnet
```

**Notes:**

- The model name must be supported by your API provider
- Changes are saved to the configuration file
- Auto-completion is available for model names (fetched from API)
- Model change affects all subsequent API requests

### `/provider <api_url>`

Sets the base URL for the API provider. This allows you to use different AI services or self-hosted endpoints.

**Usage:**

```
/provider <api_url>
```

**Parameters:**

- `api_url`: The base URL of the API provider (e.g., `https://api.openai.com`, `http://localhost:8080/v1`)

**Examples:**

```
> /provider https://api.openai.com
Provider URL set to: https://api.openai.com

> /provider http://localhost:8080/v1
Provider URL set to: http://localhost:8080/v1
```

**Notes:**

- The URL should point to the base endpoint of your API provider
- Changes are saved to the configuration file
- Ensure the provider is compatible with OpenAI's API format

## RAG (Retrieval-Augmented Generation) Commands

### `/rag on`

Enables the RAG (Retrieval-Augmented Generation) feature. When enabled, Agent-Go will search local documents for relevant context before generating responses.

**Usage:**

```
/rag on
```

**Example:**

```
> /rag on
RAG enabled.
```

**Notes:**

- Requires a valid RAG path to be set via `/rag path <path>`
- May increase response time as documents are searched
- Provides more context-aware responses for document-related queries
- Searches through subdirectories recursively
- Supports multiple file formats (txt, md, json, etc.)

### `/rag off`

Disables the RAG feature. When disabled, Agent-Go will not search local documents for context.

**Usage:**

```
/rag off
```

**Example:**

```
> /rag off
RAG disabled.
```

**Notes:**

- Responses may be less context-aware for document-related queries
- May improve response time as document search is skipped

### `/rag path <path>`

Sets the local file system path where documents for RAG are stored. This path is used when the RAG feature is enabled.

**Usage:**

```
/rag path <path>
```

**Parameters:**

- `path`: The absolute or relative path to your documents directory

**Examples:**

```
> /rag path /home/user/documents
RAG path set to: /home/user/documents

> /rag path ./docs
RAG path set to: ./docs
```

**Notes:**

- The path should contain text-based documents (txt, md, json, etc.)
- The directory must exist and be readable
- Changes are saved to the configuration file
- Subdirectories are also searched recursively
- File paths are validated for security
- Gracefully handles permission errors and inaccessible files

## Advanced Usage

### Command Chaining

While slash commands are typically executed individually, you can combine them with regular commands:

```
> /model gpt-4-turbo && /rag on
Model set to: gpt-4-turbo
RAG enabled.
```

### Error Handling

If you make a mistake in command syntax, Agent-Go will provide helpful error messages:

```
> /model
Usage: /model <model_name>

> /rag path
Usage: /rag path <path>
```

### Auto-completion

Agent-Go provides intelligent auto-completion for slash commands:

- `/` + Tab shows all available commands
- `/model` + Tab shows available models (fetched from API)
- `/rag` + Tab shows RAG options (`on`, `off`, `path`)
- `/provider` + Tab shows URL suggestions
- Dynamic model completion based on API response

### Context Management Commands

#### `/contextlength <value>`

Sets the model context length for token management. This determines when auto-compression triggers (at 75% of this value).

**Usage:**

```
/contextlength <value>
```

**Parameters:**

- `value`: Positive integer representing the context length (e.g., 131072 for GPT-4, 16384 for GPT-3.5)

**Examples:**

```
> /contextlength 16384
Model context length set to: 16384

> /contextlength 131072
Model context length set to: 131072
```

**Common Model Context Lengths:**
- GPT-4: 8192 or 131072 (depending on variant)
- GPT-3.5 Turbo: 16384
- Claude 3: 200000
- Local models: Varies (typically 2048-8192)

**Notes:**

- Must be a positive integer
- Affects auto-compression threshold calculation (triggers at 75% of this value)
- Changes are saved to the configuration file
- Should match your model's maximum context length
- Higher values allow longer conversations but use more tokens

### `/compress`

Compresses the current conversation context and starts a new chat with the compressed summary as a system message.

**Usage:**

```
/compress
```

**Example:**

```
> /compress
Context compressed. Starting new chat with compressed summary as system message.
```

**Notes:**

- Requires at least one message in the current chat
- Preserves key details and context from the previous conversation
- Completely clears the current message history
- Useful for very long conversations to avoid token limits
- Uses the same AI model for intelligent context compression
- Automatically includes AGENTS.md in the new system prompt if it exists
- Resets the token counter after compression
- Creates a fresh conversation thread while maintaining context

**When to Use:**

- When approaching token limits (check with `/config`)
- When conversation becomes too long and unwieldy
- When you want to start fresh but keep important context
- After completing a major task or project phase

## Configuration Persistence

All slash command changes (except `/help` and `/config`) are automatically saved to your configuration file at `~/.config/agent-go/config.json`. This ensures your preferences persist across sessions.

## Troubleshooting

### Common Issues

**Model not found:**

```
> /model invalid-model
Unknown model: invalid-model
```

Solution: Use `/help` or check your provider's documentation for valid model names.

**Invalid path:**

```
> /rag path /nonexistent/path
Error: cannot access /nonexistent/path
```

Solution: Ensure the directory exists and you have read permissions.

**API connection issues:**

```
> /model
Error: could not fetch models from API
```

Solution: Check your internet connection and API provider URL.

**Context compression errors:**

```
> /compress
Error: no messages to compress
```

Solution: Start a conversation first before attempting to compress context.

**Permission errors:**

```
> /rag path /protected/path
Error: permission denied
```

Solution: Choose a directory you have read access to.

### Getting Help

If you encounter issues not covered here:

1. Use `/help` to see all available commands
2. Check the main documentation in the `/docs` directory
3. Review the [Architecture](architecture.md) and [Configuration](configuration.md) documents
4. For bugs or feature requests, check the GitHub repository

### Debug Mode

For troubleshooting configuration issues:

```bash
# Enable debug logging
export DEBUG=1
./agent-go
```

This will provide detailed logging information to help diagnose issues.

## Additional Features

### Streaming Mode

Agent-Go supports streaming mode for real-time response generation:

```
> /stream on
Streaming enabled.

> Write a Python script that calculates Fibonacci numbers
[Streaming] Writing Python script...
[Streaming] Script created successfully...
```

**Usage:**
```
/stream on|off
```

**Notes:**
- When enabled, responses are streamed token by token for better user experience
- Reduces perceived latency for long responses
- Automatically disabled when shell mode is entered
- Can be toggled at any time during a session

### `/subagents on|off`

Toggles the ability for the main agent to spawn sub-agents for complex tasks.

**Usage:**

```
/subagents on|off
```

**Examples:**

```
> /subagents off
Sub-agent spawning disabled.

> /subagents on
Sub-agent spawning enabled.
```

**Notes:**

- When disabled, the `spawn_agent` tool is not available to the AI, preventing it from delegating tasks.
- This can be useful for forcing the primary agent to handle all tasks directly.
- The setting is saved to your configuration file.

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

## Troubleshooting

### Common Issues

**Model not found:**

```
> /model invalid-model
Unknown model: invalid-model
```

Solution: Use `/help` or check your provider's documentation for valid model names.

**Invalid path:**

```
> /rag path /nonexistent/path
Error: cannot access /nonexistent/path
```

Solution: Ensure the directory exists and you have read permissions.

**API connection issues:**

```
> /model
Error: could not fetch models from API
```

Solution: Check your internet connection and API provider URL.

**Context compression errors:**

```
> /compress
Error: no messages to compress
```

Solution: Start a conversation first before attempting to compress context.

**Permission errors:**

```
> /rag path /protected/path
Error: permission denied
```

Solution: Choose a directory you have read access to.

### Getting Help

If you encounter issues not covered here:

1. Use `/help` to see all available commands
2. Check the main documentation in the `/docs` directory
3. Review the [Architecture](architecture.md) and [Configuration](configuration.md) documents
4. For bugs or feature requests, check the GitHub repository

### Debug Mode

For troubleshooting configuration issues:

```bash
# Enable debug logging
export DEBUG=1
./agent-go
```

This will provide detailed logging information to help diagnose issues.

