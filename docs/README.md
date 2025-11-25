# Agent-Go Documentation

This directory contains the full documentation for Agent-Go. For a quick start, see the main [README.md](../README.md) file.

## Table of Contents

- [**Architecture**](architecture.md): A deep dive into the system architecture, components, and data flow, including MCP integration, sub-agent systems, and notes management.
- [**Commands**](commands.md): A complete reference for all slash commands and CLI features, including MCP management, todo commands, and notes management.
- [**Configuration**](configuration.md): Detailed guidance on configuration options, environment variables, and file-based settings, including MCP server configuration.
- [**Development Guide**](development.md): Instructions for setting up the development environment, running tests, and contributing to the project.
- [**Examples and Best Practices**](examples.md): Practical examples, advanced workflows, and tips for using Agent-Go effectively, including MCP, todo management, and notes management examples.

## What's New

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

### Enhanced Features

- **Streaming mode**: Real-time response generation with `/stream on`
- **Sub-agent control**: Toggle sub-agent spawning with `/subagents on|off`
- **Unlimited context**: Automatic context compression at 75% of model context length
- **System information**: Automatic OS, architecture, and environment detection
- **Command-line mode**: Execute single tasks directly from command line

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
4. Explore [Advanced Workflows](examples.md#advanced-workflows) for complex scenarios
5. Optimize [Configuration](configuration.md) for your use case
