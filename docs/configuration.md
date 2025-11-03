# Configuration

Agent-Go uses a hierarchical configuration system that loads settings from a JSON file and allows overriding them with environment variables.

## Configuration File

- **Path**: `~/.config/agent-go/config.json`
- **Format**: JSON

The application first loads default settings, then reads the configuration file.

### Example `config.json`

```json
{
  "api_url": "https://api.openai.com",
  "model": "gpt-3.5-turbo",
  "api_key": "your_secret_key_here",
  "rag_path": "/path/to/your/documents",
  "temp": 0.1,
  "max_tokens": 1000,
  "rag_enabled": true,
  "rag_snippets": 5
}
```

## Environment Variables

After loading the configuration file, the application will override any existing settings with values from environment variables. This is useful for CI/CD environments or temporary adjustments.

- `OPENAI_KEY`: Your OpenAI API key. **(Required if not in config file)**
- `OPENAI_BASE`: The base URL for the API.
- `OPENAI_MODEL`: The model to use.
- `RAG_PATH`: The path to your RAG documents.
- `RAG_ENABLED`: Set to `1` to enable RAG.
- `RAG_SNIPPETS`: The number of snippets to retrieve.

## Priority

1.  **Environment Variables** (highest priority)
2.  **Configuration File** (`~/.config/agent-go/config.json`)
3.  **Default Values** (lowest priority)