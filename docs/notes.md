# Notes Feature

The notes feature allows the agent to create, update, and delete persistent notes that are automatically injected into the system prompt for context across sessions.

## Overview

Notes are stored in the `.agent-go/notes` directory in JSON format. They persist across sessions and are automatically loaded into the system prompt when the agent starts.

## Agent Tools

The agent has access to three tools for managing notes:

### create_note

Creates a new note with a unique name and content.

**Parameters:**
- `name` (string, required): Unique identifier for the note
- `content` (string, required): The content of the note

**Example:**
```json
{
  "name": "api_endpoint",
  "content": "The API endpoint is https://api.example.com/v1"
}
```

### update_note

Updates the content of an existing note.

**Parameters:**
- `name` (string, required): Name of the note to update
- `content` (string, required): New content for the note

**Example:**
```json
{
  "name": "api_endpoint",
  "content": "The API endpoint is now https://api.example.com/v2"
}
```

### delete_note

Deletes an existing note.

**Parameters:**
- `name` (string, required): Name of the note to delete

**Example:**
```json
{
  "name": "api_endpoint"
}
```

## User Commands

Users can interact with notes using slash commands:

### /notes list

Lists all existing notes with their last update time.

```
> /notes list
Notes:
  - api_endpoint (updated: 2025-11-25 08:30)
  - database_schema (updated: 2025-11-25 09:15)
```

### /notes view <name>

Displays the full content of a specific note.

```
> /notes view api_endpoint
=== api_endpoint ===
Created: 2025-11-25 08:30:00
Updated: 2025-11-25 08:30:00

The API endpoint is https://api.example.com/v1
```

## Autocomplete Support

The `/notes` command supports tab completion:
- Type `/notes` and press Tab to see `list` and `view` options
- Type `/notes view` and press Tab to see available note names

## System Prompt Injection

All notes are automatically injected into the system prompt in the following format:

```
=== Agent Notes ===

[api_endpoint]
The API endpoint is https://api.example.com/v1

[database_schema]
The database uses PostgreSQL with the following schema...
```

This ensures the agent has access to all stored knowledge across sessions.

## Storage Format

Notes are stored as JSON files in `.agent-go/notes/`:

```json
{
  "name": "api_endpoint",
  "content": "The API endpoint is https://api.example.com/v1",
  "created_at": "2025-11-25T08:30:00Z",
  "updated_at": "2025-11-25T08:30:00Z"
}
```

## Use Cases

- **Project Information**: Store project-specific details like API endpoints, database connections, or coding conventions
- **Repository Details**: Record repository structure, build commands, or deployment procedures
- **User Preferences**: Save user preferences, coding style, or frequently used commands

## Implementation Details

The notes feature is implemented in [`src/notes.go`](../src/notes.go) and integrated into:
- [`src/tools.go`](../src/tools.go) - Tool definitions
- [`src/processor.go`](../src/processor.go) - Tool call processing
- [`src/main.go`](../src/main.go) - System prompt building and slash commands
- [`src/completion.go`](../src/completion.go) - Autocomplete support