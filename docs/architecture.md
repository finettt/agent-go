# Architecture

This document outlines the architecture of Agent-Go, a command-line AI assistant.

## Flow Diagram

The application follows a clear, sequential flow from launch to command execution.

```mermaid
graph TD
    A[Launch] --> B{Config Exists?};
    B -- No --> C[Interactive Setup Screen];
    C --> D[Save to ~/.config/agent-go/config.json];
    B -- Yes --> E[Read Config];
    D --> E;
    E --> F[Initialize Readline with History];
    F --> G[Main Input Loop];
    G --> H{Is Input a Slash Command?};
    H -- Yes --> I[Slash Command Handler];
    I --> J{Command Requires API?};
    J -- Yes --> K[Send API Request];
    J -- No --> G;
    H -- No --> L[Process as Standard Prompt];
    L --> K;
    K --> M[Print Response];
    M --> G;
```

## Core Components

- **main.go**: The application's entry point, responsible for initialization, the main CLI loop, and graceful shutdown.
- **config.go**: Manages loading and saving application settings from a JSON file and environment variables.
- **api.go**: Handles all communication with the OpenAI API, including sending requests and processing responses.
- **executor.go**: Executes shell commands requested by the AI model.
- **rag.go**: Implements the Retrieval-Augmented Generation (RAG) feature by searching local files for context.
- **types.go**: Defines all data structures used throughout the application.