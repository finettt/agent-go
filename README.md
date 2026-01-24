<div align="center">
  <img src="https://github.com/user-attachments/assets/b8914faa-998d-487e-9173-7008d75b36df" width="200" height="200" alt="Agent-Go Logo" />
  <h1>./agent-go</h1>
  <p>
    <strong>Chat less. Execute more.</strong> <br>
    The native Go agent that lives in your terminal and does the work.
  </p>

  <!-- Branded Badges -->
  <a href="https://golang.org/">
    <img src="https://img.shields.io/badge/Language-Go_1.25-FF69F6?style=for-the-badge&logo=go&logoColor=white" alt="Go Version" />
  </a>
  <a href="https://github.com/finettt/agent-go/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/License-Apache_2.0-0F0F0F?style=for-the-badge" alt="License" />
  </a>
  <a href="#">
    <img src="https://img.shields.io/badge/Platform-Win_|_Mac_|_Linux-4AF6F6?style=for-the-badge&logo=linux&logoColor=black" alt="Platform" />
  </a>
  <a href="https://zread.ai/finettt/agent-go" target="_blank">
    <img src="https://img.shields.io/badge/Ask_Zread-_.svg?style=for-the-badge&color=00b0aa&labelColor=000000&logo=data%3Aimage%2Fsvg%2Bxml%3Bbase64%2CPHN2ZyB3aWR0aD0iMTYiIGhlaWdodD0iMTYiIHZpZXdCb3g9IjAgMCAxNiAxNiIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTQuOTYxNTYgMS42MDAxSDIuMjQxNTZDMS44ODgxIDEuNjAwMSAxLjYwMTU2IDEuODg2NjQgMS42MDE1NiAyLjI0MDFWNC45NjAxQzEuNjAxNTYgNS4zMTM1NiAxLjg4ODEgNS42MDAxIDIuMjQxNTYgNS42MDAxSDQuOTYxNTZDNS4zMTUwMiA1LjYwMDEgNS42MDE1NiA1LjMxMzU2IDUuNjAxNTYgNC45NjAxVjIuMjQwMUM1LjYwMTU2IDEuODg2NjQgNS4zMTUwMiAxLjYwMDEgNC45NjE1NiAxLjYwMDFaIiBmaWxsPSIjZmZmIi8%2BCjxwYXRoIGQ9Ik00Ljk2MTU2IDEwLjM5OTlIMi4yNDE1NkMxLjg4ODEgMTAuMzk5OSAxLjYwMTU2IDEwLjY4NjQgMS42MDE1NiAxMS4wMzk5VjEzLjc1OTlDMS42MDE1NiAxNC4xMTM0IDEuODg4MSAxNC4zOTk5IDIuMjQxNTYgMTQuMzk5OUg0Ljk2MTU2QzUuMzE1MDIgMTQuMzk5OSA1LjYwMTU2IDE0LjExMzQgNS42MDE1NiAxMy43NTk5VjExLjAzOTlDNS42MDE1NiAxMC42ODY0IDUuMzE1MDIgMTAuMzk5OSA0Ljk2MTU2IDEwLjM5OTlaIiBmaWxsPSIjZmZmIi8%2BCjxwYXRoIGQ9Ik0xMy43NTg0IDEuNjAwMUgxMS4wMzg0QzEwLjY4NSAxLjYwMDEgMTAuMzk4NCAxLjg4NjY0IDEwLjM5ODQgMi4yNDAxVjQuOTYwMUMxMC4zOTg0IDUuMzEzNTYgMTAuNjg1IDUuNjAwMSAxMS4wMzg0IDUuNjAwMUgxMy43NTg0QzE0LjExMTkgNS42MDAxIDE0LjM5ODQgNS4zMTM1NiAxNC4zOTg0IDQuOTYwMVYyLjI0MDFDMTQuMzk4NCAxLjg4NjY0IDE0LjExMTkgMS42MDAxIDEzLjc1ODQgMS42MDAxWiIgZmlsbD0iI2ZmZiIvPgo8cGF0aCBkPSJNNCAxMkwxMiA0TDQgMTJaIiBmaWxsPSIjZmZmIi8%2BCjxwYXRoIGQ9Ik00IDEyTDEyIDQiIHN0cm9rZT0iI2ZmZiIgc3Ryb2tlLXdpZHRoPSIxLjUiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIvPgo8L3N2Zz4K&logoColor=ffffff" alt="zread"/>
  </a>
</div>

<br />

<div align="center">
  <img src="https://github.com/user-attachments/assets/4775958a-5c5b-4184-8cc1-2f32aa693e86" width="100%" alt="Agent-Go Terminal Demo" />
</div>

---

## ⚡ Quick Install

Get the binary and start executing in seconds.

**Standard Install (Pre-built)**
```bash
curl -fsSL https://raw.githubusercontent.com/finettt/agent-go/main/install-agent-go.sh | bash
```

**Rolling Install (Build from Source)**
```bash
curl -fsSL https://raw.githubusercontent.com/finettt/agent-go/main/install-agent-go.sh | bash -s -- --rolling
```

---

##  What is Agent-Go?

**Agent-Go** is a brutalist, command-line AI agent written in Go. It doesn't just chat; it integrates with your shell to **execute commands**, manage files, and automate workflows.

It is a modern rewrite of the original [Agent-C](https://github.com/finettt/agent-c), re-architected for speed, concurrency, and better tooling.

### The Workflow

```mermaid
graph LR
    A[User Input] -->|Context + RAG| B(Agent Core)
    B -->|API Request| C{LLM Brain}
    C -->|Tool Call| D[Executor]
    D -->|Run Command| E[Shell / Filesystem]
    E -->|Output| B
    B -->|Final Response| F[Terminal UI]
    style B fill:#FF69F6,color:#000,stroke:#000,stroke-width:2px
    style C fill:#0F0F0F,color:#FFF,stroke:#FFF
    style D fill:#4AF6F6,color:#000,stroke:#000
```

---

##  Key Features

| Capability | Description |
| :--- | :--- |
| **Native Execution** | Executes shell commands directly (`ls`, `git`, `docker`, etc.) with `&&` chaining support. |
| **Agent Studio** | Create, manage, and spawn specialized sub-agents with the `/agent` command family. |
| **Infinite Memory** | Intelligent token compression ensures your conversation context is never lost. |
| **RAG Engine** | Enable `/rag` to let the agent read your local codebase and documentation for context. |
| **MCP Integration** | Connects to **Model Context Protocol** servers (includes `context7` for docs). |
| **Project Memory** | Persistent Notes and Todo lists that stick with your project across sessions. |
| **Sub-Agents** | Spawns autonomous background threads to handle complex tasks while you keep working. |

---

## Usage

### Interactive Mode
Start the binary to enter the loop.
```bash
agent-go
```

**Workflow Example:**
```text
> Create a Python script that prints "Hello, World!" and then run it
$ echo 'print("Hello, World!")' > hello.py && python hello.py
Created hello.py and executed successfully.
Hello, World!
```

### Single Shot (Headless)
Perfect for CI/CD or scripting.
```bash
agent-go "Create a new directory called 'test-project' and navigate into it"
```

### Agent Studio & Slash Commands
Control the environment with `/` commands.

*   `/agent studio` - Interactively build a new custom agent.
*   `/session save` - Snapshot your current workspace context.
*   `/rag on` - Activate local file awareness.
*   `/stream on` - Toggle matrix-style text streaming.
*   `/shell` - Drop into a standard system shell (bypass AI).

> Type `/help` inside the tool for the full command list.

---

## Configuration

Agent-Go looks for config in this order:
1.  **Environment Variables** (Highest Priority)
2.  **Config File** (`~/.config/agent-go/config.json`)
3.  **Defaults**

**Essential Env Vars:**
```bash
export OPENAI_KEY="sk-..."             # Required
export OPENAI_MODEL="gpt-4-turbo"      # Recommended
export RAG_ENABLED=1                   # Enable local file search
export AUTO_COMPRESS=1                 # Enable infinite memory
```

### Custom Instructions (`AGENTS.md`)
Drop an `AGENTS.md` file in your current directory to give the agent project-specific rules.
```markdown
# AGENTS.md
- Always write comments in Go code.
- Use 'main' branch for git operations.
- Be sarcastic.
```

---

## Architecture

Agent-Go is built for modularity and speed.

*   **`src/executor.go`**: The safety valve. Manages platform-aware shell execution (`sh`/`cmd.exe`).
*   **`src/api.go`**: The bridge. Handles OpenAI-compatible streams and tool definitions.
*   **`src/rag.go`**: The eyes. Indexes local content for retrieval.
*   **`src/subagent.go`**: The crew. Manages background agent lifecycles.

For deep dives, check the [`/docs`](docs/) directory.

---

## Contributing

We welcome pull requests. If you want to add a feature:

1.  Fork it.
2.  Branch it (`git checkout -b feature/cyber-upgrade`).
3.  Commit it.
4.  Push it.

**Development Build:**
```bash
git clone https://github.com/finettt/agent-go.git
make run
```

---

<div align="center">
  <p>Built with 🖤 and Go.</p>
  <img src="https://img.shields.io/badge/Maintained%3F-yes-FF69F6?style=for-the-badge" />
</div>
