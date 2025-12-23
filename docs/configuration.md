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
  "temp": 0.1,
  "max_tokens": 1000,
  "rag_enabled": true,
  "rag_path": "/home/user/documents",
  "rag_snippets": 5,
  "auto_compress": true,
  "auto_compress_threshold": 20,
  "model_context_length": 131072,
  "stream": false,
  "subagents_enabled": true,
  "mcp_servers": {
    "context7": {
      "name": "context7",
      "command": "npx -y @upstash/context7-mcp"
    }
  }
}
```

### Configuration Parameters

#### API Configuration

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `api_url` | string | `"https://api.openai.com"` | Base URL for the AI API provider |
| `api_key` | string | `""` | API key for authentication (required) |
| `model` | string | `"gpt-3.5-turbo"` | AI model to use for responses |
| `temp` | float | `0.1` | Controls randomness (0.0-1.0, lower = more deterministic) |
| `max_tokens` | int | `1000` | Maximum tokens per response |

#### RAG Configuration

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `rag_enabled` | bool | `false` | Enable/disable Retrieval-Augmented Generation |
| `rag_path` | string | `""` | Path to local documents for RAG |
| `rag_snippets` | int | `5` | Number of document snippets to include in context |

#### Context Management Configuration

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `auto_compress` | bool | `true` | Enable automatic context compression |
| `auto_compress_threshold` | int | `20` | Threshold for auto compression (percentage of context length) |
| `model_context_length` | int | `131072` | Maximum context length for the AI model |

#### Streaming and Sub-agent Configuration

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `stream` | bool | `false` | Enable/disable streaming mode for real-time responses |
| `subagents_enabled` | bool | `true` | Enable/disable sub-agent spawning capability |
| `execution_mode` | string | `"ask"` | Execution mode: `"ask"` (confirm commands) or `"yolo"` (auto-execute) |
| `operation_mode` | string | `"build"` | Operation mode: `"build"` (execute commands) or `"plan"` (plan only) |

#### MCP Server Configuration

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `mcp_servers` | object | `{"context7": {...}}` | Map of MCP server configurations |

**MCP Server Object Structure:**
```json
{
  "server_name": {
    "name": "server_name",
    "command": "command to launch server"
  }
}
```

**Default MCP Server:**
The `context7` server is automatically configured for accessing up-to-date library documentation.

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
| `AUTO_COMPRESS` | Enable auto context compression (`1`=enabled, `0`=disabled) | `1` |
| `AUTO_COMPRESS_THRESHOLD` | Threshold for auto compression | `20` |
| `MODEL_CONTEXT_LENGTH` | Model context length | `131072` |
| `STREAM_ENABLED` | Enable streaming mode (`1`=enabled, `0`=disabled) | `1` or `true` |
| `SUBAGENTS_ENABLED` | Enable sub-agent spawning (`1`=enabled, `0`=disabled) | `1` or `true` |
| `EXECUTION_MODE` | Set execution mode | `"ask"` or `"yolo"` |
| `OPERATION_MODE` | Set operation mode | `"build"` or `"plan"` |

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
export AUTO_COMPRESS=1
export AUTO_COMPRESS_THRESHOLD=15
```

#### Production Environment

```bash
export OPENAI_KEY="${OPENAI_API_KEY}"
export OPENAI_MODEL="gpt-4-turbo"
export RAG_ENABLED=0
export AUTO_COMPRESS=1
export MODEL_CONTEXT_LENGTH=16384
```

#### High-Performance Configuration

```bash
export OPENAI_KEY="your-api-key"
export OPENAI_MODEL="gpt-4-turbo"
export RAG_ENABLED=1
export RAG_PATH="/optimized/documents"
export RAG_SNIPPETS=8
export AUTO_COMPRESS=1
export AUTO_COMPRESS_THRESHOLD=25
export MODEL_CONTEXT_LENGTH=131072
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
  "auto_compress": true,                       // from config file
  "auto_compress_threshold": 20,               // from config file
  "model_context_length": 131072,              // from config file
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
  "temperature": 0.2,
  "auto_compress": true,
  "auto_compress_threshold": 15,
  "model_context_length": 32768
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
  "temperature": 0.0,
  "auto_compress": true,
  "auto_compress_threshold": 25,
  "model_context_length": 16384
}
```

### Scenario 3: Multi-Environment Setup

Use environment variables for different environments:

```bash
# Development
export OPENAI_BASE="https://api.openai.com"
export OPENAI_MODEL="gpt-4-turbo"
export RAG_ENABLED=1
export RAG_PATH="./project-docs"
export AUTO_COMPRESS=1
export AUTO_COMPRESS_THRESHOLD=15

# Production
export OPENAI_BASE="https://api.openai.com"
export OPENAI_MODEL="gpt-4"
export RAG_ENABLED=0
export AUTO_COMPRESS=1
export MODEL_CONTEXT_LENGTH=131072
```

### Scenario 4: RAG-Optimized Setup

For heavy RAG usage with large document collections:

```json
{
  "api_url": "https://api.openai.com",
  "model": "gpt-4-turbo",
  "api_key": "your-api-key",
  "rag_enabled": true,
  "rag_path": "/large/documents/collection",
  "rag_snippets": 10,
  "temperature": 0.1,
  "auto_compress": true,
  "auto_compress_threshold": 30,
  "model_context_length": 131072
}
```

## Configuration Management

### First-Time Setup

When you run Agent-Go for the first time:

1. The application detects missing configuration
2. Launches an interactive setup wizard
3. Prompts for required API key
4. Saves configuration to `~/.config/agent-go/config.json`
5. Provides feedback on successful configuration

### Modifying Configuration

#### Using Slash Commands (Interactive)

```bash
> /model gpt-4-turbo
Model set to: gpt-4-turbo

> /rag on
RAG enabled.

> /rag path /home/user/documents
RAG path set to: /home/user/documents

> /contextlength 131072
Model context length set to: 131072

> /stream on
Streaming enabled.

> /subagents on
Sub-agent spawning enabled.
```

#### Managing MCP Servers (Interactive)

```bash
> /mcp list
Configured MCP servers:
- context7: npx -y @upstash/context7-mcp

> /mcp add time uvx mcp-server-time
MCP server 'time' added.

> /mcp add weather npx -y @weather/mcp-server
MCP server 'weather' added.

> /mcp remove weather
MCP server 'weather' removed.
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
  AUTO_COMPRESS: 1
  MODEL_CONTEXT_LENGTH: 131072
```

#### Docker Environment Setup

The Dockerfile is configured to create a `/workspace` volume where you can mount your current directory:

```bash
# Basic usage with workspace mounted
docker run -it -v $(pwd):/workspace agent-go

# With full environment configuration
docker run -it \
  -v $(pwd):/workspace \
  -e OPENAI_KEY="your-api-key" \
  -e OPENAI_MODEL="gpt-4-turbo" \
  -e RAG_ENABLED=1 \
  -e RAG_PATH="/workspace" \
  -e AUTO_COMPRESS=1 \
  agent-go

# With separate documents directory
docker run -it \
  -v $(pwd):/workspace \
  -v /path/to/documents:/documents \
  -e OPENAI_KEY="your-api-key" \
  -e RAG_ENABLED=1 \
  -e RAG_PATH="/documents" \
  agent-go

# With persistent configuration
docker run -it \
  -v $(pwd):/workspace \
  -v ~/.config/agent-go:/home/appuser/.config/agent-go \
  -e OPENAI_KEY="your-api-key" \
  agent-go

# With MCP server configuration
docker run -it \
  -v $(pwd):/workspace \
  -e OPENAI_KEY="your-api-key" \
  agent-go /mcp add custom-server "npx -y @custom/mcp-server"
```

**Docker Volume Configuration:**
- `/workspace` - Mount your current directory here for file access
- `/home/appuser/.config/agent-go` - Configuration directory (can be mounted for persistence)
- Additional volumes can be mounted for RAG documents or other data
- The container runs as non-root user `appuser` for security
- Ensure mounted directories have appropriate permissions

## Security Considerations

### API Key Protection

- **Never commit API keys to version control**
- **Use environment variables for sensitive data**
- **Set proper file permissions**:

  ```bash
  chmod 600 ~/.config/agent-go/config.json
  ```

- **Consider using secret management tools** for production
- **Use encrypted configuration files** for sensitive environments

### Configuration File Security

```bash
# Secure configuration file
mkdir -p ~/.config/agent-go
touch ~/.config/agent-go/config.json
chmod 600 ~/.config/agent-go/config.json
```

### Environment Variable Security

```bash
# Use .env files for development
echo "OPENAI_KEY=your-api-key" > .env
echo "OPENAI_MODEL=gpt-4-turbo" >> .env
# Add .env to .gitignore

# Use secure shell sessions for sensitive operations
read -s -p "Enter API key: " api_key
export OPENAI_KEY="$api_key"
```

## Troubleshooting

### Common Configuration Issues

### Missing API Key

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

**Context Compression Issues**

```
Error: context compression failed
Solution: Check API connectivity and ensure you have sufficient messages to compress
```

**Model Context Length Issues**

```
Error: invalid context length
Solution: Set MODEL_CONTEXT_LENGTH to a positive integer matching your model's capabilities
```

**Auto-compression Threshold Issues**

```
Error: invalid auto-compress threshold
Solution: Set AUTO_COMPRESS_THRESHOLD to a positive integer (typically 10-50)
```

### Configuration Validation

Agent-Go validates configuration on startup:

1. **Required fields**: `api_key` must be set
2. **URL validation**: `api_url` must be a valid URL
3. **Path validation**: `rag_path` must exist if RAG is enabled
4. **Range validation**: `temperature` must be 0.0-1.0
5. **Type validation**: All values must match expected types
6. **Context length validation**: `model_context_length` must be positive
7. **Auto-compression validation**: `auto_compress_threshold` must be positive
8. **RAG snippets validation**: `rag_snippets` must be positive

### Debug Mode

For troubleshooting configuration issues:

```bash
# Enable debug logging
export DEBUG=1
./agent-go

# Or use verbose mode
./agent-go -v
```

### Configuration Testing

```bash
# Test configuration without running the full application
go run src/main.go --test-config

# Validate configuration file
cat ~/.config/agent-go/config.json | jq .
```

## Advanced Configuration Management

### Environment-Specific Configurations

Use different configuration files for different environments:

```bash
# Development
export CONFIG_FILE=~/.config/agent-go/dev.json

# Production
export CONFIG_FILE=~/.config/agent-go/prod.json

# Override with environment variables
export OPENAI_KEY="dev-key" && ./agent-go
```

### MCP Server Management

Configure MCP servers for extended functionality:

```bash
# Add MCP servers programmatically
export MCP_SERVERS='{"time":{"name":"time","command":"uvx mcp-server-time"}}'

# Or via configuration file
cat > ~/.config/agent-go/config.json <<EOF
{
  "api_key": "${OPENAI_KEY}",
  "mcp_servers": {
    "context7": {
      "name": "context7",
      "command": "npx -y @upstash/context7-mcp"
    },
    "time": {
      "name": "time",
      "command": "uvx mcp-server-time"
    }
  }
}
EOF
```

**Common MCP Servers:**
- **context7**: Library documentation (default)
- **time**: Time and timezone utilities
- **filesystem**: File system operations
- **database**: Database query tools
- **weather**: Weather information

### Todo List Storage

Todo lists are stored per-agent:

```bash
# Todo list location
~/.config/agent-go/todos/{agent-id}.json

# Main agent todos
~/.config/agent-go/todos/main.json

# Sub-agent todos (UUID-based)
~/.config/agent-go/todos/550e8400-e29b-41d4-a716-446655440000.json
```

**Todo List Structure:**
```json
{
  "agent_id": "main",
  "todos": [
    {
      "id": 1,
      "task": "Set up development environment",
      "status": "completed"
    },
    {
      "id": 2,
      "task": "Write documentation",
      "status": "in-progress"
    }
  ],
  "next_id": 3
}
```

### Secret Management

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

# Using HashiCorp Vault
vault read -field=openai_key secret/agent-go
```

### Configuration Backup and Recovery

Regularly backup your configuration:

```bash
# Backup configuration
cp ~/.config/agent-go/config.json ~/.config/agent-go/config.json.backup

# Restore configuration
cp ~/.config/agent-go/config.json.backup ~/.config/agent-go/config.json

# Compress old configurations
find ~/.config/agent-go -name "config.json.*" -mtime +30 -delete
```

### Version Control Integration

Exclude sensitive configuration from version control:

```gitignore
# .gitignore
~/.config/agent-go/config.json
*.key
*.secret
.env
.env.local
```

### Configuration Templates

Create configuration templates for different use cases:

```bash
# Template for new users
cat > ~/.config/agent-go/config.template.json << EOF
{
  "api_url": "https://api.openai.com",
  "model": "gpt-3.5-turbo",
  "temperature": 0.1,
  "max_tokens": 1000,
  "rag_enabled": false,
  "auto_compress": true,
  "auto_compress_threshold": 20,
  "model_context_length": 131072
}
EOF

# Copy template for new user
cp ~/.config/agent-go/config.template.json ~/.config/agent-go/config.json
```

### Configuration Monitoring and Logging

Monitor configuration changes and usage:

```bash
# Log configuration changes
echo "$(date): Configuration modified" >> ~/.config/agent-go/config.log

# Monitor API usage
grep "API request" ~/.config/agent-go/agent.log | tail -10
```

## Configuration Best Practices

### Security Best Practices

1. **Never hardcode API keys** in configuration files
2. **Use environment variables** for sensitive data in production
3. **Set proper file permissions** (600) for configuration files
4. **Use encrypted configuration files** for sensitive environments
5. **Implement audit logging** for configuration changes
6. **Regularly rotate API keys** and update configuration

### Performance Optimization

1. **Use appropriate context lengths** for your model
2. **Optimize RAG search parameters** for your document collection
3. **Enable auto-compression** to manage long conversations
4. **Set reasonable token limits** to control API costs
5. **Use caching** where appropriate for frequently accessed data

### Development vs Production

**Development Configuration:**
- More verbose logging
- Debug mode enabled
- Local API endpoints
- Larger context limits
- Experimental features enabled

**Production Configuration:**
- Minimal logging
- Strict security settings
- Production API endpoints
- Optimized context limits
- Experimental features disabled

### Configuration Validation

Agent-Go performs comprehensive validation on startup:

```bash
# Validate configuration file structure
cat ~/.config/agent-go/config.json | jq .

# Test API connectivity
curl -H "Authorization: Bearer $OPENAI_KEY" $OPENAI_BASE/v1/models

# Validate RAG path
ls -la $RAG_PATH 2>/dev/null && echo "RAG path is accessible"
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

4. **Migrate configuration**:

   ```bash
   # Check for deprecated options
   grep -i "message_history_limit" ~/.config/agent-go/config.json && echo "Warning: message_history_limit is deprecated, use model_context_length instead"
   ```

### Configuration File Format Changes

If you're upgrading from an older version, the configuration file format may have changed. Review the example configuration above and update your file accordingly.

**Key Changes in Recent Versions:**

- **Removed**: `message_history_limit` - replaced with `model_context_length`
- **Added**: `auto_compress` - automatic context compression
- **Added**: `auto_compress_threshold` - threshold for auto compression
- **Added**: `model_context_length` - configurable context length

**Migration Guide:**

```bash
# Old configuration format
{
  "message_history_limit": 20
}

# New configuration format
{
  "model_context_length": 131072,
  "auto_compress": true,
  "auto_compress_threshold": 20
}
```
