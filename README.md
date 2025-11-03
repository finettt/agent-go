# Agent-Go

A lightweight AI agent written in Go that communicates with the OpenAI API and executes shell commands. This is a rewrite of the original [Agent-C](https://github.com/finettt/agent-c) project.

## Features

- **Tool Calling**: Execute shell commands directly through AI responses.
- **Conversation Memory**: Manages a sliding window of the last 20 messages.
- **Cross-Platform**: Works on macOS, Linux, and Windows.
- **RAG (Retrieval-Augmented Generation)**: Searches local files to provide context-aware responses.

## Quick Start

### Prerequisites

- Go 1.18+
- An OpenAI API key

### Build

```bash
make
```

### Setup

Set your OpenAI API key and other configurations as environment variables:

```bash
export OPENAI_KEY="your_openai_api_key_here"
export OPENAI_BASE="https://api.openai.com"  # Optional, defaults to OpenAI
export OPENAI_MODEL="gpt-3.5-turbo"             # Optional
export RAG_PATH="/path/to/your/documents"       # Optional, for RAG
export RAG_ENABLED=1                            # Optional, 1 to enable RAG
export RAG_SNIPPETS=5                           # Optional, number of snippets
```

### Run

```bash
./agent-go
```

## License

CC0 - "No Rights Reserved"