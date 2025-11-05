# Configuration

Agent-Go uses a hierarchical configuration system that loads settings from multiple sources, allowing you to customize the application behavior for different environments and use cases.

## Configuration Sources

Agent-Go uses a three-tier configuration system with the following priority (highest to lowest):

1. **Environment Variables** - Override all other settings
2. **Configuration File** - Persistent settings stored in JSON format
3. **Default Values** - Built-in defaults for all settings

## Configuration File

### Location and Format

- **Path**: `~/.config/agent-go/config.json`
- **Format**: JSON
- **Permissions**: Secure file permissions (600) to protect sensitive data

The configuration file is automatically created on first run or when settings are modified through slash commands.

### Complete Configuration Example

```json
{
  "api_url": "https://api.openai.com",
  "model": "gpt-4-turbo",
  "api_key": "sk-your-secret-api-key-here",
  "temperature": 0.1,
  "max_tokens": 1000,
  "rag_enabled": true,
  "rag_path": "/home/user/documents",
  "rag_snippets": 5,
  "message_history_limit": 20
}
```

### Configuration Parameters

#### API Configuration

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `api_url` | string | `"https://api.openai.com"` | Base URL for the AI API provider |
| `api_key` | string | `""` | API key for authentication (required) |
| `model` | string | `"gpt-3.5-turbo"` | AI model to use for responses |
| `temperature` | float | `0.1` | Controls randomness (0.0-1.0, lower = more deterministic) |
| `max_tokens` | int | `1000` | Maximum tokens per response |

#### RAG Configuration

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `rag_enabled` | bool | `false` | Enable/disable Retrieval-Augmented Generation |
| `rag_path` | string | `""` | Path to local documents for RAG |
| `rag_snippets` | int | `5` | Number of document snippets to include in context |

#### Application Configuration

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `message_history_limit` | int | `20` | Maximum number of messages to retain in conversation history |

## Environment Variables

Environment variables provide a way to override configuration settings without modifying the configuration file. This is particularly useful for:

- CI/CD pipelines
- Docker containers
- Temporary testing
- Security-sensitive environments

### Available Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `OPENAI_KEY` | API key for authentication | `sk-proj-abc123...` |
| `OPENAI_BASE` | Base URL for API provider | `https://api.openai.com` |
| `OPENAI_MODEL` | AI model to use | `gpt-4-turbo` |
| `RAG_PATH` | Path to RAG documents | `/home/user/documents` |
| `RAG_ENABLED` | Enable RAG feature (`1`=enabled, `0`=disabled) | `1` |
| `RAG_SNIPPETS` | Number of RAG snippets | `5` |

### Environment Variable Examples

#### Basic Setup

```bash
export OPENAI_KEY="sk-proj-your-api-key-here"
export OPENAI_MODEL="gpt-4-turbo"
```

#### Development Environment

```bash
export OPENAI_KEY="dev-api-key"
export OPENAI_BASE="http://localhost:8080/v1"
export RAG_ENABLED=1
export RAG_PATH="./project-docs"
```

#### Production Environment

```bash
export OPENAI_KEY="${OPENAI_API_KEY}"
export OPENAI_MODEL="gpt-4-turbo"
export RAG_ENABLED=0
```

## Configuration Priority and Merging

The configuration system merges settings from all sources with the following precedence:

1. **Environment Variables** (highest priority)
2. **Configuration File** (`~/.config/agent-go/config.json`)
3. **Default Values** (lowest priority)

### Example Priority Scenario

**Default Configuration:**

```json
{
  "api_url": "https://api.openai.com",
  "model": "gpt-3.5-turbo",
  "temperature": 0.1
}
```

**Configuration File:**

```json
{
  "model": "gpt-4-turbo",
  "rag_enabled": true
}
```

**Environment Variables:**

```bash
export OPENAI_BASE="https://api.anthropic.com"
export RAG_SNIPPETS=10
```

**Final Merged Configuration:**

```json
{
  "api_url": "https://api.openai.com",          // from config file
  "model": "gpt-4-turbo",                       // from config file
  "temperature": 0.1,                          // from config file
  "rag_enabled": true,                         // from config file
  "rag_snippets": 10,                          // from environment variable
  "api_key": ""                                // from environment variable (if set)
}
```

## Configuration Scenarios

### Scenario 1: Development Setup

For local development with custom documents:

```json
{
  "api_url": "http://localhost:8080/v1",
  "model": "gpt-4-turbo",
  "api_key": "dev-key",
  "rag_enabled": true,
  "rag_path": "./docs",
  "temperature": 0.2
}
```

### Scenario 2: Production Deployment

For production with minimal features:

```json
{
  "api_url": "https://api.openai.com",
  "model": "gpt-4-turbo",
  "api_key": "${OPENAI_API_KEY}",
  "rag_enabled": false,
  "temperature": 0.0
}
```

### Scenario 3: Multi-Environment Setup

Use environment variables for different environments:

```bash
# Development
export OPENAI_BASE="http://localhost:8080/v1"
export OPENAI_MODEL="gpt-4-turbo"
export RAG_ENABLED=1

# Production
export OPENAI_BASE="https://api.openai.com"
export OPENAI_MODEL="gpt-4"
export RAG_ENABLED=0
```

## Configuration Management

### First-Time Setup

When you run Agent-Go for the first time:

1. The application detects missing configuration
2. Launches an interactive setup wizard
3. Prompts for required API key
4. Saves configuration to `~/.config/agent-go/config.json`

### Modifying Configuration

#### Using Slash Commands (Interactive)

```bash
> /model gpt-4-turbo
Model set to: gpt-4-turbo

> /rag on
RAG enabled.

> /rag path /home/user/documents
RAG path set to: /home/user/documents
```

#### Manual Configuration File Editing

```bash
# Create/edit configuration file
nano ~/.config/agent-go/config.json
```

#### Environment Variables for CI/CD

```bash
# GitHub Actions example
env:
  OPENAI_KEY: ${{ secrets.OPENAI_API_KEY }}
  OPENAI_MODEL: gpt-4-turbo
  RAG_ENABLED: 0
```

## Security Considerations

### API Key Protection

- **Never commit API keys to version control**
- **Use environment variables for sensitive data**
- **Set proper file permissions**:

  ```bash
  chmod 600 ~/.config/agent-go/config.json
  ```

- **Consider using secret management tools** for production

### Configuration File Security

```bash
# Secure configuration file
mkdir -p ~/.config/agent-go
touch ~/.config/agent-go/config.json
chmod 600 ~/.config/agent-go/config.json
```

## Troubleshooting

### Common Configuration Issues

**Missing API Key**

```
Error: OpenAI API key is not set
Solution: Set OPENAI_KEY environment variable or run interactive setup
```

**Invalid API URL**

```
Error: could not connect to API
Solution: Verify OPENAI_BASE URL is correct and accessible
```

**RAG Path Issues**

```
Error: cannot access RAG documents
Solution: Ensure RAG_PATH exists and is readable
```

### Configuration Validation

Agent-Go validates configuration on startup:

1. **Required fields**: `api_key` must be set
2. **URL validation**: `api_url` must be a valid URL
3. **Path validation**: `rag_path` must exist if RAG is enabled
4. **Range validation**: `temperature` must be 0.0-1.0
5. **Type validation**: All values must match expected types

### Debug Mode

For troubleshooting configuration issues:

```bash
# Enable debug logging
export DEBUG=1
./agent-go
```

## Best Practices

### 1. Environment-Specific Configurations

Use different configuration files for different environments:

```bash
# Development
export CONFIG_FILE=~/.config/agent-go/dev.json

# Production
export CONFIG_FILE=~/.config/agent-go/prod.json
```

### 2. Secret Management

Use secret management tools for production:

```bash
# Using Docker secrets
docker run --env-file /run/secrets/openai_key ...

# Using Kubernetes
env:
  - name: OPENAI_KEY
    valueFrom:
      secretKeyRef:
        name: ai-secrets
        key: openai-key
```

### 3. Configuration Backup

Regularly backup your configuration:

```bash
# Backup configuration
cp ~/.config/agent-go/config.json ~/.config/agent-go/config.json.backup

# Restore configuration
cp ~/.config/agent-go/config.json.backup ~/.config/agent-go/config.json
```

### 4. Version Control

Exclude sensitive configuration from version control:

```gitignore
# .gitignore
~/.config/agent-go/config.json
*.key
*.secret
```

## Migration Guide

### Upgrading from Previous Versions

1. **Backup existing configuration**:

   ```bash
   cp ~/.config/agent-go/config.json ~/.config/agent-go/config.json.backup
   ```

2. **Review new configuration options**:
   - Check for new parameters in the documentation
   - Update your configuration file if needed

3. **Test with new version**:

   ```bash
   make test
   ./agent-go --help
   ```

### Configuration File Format Changes

If you're upgrading from an older version, the configuration file format may have changed. Review the example configuration above and update your file accordingly.
