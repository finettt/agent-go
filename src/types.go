package main

type Message struct {
	Role       string     `json:"role"`
	Content    *string    `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type Config struct {
	APIURL                string  `json:"api_url"`
	Model                 string  `json:"model"`
	APIKey                string  `json:"api_key"`
	RAGPath               string  `json:"rag_path"`
	Temp                  float32 `json:"temp"`
	MaxTokens             int     `json:"max_tokens"`
	RAGEnabled            bool    `json:"rag_enabled"`
	RAGSnippets           int     `json:"rag_snippets"`
	AutoCompress          bool    `json:"auto_compress"`
	AutoCompressThreshold int     `json:"auto_compress_threshold"`
	ModelContextLength    int     `json:"model_context_length"`
	Stream                bool    `json:"stream"`
	SubagentsEnabled      bool    `json:"subagents_enabled"`
}

type Agent struct {
	ID       string
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
	Usage   Usage    `json:"usage"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choice struct {
	Message Message `json:"message"`
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

type SubAgentTask struct {
	Task string `json:"task"`
}

// Streaming response types
type StreamChunk struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

type StreamChoice struct {
	Index        int     `json:"index"`
	Delta        Delta   `json:"delta"`
	FinishReason *string `json:"finish_reason"`
}

type Delta struct {
	Role      *string    `json:"role,omitempty"`
	Content   *string    `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}
