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
  /session new       - Create new session and save current context
  /session list      - View saved sessions
  /session restore <name> - Restore previous session
  /session rm <name> - Delete saved session
  /agent studio      - Start Agent Studio for creating custom agents
  /agent list        - List saved agent definitions
  /agent view <name> - View a specific agent definition
  /agent use <name>  - Activate a specific agent
  /agent clear       - Clear agent-specific context
  /agent rm <name>   - Delete a saved agent definition
  /todo              - Display the current todo list
  /notes list        - List all notes
  /notes view <name> - View a specific note
  /mcp add <name> <command> - Add an MCP server
  /mcp remove <name> - Remove an MCP server
  /mcp list          - List MCP servers
  /usage             - Display detailed token usage statistics
  /cost              - Display cost tracking information
  /verbose on|off    - Toggle verbose logging mode
  /security          - Spawn subagent for security review
  /edit              - Open nano editor to compose prompt
  /quit              - Exit the application
 ````

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

### `/clear`

Clears the current conversation context (message history) without creating a new session or compressing context.

**Usage:**

```
/clear
```

**Example:**

```
> /clear
Context cleared.
```

### `/sandbox`

Relaunches Agent-Go in a Docker sandbox environment for isolated execution.

**Usage:**

```
/sandbox
```

**Example:**

```
> /sandbox
Building Docker image...
Starting sandbox environment...
```

**Notes:**

- Requires Docker to be installed and running
- Mounts the current working directory and configuration
- Useful for executing potentially unsafe commands in isolation

### `/bg` (Background Commands)

Manage background processes directly from the agent interface.

- **`/bg list`**: List all running background processes
- **`/bg view <pid>`**: View the logs (stdout/stderr) for a specific process ID
- **`/bg kill <pid>`**: Terminate a background process

**Usage:**

```
/bg list
/bg view 12345
/bg kill 12345
```

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

## Session Management Commands

### `/session new`

Creates a new session by saving the current conversation context. This allows you to return to your current state later or switch between different projects.

**Usage:**

```
/session new
```

**Example:**

```
> /session new
Session 'session-20251218-123456' saved successfully.
```

**Notes:**

- Automatically saves the current conversation history
- Generates a unique session name with timestamp
- Session data includes compressed context and metadata
- Sessions are stored in `~/.config/agent-go/sessions/`
- Useful for project switching and context preservation

### `/session list`

Lists all saved sessions with metadata including creation time, last access time, and conversation statistics.

**Usage:**

```
/session list
```

**Example:**

```
> /session list
Saved Sessions:
1. session-20251218-123456 (Created: 2025-12-18 10:30:00, Last Accessed: 2025-12-18 10:45:00, Messages: 25, Tokens: 15,432)
2. project-alpha (Created: 2025-12-15 09:15:00, Last Accessed: 2025-12-16 16:20:00, Messages: 150, Tokens: 89,215)
3. api-integration (Created: 2025-12-10 14:45:00, Last Accessed: 2025-12-12 11:30:00, Messages: 89, Tokens: 52,178)

Current Session: session-20251218-123456
```

**Notes:**

- Displays creation time, last access time, message count, and token usage
- Shows the current active session
- Sessions are ordered by last access time
- Provides quick overview of available sessions

### `/session restore <name>`

Restores a previously saved session, loading its conversation history and context.

**Usage:**

```
/session restore <name>
```

**Parameters:**

- `name`: The name of the session to restore

**Example:**

```
> /session restore project-alpha
Session 'project-alpha' restored successfully.
Loaded 150 messages with 89,215 tokens.
```

**Notes:**

- Replaces current conversation with the restored session
- Compressed context is automatically decompressed
- All session metadata is preserved
- Current session is saved before restoration
- Useful for continuing previous work

### `/session rm <name>`

Deletes a saved session permanently.

**Usage:**

```
/session rm <name>
```

**Parameters:**

- `name`: The name of the session to delete

**Example:**

```
> /session rm old-session
Session 'old-session' deleted successfully.
```

**Notes:**

- Permanently removes the session file
- Cannot be undone
- Useful for cleaning up old or unused sessions
- Validates session existence before deletion

## Agent Studio Commands

Agent Studio is a complete agent management system that allows you to create, manage, and use task-specific agents.

### `/agent studio [spec]`

Starts the Agent Studio interface for creating custom agents. You can optionally provide an initial specification.

**Usage:**

```
/agent studio
/agent studio [spec]
```

**Parameters:**

- `spec` (optional): Initial agent specification

**Example:**

```
> /agent studio
Welcome to Agent Studio!

Describe the agent you want to create:
1. What is the agent's primary goal or purpose?
2. What are its constraints or limitations?
3. What is its workflow or process?
4. What tools or capabilities should it have?

Enter your agent description, or type 'help' for guidance.
```

**Agent Studio Features:**

- Interactive chat interface for agent creation
- Validates agent specifications
- Only permits agent creation (rejects all other tools)
- Automatically saves agent definitions
- Protected `default` agent that cannot be deleted

**Example Agent Creation:**

```
> /agent studio

User: I want to create a code review agent that analyzes pull requests

Agent Studio: Great! Let's create a code review agent. Please specify:
1. Primary goal: Analyze code changes for quality, security, and best practices
2. Constraints: Should not approve changes that break existing functionality
3. Workflow: Review diffs, run linting, check for security issues, provide feedback
4. Tools: execute_command, spawn_agent, get_todo_list, create_todo, update_todo

Please confirm this specification or make adjustments.
```

### `/agent list`

Lists all saved agent definitions with their metadata.

**Usage:**

```
/agent list
```

**Example:**

```
> /agent list
Available Agents:

1. default (Built-in)
   - Purpose: General-purpose AI assistant
   - Model: gpt-4-turbo
   - Temperature: 0.1
   - Max Tokens: 1000
   - Status: Protected (cannot be deleted)

2. code-reviewer
   - Purpose: Analyze code changes for quality and security
   - Model: gpt-4-turbo
   - Temperature: 0.2
   - Max Tokens: 2000
   - Created: 2025-12-18 10:30:00
   - Updated: 2025-12-18 10:30:00

3. documentation-writer
   - Purpose: Generate technical documentation
   - Model: gpt-4-turbo
   - Temperature: 0.7
   - Max Tokens: 1500
   - Created: 2025-12-17 14:20:00
   - Updated: 2025-12-17 14:20:00
```

### `/agent view <name>`

Displays the detailed definition of a specific agent, including its system prompt and configuration.

**Usage:**

```
/agent view <name>
```

**Parameters:**

- `name`: The name of the agent to view

**Example:**

```
> /agent view code-reviewer
=== Agent: code-reviewer ===
Purpose: Analyze code changes for quality and security
Model: gpt-4-turbo
Temperature: 0.2
Max Tokens: 2000
Created: 2025-12-18 10:30:00
Updated: 2025-12-18 10:30:00

System Prompt:
You are a code review agent specialized in analyzing pull requests. Your purpose is to:
1. Analyze code changes for quality, security, and best practices
2. Review diffs and identify potential issues
3. Run linting and static analysis tools
4. Check for security vulnerabilities
5. Provide constructive feedback to developers

When reviewing code:
- Be thorough but constructive
- Focus on logic, performance, and security
- Suggest improvements and alternatives
- Never approve changes that could break existing functionality
- Use the available tools to assist in your analysis

Available Tools:
- execute_command: Run shell commands for testing and analysis
- spawn_agent: Delegate complex tasks to specialized sub-agents
- get_todo_list: Check current todo items
- create_todo: Add new todo items
- update_todo: Update existing todo items
- use_mcp_tool: Access MCP server tools
```

### `/agent use <name>`

Activates a specific agent for the current chat session. This rebuilds the system prompt with the agent's configuration and clears the current context.

**Usage:**

```
/agent use <name>
```

**Parameters:**

- `name`: The name of the agent to activate

**Example:**

```
> /agent use code-reviewer
Agent 'code-reviewer' activated.
System prompt rebuilt with agent configuration.
Context cleared for focused agent operation.
```

**Notes:**

- Clears current conversation context
- Rebuilds system prompt with agent's configuration
- Applies agent-specific model settings (if any)
- Context is isolated from previous agent usage
- Useful for switching between different tasks or projects

### `/agent clear`

Deactivates the current agent and restores the previous model settings, clearing the context.

**Usage:**

```
/agent clear
```

**Example:**

```
> /agent clear
Agent deactivated.
Previous model settings restored.
Context cleared.
```

**Notes:**

- Returns to the default agent behavior
- Restores previous model configuration
- Clears conversation context
- Useful for resetting the agent state

### `/agent rm <name>`

Deletes a saved agent definition.

**Usage:**

```
/agent rm <name>
```

**Parameters:**

- `name`: The name of the agent to delete

**Example:**

```
> /agent rm test-agent
Agent 'test-agent' deleted successfully.
```

**Notes:**

- Permanently removes the agent definition
- Cannot delete the built-in `default` agent
- Validates agent existence before deletion
- Useful for cleaning up unused agents

## Usage Tracking Commands

### `/usage`

Displays detailed token usage statistics for the current session, including cumulative counts and breakdowns.

**Usage:**

```
/usage
```

**Example:**

```
> /usage
=== Token Usage Statistics ===

Current Session:
- Total Tokens: 15,432
  * Prompt Tokens: 8,241
  * Completion Tokens: 7,191
  * Tokens (Reasoning): 0

Cumulative Totals:
- Total: 15,432 tokens
- Estimated Cost: $0.03 (based on current pricing)

Session Information:
- Start Time: 2025-12-18 10:30:00
- Current Session: session-20251218-123456
- Active Agent: default
```

**Notes:**

- Shows both prompt and completion token counts
- Displays tokens used for reasoning (if applicable)
- Provides cost estimation based on current model pricing
- Cumulative counts include all API calls in the session
- Resets after context compression

### `/cost`

Displays detailed cost tracking information, including current session costs and historical data.

**Usage:**

```
/cost
```

**Example:**

```
> /cost
=== Cost Tracking Information ===

Current Session:
- API Cost: $0.03
- Tokens Used: 15,432
- Model: gpt-4-turbo

Historical Costs:
- Today: $0.45
- This Week: $2.15
- This Month: $15.80

Cost Breakdown:
- Chat Completion: $0.03 (100%)

Session: session-20251218-123456
```

**Notes:**

- Tracks costs across sessions and time periods
- Provides historical cost data for budgeting
- Breaks down costs by operation type
- Useful for monitoring API usage expenses

### `/verbose on|off`

Toggles verbose logging mode for enhanced debugging and monitoring.

**Usage:**

```
/verbose on
/verbose off
```

**Example:**

```
> /verbose on
Verbose mode enabled.
Detailed logging active for agent operations.

> /verbose off
Verbose mode disabled.
Reduced logging output.
```

**Notes:**

- Enables detailed logging for debugging
- Shows internal agent operations and decision-making
- Useful for troubleshooting complex issues
- Can be configured via environment variables
- Affects subagent and tool execution logging

### `/security`

Spawns a specialized subagent for security code review and analysis.

**Usage:**

```
/security
```

**Example:**

```
> /security
Security subagent spawned.
Analyzing code for potential security vulnerabilities...

Security Review Results:
- No critical vulnerabilities detected
- 2 potential security issues identified
- Recommendations provided for secure coding practices

Security analysis complete.
```

**Notes:**

- Creates a specialized subagent for security analysis
- Reviews code for common vulnerabilities
- Provides security recommendations
- Useful for pre-deployment code review
- Integrates with existing security tools via MCP

### `/edit`

Opens the nano text editor to compose multi-line prompts or commands.

**Usage:**

```
/edit
```

**Example:**

```
> /edit
Opening nano editor...

[User edits prompt in nano editor...]

Prompt saved and ready for execution.
```

**Notes:**

- Opens nano text editor for complex input
- Useful for multi-line prompts or commands
- Automatically saves and returns to Agent-Go
- Requires nano to be installed on the system
- Alternative to direct command line input for complex tasks

### Background Command Tools

Agent-Go supports background command execution through specialized tools:

**`execute_command` with `background` parameter:**

```
{
  "command": "long-running-task.sh",
  "background": true
}
```

**`get_background_logs`:**
Retrieves output from running background processes.

**`list_background_commands`:**
Lists all currently running background commands.

**`kill_background_command`:**
Terminates a specific background command.

**Notes:**

- Background commands run asynchronously
- Process IDs are tracked for management
- Output is streamed in real-time
- Application prevents exit while background tasks run
- Useful for long-running operations

## MCP (Model Context Protocol) Commands

### `/mcp add <name> <command>`

Adds a new MCP (Model Context Protocol) server to extend Agent-Go's capabilities with external tools and resources.

**Usage:**

```
/mcp add <name> <command>
```

**Parameters:**

- `name`: Unique identifier for the MCP server
- `command`: Command to launch the MCP server

**Examples:**

```
> /mcp add time uvx mcp-server-time
MCP server 'time' added.

> /mcp add weather npx -y @weather/mcp-server
MCP server 'weather' added.

> /mcp add filesystem npx -y @modelcontextprotocol/server-filesystem /path/to/files
MCP server 'filesystem' added.
```

**Notes:**

- The server is immediately available for use by the AI agent
- Server configuration is persisted to the config file
- The AI automatically discovers available tools from connected servers
- Common MCP server commands use `npx -y` or `uvx` for package execution

### `/mcp remove <name>`

Removes a configured MCP server.

**Usage:**

```
/mcp remove <name>
```

**Parameters:**

- `name`: Name of the MCP server to remove

**Examples:**

```
> /mcp remove weather
MCP server 'weather' removed.

> /mcp remove old-server
MCP server 'old-server' removed.
```

**Notes:**

- Removes the server from configuration
- Changes are saved immediately
- The server will no longer be available for tool calls

### `/mcp list`

Lists all configured MCP servers and their commands.

**Usage:**

```
/mcp list
```

**Example:**

```
> /mcp list
Configured MCP servers:
- context7: npx -y @upstash/context7-mcp
- time: uvx mcp-server-time
- weather: npx -y @weather/mcp-server
```

**Notes:**

- Shows all configured servers
- Displays the command used to launch each server
- The `context7` server is configured by default for library documentation

**Default MCP Server:**

Agent-Go comes pre-configured with the **context7** MCP server for accessing up-to-date library documentation:

```
> Ask the agent to get React documentation
[Using MCP: Connecting to context7 server]
[Fetching documentation for React...]
Based on the latest React documentation...
```

**Available Tools from context7:**

- `resolve-library-id`: Finds the correct library identifier
- `get-library-docs`: Retrieves up-to-date documentation for a library

### `/todo`

Displays the current todo list for the active agent.

**Usage:**

```
/todo
```

**Example:**

```
> /todo
Current Todo List:
- [ID: 1] Set up development environment (in-progress)
- [ID: 2] Write documentation (pending)
- [ID: 3] Run tests (completed)
```

**Notes:**

- Todo lists are persistent across sessions
- Each agent (main and sub-agents) has its own todo list
- Todo lists are stored in `~/.config/agent-go/todos/`
- You can also ask the AI to create, update, or view todos

**Todo Management via AI:**

```
> Create a todo for setting up the database
Created new todo:
- [ID: 4] Set up the database (pending)

> Update todo 4 to in-progress
Updated todo:
- [ID: 4] Set up the database (in-progress)

> Mark todo 4 as completed
Updated todo:
- [ID: 4] Set up the database (completed)
```

**Todo Statuses:**

- `pending`: Not yet started
- `in-progress`: Currently being worked on
- `completed`: Finished

## Notes Management Commands

### `/notes list`

Lists all existing notes with their last update time.

**Usage:**

```
/notes list
```

**Example:**

```
> /notes list
Notes:
  - api_endpoint (updated: 2025-11-25 08:30)
  - database_schema (updated: 2025-11-25 09:15)
  - deployment_instructions (updated: 2025-11-25 10:45)
```

**Notes:**

- Shows all stored notes with their last modification timestamp
- Notes are stored in `.agent-go/notes/` directory
- Useful for quickly finding existing notes

### `/notes view <name>`

Displays the full content of a specific note.

**Usage:**

```
/notes view <name>
```

**Parameters:**

- `name`: Name of the note to view

**Example:**

```
> /notes view api_endpoint
=== api_endpoint ===
Created: 2025-11-25 08:30:00
Updated: 2025-11-25 08:30:00

The API endpoint is https://api.example.com/v1
```

**Notes:**

- Shows the complete note content with metadata
- Includes creation and update timestamps
- Useful for reviewing stored information

### Notes Management via AI

You can also ask the AI to create, update, or delete notes through natural conversation:

> Create a note called "api_endpoint" with content "The API endpoint is <https://api.example.com/v1>"
Created note: api_endpoint

> Update the api_endpoint note to use v2
Updated note: api_endpoint

> Delete the old_database_config note
Deleted note: old_database_config

> Show me all my notes
Notes:

- api_endpoint (updated: 2025-11-25 08:30)
- database_schema (updated: 2025-11-25 09:15)

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

### `/mode` (or `/plan`)

Toggles between **Plan** mode and **Build** mode.

- **Build Mode** (Default): The agent has access to all tools, including `execute_command`. This is the standard mode for getting things done.
- **Plan Mode**: The agent is restricted from executing commands. It can only think and use the `suggest_plan` tool to propose a detailed plan to the user. This is useful for complex tasks where you want to agree on a strategy before any code is modified.

**Usage:**

```
/mode
# or
/plan
```

**Example:**

```
> /mode
Switched to Plan mode.

> /plan
Switched to Build mode.
```

### `/ask on|off`

Toggles between **Ask** mode and **YOLO** mode for command execution.

- **Ask Mode** (Default): The agent will ask for your confirmation before executing potentially dangerous tools (like `execute_command`).
- **YOLO Mode**: "You Only Look Once" - The agent will execute commands immediately without asking for confirmation. Use with caution!

**Usage:**

```
/ask on
/ask off
```

**Example:**

```
> /ask off
Switched to YOLO mode.

> /ask on
Switched to Ask mode.
```

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

### MCP Tool Usage

The AI agent can automatically use tools from connected MCP servers:

```

> What time is it in Tokyo?
[Using MCP tool: get_current_time from server 'time']
Current time in Tokyo: 3:45 PM JST

> Get the latest documentation for Next.js routing
[Using MCP tool: resolve-library-id from server 'context7']
[Using MCP tool: get-library-docs from server 'context7']
Based on the Next.js documentation, routing works as follows...

```

**How MCP Tools Work:**
- The AI automatically detects when an MCP tool can help
- Tools are called transparently during conversation
- Tool results are incorporated into the AI's response
- No special syntax required - just ask naturally

### Todo List Management

Create and manage todos through natural conversation:

```

> Create a todo list for today's tasks
Created todo list:

- [ID: 1] Review pull requests (pending)
- [ID: 2] Update documentation (pending)
- [ID: 3] Deploy to staging (pending)

> Mark todo 1 as in-progress
Updated: Review pull requests (in-progress)

> Show my todo list
Current Todo List:

- [ID: 1] Review pull requests (in-progress)
- [ID: 2] Update documentation (pending)
- [ID: 3] Deploy to staging (pending)

> Complete todo 1
Updated: Review pull requests (completed)

```
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

### MCP Server Debugging

If MCP servers aren't connecting:

```
> /mcp list
Configured MCP servers:
- context7: npx -y @upstash/context7-mcp

# Check if the MCP command works standalone
> /shell
shell> npx -y @upstash/context7-mcp
# Should see MCP server output
shell> exit
```

**Common MCP Issues:**

- **Server not found**: Ensure the command is correct and the package is available
- **Connection failed**: Check network connectivity and package installation
- **Tool call errors**: Verify the tool name and arguments match the server's schema

### Debug Mode

For troubleshooting configuration issues:

```bash
# Enable debug logging
export DEBUG=1
./agent-go
```

This will provide detailed logging information to help diagnose issues.
