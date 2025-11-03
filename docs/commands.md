# Slash Commands

Agent-Go supports a set of slash commands for managing its configuration and features directly from the CLI.

## General Commands

- `/help`: Displays a list of all available commands.
- `/config`: Shows the current configuration, including the model, provider URL, and RAG status.

## Model and Provider

- `/model <model_name>`: Sets the AI model to be used for generating responses.
  - Example: `/model gpt-4-turbo`
- `/provider <api_url>`: Sets the base URL for the API provider.
  - Example: `/provider https://api.openai.com`

## RAG (Retrieval-Augmented Generation)

- `/rag on`: Enables the RAG feature.
- `/rag off`: Disables the RAG feature.
- `/rag path <path>`: Sets the local file system path where documents for RAG are stored.
  - Example: `/rag path /path/to/my/documents`