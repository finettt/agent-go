# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0](https://github.com/finettt/agent-go/compare/v1.3.0...v2.0.0) (2026-01-25)


### ⚠ BREAKING CHANGES

* **agents:** OperationMode config is deprecated. Use /plan command to toggle between plan/build agents instead.

### Added

* **agent:** allow sub-agents to select model type via spawn_agent tool ([59af5fe](https://github.com/finettt/agent-go/commit/59af5fe05c463322f940626e5889fac4e93e879a))
* **agents:** add per-agent tool management with allow/deny policies ([d83e58e](https://github.com/finettt/agent-go/commit/d83e58eb7bf2f5e3d805576e3145dbcd8affe587))
* **agents:** implement agent-based plan/build mode architecture ([90ada5b](https://github.com/finettt/agent-go/commit/90ada5b486b00135614111017f6d60963cca3c89))
* **ci:** Add pr-review from QluxLab/Reviewer ([f504ad1](https://github.com/finettt/agent-go/commit/f504ad14609cede9220a8eb1846b05ec77e75c1e))
* **ci:** Add pr-review from QluxLab/Reviewer ([#25](https://github.com/finettt/agent-go/issues/25)) ([f81f144](https://github.com/finettt/agent-go/commit/f81f144b412cf3614438a4616065f74d0b6c5e8f))
* **compression:** preserve system prompts and task continuity during context compression ([cd8d10e](https://github.com/finettt/agent-go/commit/cd8d10ec74619b09a8e8c28ab76d1d7c2871a5d8))
* **core:** implement checkpoint and rollback system ([f8b13e1](https://github.com/finettt/agent-go/commit/f8b13e139c30d63992c359a43cbc36047f58a9c2))
* **core:** implement checkpoint and rollback system ([#21](https://github.com/finettt/agent-go/issues/21)) ([971a39a](https://github.com/finettt/agent-go/commit/971a39ae5df56a7455a065ead6949ae32e734ae5))
* **core:** integrate mini model for utility tasks ([23499c0](https://github.com/finettt/agent-go/commit/23499c032a3a023d68bc635c0ec35f12b372dfd7))


### Fixed

* **tokens:** implement "Last Usage" algorithm for accurate context tracking ([deb0d37](https://github.com/finettt/agent-go/commit/deb0d370be7ffe16023245c461cbb9a876f36d91))


### Changed

* add PR review workflow ([896e0d8](https://github.com/finettt/agent-go/commit/896e0d83bc7d19ca81d2a6b28db457e368a1a821))
* add release-please automation ([2e3c158](https://github.com/finettt/agent-go/commit/2e3c1585f5e4aabc090fe5f4ad7bff0850bd0dbb))
* **agents:** replace global activeAgentDef with agent.AgentDefName ([5ee620d](https://github.com/finettt/agent-go/commit/5ee620d4cd2c8e1d8e0a5a51431bc807af90682e))
* **claude:** add claude code router workflow ([19ff0b0](https://github.com/finettt/agent-go/commit/19ff0b05e00c4c69248f4ed487de193f978a3903))
* **claude:** add claude code router workflow ([#23](https://github.com/finettt/agent-go/issues/23)) ([816b185](https://github.com/finettt/agent-go/commit/816b185d41fcdd1e8b014ea1ce18436e94f7af52))
* **claude:** add github token input to action step ([124d1e4](https://github.com/finettt/agent-go/commit/124d1e4abc9beb6057ef69de80a0dad14ee6c822))
* **claude:** add triggers for pr reviews and issues ([ef2734a](https://github.com/finettt/agent-go/commit/ef2734a7ea3214a37402a693b93a8ff9d69ce835))
* **claude:** configure workflow for interactive pr reviews ([b74690c](https://github.com/finettt/agent-go/commit/b74690cc1a5116dacfc31e2685d4ffd8e26d7205))
* **claude:** fix environment variable validation failure ([8991021](https://github.com/finettt/agent-go/commit/89910217bd44370bcc01a34ab15ec40dc65c5fb9))
* **claude:** pass github token to action input ([96d0220](https://github.com/finettt/agent-go/commit/96d022009884d01412a2807966267ca8b272e0ce))
* **claude:** refactor workflow for general assistant capabilities ([598de39](https://github.com/finettt/agent-go/commit/598de395afc3c6674ac8d37e659aee97d2ffa688))
* **claude:** remove failed experimental workflow ([a76a27e](https://github.com/finettt/agent-go/commit/a76a27e701b38fce70d898c8e722067b0246e0e9))
* **claude:** restrict triggers to explicit [@claude](https://github.com/claude) mentions ([e151850](https://github.com/finettt/agent-go/commit/e151850c78cee6e275e6f4f682e18f48fca7a11c))
* **claude:** unify and configure claude workflow ([2be782d](https://github.com/finettt/agent-go/commit/2be782dc156d3630f275b230280e6090a2e37428))
* **cli:** polish CLI output and logo alignment ([1f98007](https://github.com/finettt/agent-go/commit/1f9800754d9163db6831e854620858a61301abdc))
* **cli:** update greeting and colorize warnings ([7f3ff07](https://github.com/finettt/agent-go/commit/7f3ff076da603c9923b546314232f09bc5bdc8fd))
* **cli:** update subagents help text to mention mini model support ([fc9c867](https://github.com/finettt/agent-go/commit/fc9c8679b8ab4e794852c2c62200b78d3d28dfa0))
* **install:** update rolling install syntax to use environment variable ([3f705ac](https://github.com/finettt/agent-go/commit/3f705ac4ed30168d22973a741a70c7548daac65e))
* **install:** update rolling install syntax to use environment variable ([#31](https://github.com/finettt/agent-go/issues/31)) ([8048b48](https://github.com/finettt/agent-go/commit/8048b48f95156bc06d696713db0c720bf49bcce3))
* **pr-review:** migrate from QluxLab/Reviewer to qodo-ai/pr-agent ([c8b8040](https://github.com/finettt/agent-go/commit/c8b804039a4ed6be57c65afe0c63fd6f38512180))
* **pr-review:** upgrade QluxLab/Reviewer to v1.1 ([8922d99](https://github.com/finettt/agent-go/commit/8922d993aa611d837fc1cabe7d86f6386bdf00db))
* **pr-review:** upgrade QluxLab/Reviewer to v1.2 ([8a9c982](https://github.com/finettt/agent-go/commit/8a9c9825ee1043f097880321b13aecf68434a143))
* **pr-review:** upgrade QluxLab/Reviewer to v1.3 ([054e535](https://github.com/finettt/agent-go/commit/054e53586af266dfab8c67ec9444ad3cefff67a1))
* **pr-review:** upgrade QluxLab/Reviewer to v2.0 ([85b21f3](https://github.com/finettt/agent-go/commit/85b21f3e67a6bba864e48cbba366c6c456129c45))
* **pr-review:** upgrade QluxLab/Reviewer to v2.1 ([04531c8](https://github.com/finettt/agent-go/commit/04531c864ba96e253f317e3898b74edfad7c3a29))
* **pr-review:** upgrade QluxLab/Reviewer to v2.2 ([c526047](https://github.com/finettt/agent-go/commit/c526047281080fb626b3599b8f158d059ef2e09a))
* **pr-review:** upgrade QluxLab/Reviewer to v2.3 ([6a13273](https://github.com/finettt/agent-go/commit/6a13273c60971a1c24c08bdcde3a1c1d44163289))
* **pr-review:** upgrade QluxLab/Reviewer to v2.4 ([bf0c871](https://github.com/finettt/agent-go/commit/bf0c871cb958c62c3eb26006769573d7dd9cc243))
* **pr-review:** upgrade QluxLab/Reviewer to v2.4 ([#27](https://github.com/finettt/agent-go/issues/27)) ([4553580](https://github.com/finettt/agent-go/commit/4553580010b47f77009b084f21933dca030dfbcc))
* **tools:** unify operation mode and agent policy filtering ([0618e3e](https://github.com/finettt/agent-go/commit/0618e3eaf7a748332c134f49e7553aaa170c9e7e))

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
