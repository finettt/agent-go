# Agent-Go Documentation

This directory contains the full documentation for Agent-Go. For a quick start, see the main [README.md](../README.md) file.

## Table of Contents

- [**Architecture**](architecture.md): A deep dive into the system architecture, components, and data flow, including MCP integration, sub-agent systems, and notes management.
- [**Commands**](commands.md): A complete reference for all slash commands and CLI features, including MCP management, todo commands, and notes management.
- [**Configuration**](configuration.md): Detailed guidance on configuration options, environment variables, and file-based settings, including MCP server configuration.
- [**Development Guide**](development.md): Instructions for setting up the development environment, running tests, and contributing to the project.
- [**Examples and Best Practices**](examples.md): Practical examples, advanced workflows, and tips for using Agent-Go effectively, including MCP, todo management, and notes management examples.

## What's New

### Agent Studio

Complete agent management system for creating, managing, and using task-specific agents:

- **Agent Studio Interface**: `/agent studio` command for interactive agent creation
- **Agent Management**: `/agent list`, `/agent view <name>`, `/agent use <name>`, `/agent clear`, `/agent rm <name>`
- **Persistent Storage**: Agent definitions stored in `~/.config/agent-go/agents/*.json`
- **Subagent Support**: Task-specific agents can be used by subagents via `{"agent": "name"}` parameter
- **Built-in Agent**: Protected `default` agent that cannot be deleted or overwritten

### Session Management

Save and restore conversation sessions for seamless context switching:

- **Session Creation**: `/session new` - save current context and create new session
- **Session Listing**: `/session list` - view all saved sessions
- **Session Restoration**: `/session restore <name>` - restore previous session
- **Session Deletion**: `/session rm <name>` - delete saved session
- **Agent Tool**: `name_session` tool for agents to rename sessions

### Background Command Execution

Run shell commands in the background with full monitoring and management:

- **Background Execution**: `execute_command` with `background` parameter
- **Command Management**: `kill_background_command`, `get_background_logs`, `list_background_commands` tools
- **Application Safety**: Prevents exit when background tasks are running
- **Real-time Monitoring**: Track command output and status

### Tool Loop Detection

Automatically detects and prevents infinite loops where the AI repeatedly calls the same tool:

- **Loop Detection**: Monitors tool calls within API responses and across iterations
- **Threshold**: Triggers after 3 repeated calls of the same tool with identical arguments
- **Intervention**: Injects a stop message to guide the model toward alternative approaches
- **Scope**: Applied to all execution modes (CLI, task, pipeline, editor) and sub-agents
- **Implementation**: Uses `checkToolLoop()` function with `getToolCallSignature()` helper

### Empty Response Retry Logic

Intelligent retry mechanism for handling empty model responses:

- **Detection**: Identifies empty responses (zero choices, empty content, no tool calls)
- **Retry Limit**: Up to 3 automatic retry attempts
- **Feedback**: Warns user about empty responses and retry attempts via stderr
- **Scope**: Applied to all execution modes (interactive CLI, single task mode, pipeline mode)
- **Resilience**: Improves stability by handling transient API issues gracefully

### Time Context Injection

Provides temporal awareness to the AI:

- **Time Format**: UTC time in ISO 8601 format, local time in RFC1123 format
- **Timezone Info**: Includes timezone name and UTC offset
- **Format**: "Current Time: 2026-01-29T14:30:00Z (UTC) | Local: Thu, 29 Jan 2026 17:30:00 MSK | Timezone: MSK (UTC+03:00)"
- **Implementation**: `getCurrentTimeContext()` function in `system.go`
- **Integration**: Injected as a system message in API requests
- **Benefits**: Enables time-aware operations and improves context for scheduling tasks

### Enhanced Features

- **Usage Tracking**: `/usage <1|2|3>` command for verbosity control (Silent/Basic/Detailed) and `/cost` for detailed context and session statistics
- **Security Review**: `/security` command to spawn subagent for code review
- **Note Mentions**: Support for `#note-name` syntax to inject note content
- **Plan/Build Agents**: Separate agent definitions for planning vs. implementation with automatic mode switching
- **Reasoning Support**: Chain-of-thought reasoning with "Think..." indicator
- **Token Formatting**: Human-readable K/M suffixes for large token counts
- **Dual Token Tracking**: "Last Usage" algorithm tracks both current context size and cumulative session usage

### MCP (Model Context Protocol) Integration

Agent-Go now integrates with MCP servers to extend functionality with external tools and resources:

- **Default context7 server**: Automatically configured for accessing up-to-date library documentation
- **Custom MCP servers**: Add any MCP-compatible server with `/mcp add`
- **Tool discovery**: AI automatically discovers and uses available tools from connected servers
- **Dynamic configuration**: Manage servers via `/mcp add`, `/mcp remove`, and `/mcp list` commands

### Todo List Management

Built-in todo list management for tracking tasks across sessions:

- **Persistent storage**: Todos stored in `~/.config/agent-go/todos/`
- **Per-agent lists**: Each agent (main and sub-agents) maintains separate todo lists
- **Status tracking**: Supports pending, in-progress, and completed statuses
- **AI-driven management**: Create, update, and view todos through natural conversation
- **Slash command**: Quick access with `/todo` command

### Notes Management

Persistent notes feature for storing and retrieving important information across sessions:

- **Persistent storage**: Notes stored in `.agent-go/notes/` in JSON format
- **AI-driven management**: Create, update, and delete notes through natural conversation
- **System integration**: Notes automatically injected into system prompt for context
- **Slash commands**: Quick access with `/notes list` and `/notes view <name>`
- **Autocomplete support**: Tab completion for note names and commands
- **Note Mentions**: Support for `#note-name` syntax to inject note content

### Enhanced CLI Features

- **Streaming mode**: Real-time response generation with `/stream on`
- **Sub-agent control**: Toggle sub-agent spawning with `/subagents on|off`
- **Unlimited context**: Automatic context compression at 75% of model context length
- **System information**: Automatic OS, architecture, and environment detection
- **Command-line mode**: Execute single tasks directly from command line
- **New Color Scheme**: Enhanced theming with consistent colors (#FF93FB, #FFF, #AAA)
- **Improved Autocomplete**: Enhanced completion for `/agent` commands and agent names
- **/edit Command**: New slash command to compose prompts using nano text editor

## Quick Navigation

### For New Users

1. Start with the main [README.md](../README.md) for installation and setup
2. Learn about [Commands](commands.md) for interactive usage
3. Review [Configuration](configuration.md) for customization options
4. Check [Examples](examples.md) for practical use cases
5. Explore [Notes Management](notes.md) for persistent information storage

### For Developers

1. Review [Architecture](architecture.md) to understand the system design
2. Read [Development Guide](development.md) for contribution guidelines
3. Explore [Examples](examples.md) for advanced workflows
4. Check [Commands](commands.md) for testing and debugging commands

### For Advanced Users

1. Learn about [MCP Integration](commands.md#mcp-model-context-protocol-commands) to extend functionality
2. Master [Todo Management](commands.md#todo) for project tracking
3. Explore [Notes Management](notes.md) for persistent knowledge storage
4. Discover [Agent Studio](commands.md#agent-studio) for creating custom agents
5. Utilize [Session Management](commands.md#session) for context switching and continuity
6. Explore [Background Commands](commands.md#background-commands) for parallel task execution
7. Use [Plan/Build Workflow](commands.md#plan) - `/plan view` & `/plan edit` for strategic planning
8. Explore [Advanced Workflows](examples.md#advanced-workflows) for complex scenarios
9. Optimize [Configuration](configuration.md) for your use case
10. Learn about [Usage Tracking](commands.md#usage-tracking) - `/usage` for verbosity, `/cost` for detailed stats
