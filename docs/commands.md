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
  /shell             - Enter shell mode
  /quit              - Exit the application
```

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
- Auto-completion is available for model names

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

- Requires a valid RAG path to be set
- May increase response time as documents are searched
- Provides more context-aware responses for document-related queries

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

### Getting Help

If you encounter issues not covered here:

1. Use `/help` to see all available commands
2. Check the main documentation in the `/docs` directory
3. Review the [Architecture](architecture.md) and [Configuration](configuration.md) documents
4. For bugs or feature requests, check the GitHub repository

### `/compress`

Сжимает текущий контекст беседы и начинает новый чат с сжатым резюме в качестве системного сообщения.

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

- Требует наличия хотя бы одного сообщения в текущем чате
- Сохраняет ключевые детали и контекст предыдущего разговора
- Полностью очищает текущую историю сообщений
- Полезно для очень длинных бесед, чтобы избежать токен-лимитов
- Использует тот же API для интеллектуального сжатия контекста
- Автоматически включает AGENTS.md в новый системный промпт, если он существует
