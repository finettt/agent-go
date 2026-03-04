# Interactive Terminal Sessions

The agent-go framework provides interactive terminal session capabilities, allowing agents to open, control, and interact with terminal applications like `nano`, `vim`, `bash`, and more.

## Overview

Terminal sessions use pseudo-terminals (PTY) to provide real terminal behavior, including:
- Full ANSI escape sequence support (colors, cursor positioning)
- Interactive input/output
- Support for terminal-based applications
- Session management and tracking

## Available Tools

### 1. `open_terminal_session`

Start an interactive terminal session with a command.

**Parameters:**
- `command` (string, required): The command to run in the terminal session

**Returns:**
```json
{
  "session_id": "sess_1",
  "pid": 12345
}
```

**Example:**
```json
{
  "command": "nano myfile.txt"
}
```

### 2. `send_terminal_input`

Send input (text or keycodes) to an active terminal session.

**Parameters:**
- `session_id` (string, required): The terminal session ID
- `input` (string, required): Text to type or keycode (e.g., "Ctrl+S", "Enter", "Hello World")

**Returns:**
```json
{
  "status": "sent"
}
```

**Example:**
```json
{
  "session_id": "sess_1",
  "input": "Ctrl+S"
}
```

### 3. `read_terminal_output`

Read output from an active terminal session.

**Parameters:**
- `session_id` (string, required): The terminal session ID
- `bytes` (integer, optional): Number of bytes to read (default: 4096)
- `read_all` (boolean, optional): Read all available output (default: false)

**Returns:**
```json
{
  "output": "\x1b[32mHello World\x1b[0m"
}
```

**Example:**
```json
{
  "session_id": "sess_1",
  "bytes": 2048
}
```

### 4. `close_terminal_session`

Close an active terminal session.

**Parameters:**
- `session_id` (string, required): The terminal session ID to close

**Returns:**
```json
{
  "status": "closed"
}
```

**Example:**
```json
{
  "session_id": "sess_1"
}
```

### 5. `list_terminal_sessions`

List all active terminal sessions.

**Parameters:** None

**Returns:**
```
Active Terminal Sessions:
- Session: sess_1 | PID: 12345 | Command:  | Duration: 1m30s | Status: Running
- Session: sess_2 | PID: 12346 | Command:  | Duration: 45s | Status: Running
```

## Human-Readable Keycodes

The agent uses simple, readable key names instead of escape codes:

| Key Name | Escape Code | Use Case |
|----------|-------------|----------|
| `Ctrl+A` through `Ctrl+Z` | `\x01` - `\x1a` | Control characters |
| `Enter` | `\n` | Confirm/Submit |
| `Tab` | `\t` | Tab completion |
| `Escape` | `\x1b` | Cancel/Exit |
| `Backspace` | `\x7f` | Delete character |
| `ArrowUp` | `\x1b[A` | Navigate up |
| `ArrowDown` | `\x1b[B` | Navigate down |
| `ArrowRight` | `\x1b[C` | Navigate right |
| `ArrowLeft` | `\x1b[D` | Navigate left |
| `F1` - `F12` | `\x1b[11~` - `\x1b[24~` | Function keys |
| `PageUp` | `\x1b[5~` | Page up |
| `PageDown` | `\x1b[6~` | Page down |
| `Home` | `\x1b[H` | Go to start |
| `End` | `\x1b[F` | Go to end |

Any text that doesn't match a keycode is sent as plain text.

## Common Editor Keycodes

### Nano
- `Ctrl+S`: Save file
- `Ctrl+X`: Exit
- `Ctrl+K`: Cut line
- `Ctrl+U`: Paste
- `Ctrl+W`: Search

### Vim
- `Escape`: Exit insert mode
- `:wq` + `Enter`: Save and quit
- `:q!` + `Enter`: Quit without saving
- `i`: Enter insert mode
- `dd`: Delete line (in normal mode)

### Bash/Shell
- `Ctrl+C`: Interrupt current command
- `Ctrl+D`: End of input / Exit shell
- `Ctrl+L`: Clear screen
- `Tab`: Auto-complete
- `ArrowUp`/`ArrowDown`: Command history

## Usage Examples

### Example 1: Edit a file with nano

```
1. open_terminal_session({"command": "nano myfile.txt"})
   → Returns: {"session_id": "sess_1", "pid": 12345}

2. send_terminal_input({"session_id": "sess_1", "input": "Hello World"})
   → Sends text to nano

3. send_terminal_input({"session_id": "sess_1", "input": "Ctrl+S"})
   → Saves file

4. send_terminal_input({"session_id": "sess_1", "input": "Ctrl+X"})
   → Exits nano

5. close_terminal_session({"session_id": "sess_1"})
   → Closes session
```

### Example 2: Run interactive Python

```
1. open_terminal_session({"command": "python3"})
   → Returns: {"session_id": "sess_2", "pid": 12346}

2. send_terminal_input({"session_id": "sess_2", "input": "print('Hello from Python')"})
   → Sends Python command

3. send_terminal_input({"session_id": "sess_2", "input": "Enter"})
   → Executes command

4. read_terminal_output({"session_id": "sess_2", "read_all": true})
   → Reads Python output

5. send_terminal_input({"session_id": "sess_2", "input": "exit()"})
   → Exits Python

6. close_terminal_session({"session_id": "sess_2"})
   → Closes session
```

### Example 3: Interactive bash session

```
1. open_terminal_session({"command": "bash"})
   → Returns: {"session_id": "sess_3", "pid": 12347}

2. send_terminal_input({"session_id": "sess_3", "input": "ls -la"})
   → Types command

3. send_terminal_input({"session_id": "sess_3", "input": "Enter"})
   → Executes command

4. read_terminal_output({"session_id": "sess_3", "bytes": 4096})
   → Reads command output

5. send_terminal_input({"session_id": "sess_3", "input": "exit"})
   → Exits bash

6. close_terminal_session({"session_id": "sess_3"})
   → Closes session
```

## Safety and Limitations

### Session Limits
- **Maximum session duration**: 5 minutes per session
- **Maximum concurrent sessions**: 10 active sessions
- **Auto-cleanup**: Sessions are cleaned up on agent exit

### Security
- Commands are validated before execution
- Dangerous patterns (e.g., `rm -rf /`) are blocked
- All sessions run under the agent's user permissions

### Best Practices
1. Always close sessions when done to free resources
2. Read output regularly to prevent buffer overflow
3. Use appropriate timeouts for long-running commands
4. Check session status with `list_terminal_sessions`

## Implementation Details

### Platform Support
- **Linux**: Uses `github.com/creack/pty` for PTY support
- **macOS**: Uses `github.com/creack/pty` for PTY support
- **Windows**: Uses pseudo-console via `golang.org/x/sys/windows`

### Architecture
- Each session runs in its own goroutine
- Output is buffered in memory
- Input is sent via channels
- PTY handles ANSI escape sequences automatically

### Session Storage
Sessions are stored in memory with the following structure:
```go
type TerminalSession struct {
    ID           string
    PID          int
    PTY          *os.File
    StdoutBuffer *bytes.Buffer
    InputChan    chan string
    DoneChan     chan struct{}
    StartTime    time.Time
}
```

## Troubleshooting

### Session not responding
- Check if the session is still active with `list_terminal_sessions`
- Ensure you're sending the correct keycodes
- Try sending `Ctrl+C` to interrupt hung processes

### Output appears garbled
- The output includes ANSI escape sequences (colors, formatting)
- Use a terminal emulator or strip ANSI codes if needed
- Read output after each significant action

### Session timeout
- Sessions automatically close after 5 minutes
- Close and reopen if you need longer sessions
- Consider using background commands for long-running processes

## See Also

- [Tool Management](tool-management.md) - Managing agent tools
- [Commands](commands.md) - Agent command reference
- [Development](development.md) - Contributing to agent-go