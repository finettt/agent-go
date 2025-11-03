package main

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Config struct {
	APIURL      string  `json:"api_url"`
	Model       string  `json:"model"`
	APIKey      string  `json:"api_key"`
	RAGPath     string  `json:"rag_path"`
	Temp        float32 `json:"temp"`
	MaxTokens   int     `json:"max_tokens"`
	RAGEnabled  bool    `json:"rag_enabled"`
	RAGSnippets int     `json:"rag_snippets"`
}

type Agent struct {
	Messages []Message
	MsgCount int
}

type APIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
	Stream      bool      `json:"stream"`
	Tools       []Tool    `json:"tools"`
	ToolChoice  string    `json:"tool_choice"`
}

type Tool struct {
	Type     string             `json:"type"`
	Function FunctionDefinition `json:"function"`
}

type FunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

type APIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message ResponseMessage `json:"message"`
}

type ResponseMessage struct {
	Role      string     `json:"role"`
	Content   *string    `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type CommandArgs struct {
	Command string `json:"command"`
}
