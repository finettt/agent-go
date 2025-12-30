# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] - 2025-12-30

### Added
- **Skills System**:
  - Implement custom skill system with script-based tools (9790d8c)
  - Enable direct execution of .sh files (c30f2ab)
- **Session Management**:
  - New `/session view <name>` command to inspect session details, metadata, and recent messages (e0139f2)
  - Persist token usage statistics in session data (87d7596)
- **CLI Enhancements**:
  - New `/current` command to display the current in-progress task (faa638a)
  - Added 'all' option to switch to YOLO mode during execution (20f2804)
  - Plan viewing and editing capabilities (3d4a9cb)
  - Progress bar for cost tracking (3338123)
  - Enhanced output readability with bullets (61ee556)
  - Changed prompt to '?' when in plan mode (69210f3)
  - Sandbox mode and environment detection (85cd03c)
- **Tool Capabilities**:
  - Added `clear_todo` tool for cleaning up tasks (2de9760)
  - Enhanced planning tool and prompt display (5c59ebe)
  - Background command support via `/bg` commands (eee1099)
- **Installer**:
  - Rolling update support (e72b5a0)
  - Pre-built binary downloads (67a3bd6)

### Changed
- **CLI/UX**:
  - Refined CLI user experience and interface (84419bc)
  - Redesigned `/cost` command with progress bar visualization (28a363d)
  - Set usage display to silent by default (77b878b)
- **Configuration**:
  - Increased default context limit (77b878b)
- **Documentation**:
  - Removed duplicate content and fixed formatting errors across docs (aa62567)
  - Updated documentation for modes and commands (b449120)
  - Added rolling install instructions (a67dcd0)
- **License**:
  - Updated license from MIT to Apache 2.0 with Commons Clause (1062479)

### Removed
- **Streaming**:
  - Removed streaming functionality for better stability (ddf72e7)

### Fixed
- **Windows Compatibility**:
  - Improved compatibility for installer and edit command (cd92922)
- **Security**:
  - Fix for potentially unsafe quoting (044acb9)

## [1.2.0] - 2025-12-18

### Added
- **Agent Studio**: Complete agent management system with `/agent` commands
  - `/agent studio` - open agent studio interface
  - `/agent list` - view available agents
  - `/agent view <name>` - inspect agent configuration
  - `/agent use <name>` - switch to specific agent
  - `/agent clear` - clear agent-specific context
  - `/agent rm <name>` - remove custom agent
- **Task-Specific Agents**: Persistent agents stored in `~/.config/agent-go/agents`
- **Built-in Default Agent**: Protected default agent that cannot be deleted or overwritten
- **Subagent Agent Selection**: Support for selecting specific agents via `{"agent": "name"}` parameter
- **Agent Discovery**: Available agent names injected into system prompt for easy discovery
- **Enhanced Autocompletion**: Readline completion for `/agent` commands and agent names
- **/edit Command**: New slash command to compose prompts using nano text editor
- **Session Management**: Complete session save/restore functionality
  - `/session new` - create new session and save current context
  - `/session list` - view saved sessions
  - `/session restore <name>` - restore previous session
  - `/session rm <name>` - delete saved session
  - `name_session` tool for agent to rename sessions
- **Background Command Execution**: Support for running shell commands in background
  - `execute_command` with `background` parameter
  - `kill_background_command`, `get_background_logs`, `list_background_commands` tools
  - Application prevents exit when background tasks are running
- **Usage Tracking**: Granular token and tool usage monitoring
  - `/usage` command for detailed token usage
  - `/cost` command for cost tracking
- **Security Review**: `/security` command to spawn subagent for code review
- **Note Mentions**: Support for `#note-name` syntax to inject note content
- **Verbose Mode**: Enhanced logging control with `/verbose` command
- **Plan Mode**: Implementation with execution safety for complex multi-step tasks
- **Reasoning Support**: Chain-of-thought reasoning with "Think..." indicator
- **Token Formatting**: Human-readable K/M suffixes for large token counts

### Changed
- **CLI Enhancements**:
  - New color scheme with consistent theming (#FF93FB, #FFF, #AAA)
  - Simplified tool output with concise status messages
  - Improved startup banner with model and working directory info
  - Enhanced autocomplete for `/help`, `/ask`, `/mode`, `/plan` commands
  - Numeric completion options for `/contextlength`
  - Reduced noisy subagent logging (configurable verbosity)
- **API Improvements**:
  - Support for `reasoning_content` in responses
  - Enhanced message and delta structures
  - Better streaming response handling
- **Subagent Enhancements**:
  - Parallel execution support with goroutines and WaitGroup
  - Thread-safe console I/O and message history updates
  - Verbose mode configuration for detailed logging
- **Security Improvements**:
  - Removed verbose mode configuration flag
  - Made logging unconditional for transparency
  - Enhanced command validation and platform handling
- **Session Management**:
  - Automatic session saving on exit and context clear
  - Enhanced startup with session restoration
  - Improved context switching between sessions

### Technical Details
- Agent configurations stored as JSON in `~/.config/agent-go/agents/`
- Background processes managed with process IDs and output buffers
- Usage tracking with granular per-command and per-agent metrics
- Session data stored with timestamps and context compression
- Enhanced system prompt injection for agent discovery and cross-session context
- Improved memory management for large conversations and background tasks

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