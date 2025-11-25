# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-11-25

### Added
- **Notes Feature**: Persistent project memory with agent tools (`create_note`, `update_note`, `delete_note`)
- Notes are automatically injected into the system prompt for cross-session context
- User commands: `/notes list` and `/notes view <name>` for viewing stored notes
- Tab autocomplete support for `/notes` commands and note names
- New documentation: [`docs/notes.md`](docs/notes.md) with complete feature guide

### Changed
- Updated tool definitions in `src/tools.go` to include note management tools
- Enhanced system prompt to include agent notes section when notes exist
- Extended autocomplete system in `src/completion.go` for notes commands

### Technical Details
- Notes stored as JSON files in `.agent-go/notes/` directory
- Each note includes `name`, `content`, `created_at`, and `updated_at` fields
- Sanitized filenames to prevent path traversal attacks
- Notes automatically loaded at startup and injected into system prompt

## [1.0.0] - Initial Release

### Added
- Intelligent command execution with shell support
- MCP (Model Context Protocol) integration with context7 default server
- Sub-agent delegation for complex multi-step tasks
- Unlimited conversation memory with automatic context compression
- Retrieval-Augmented Generation (RAG) for local knowledge base
- Todo list management for task tracking
- Shell mode for direct command execution
- Streaming mode for real-time response generation
- Cross-platform support (macOS, Linux, Windows)
- Configuration via environment variables, config file, or CLI arguments
- Custom agent behavior via `AGENTS.md` file
- Real-time token tracking