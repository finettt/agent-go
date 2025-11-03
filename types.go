package main

// Message представляет собой одно сообщение в истории диалога.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Config содержит конфигурацию приложения.
type Config struct {
	APIURL       string
	Model        string
	APIKey       string
	RAGPath      string
	Temp         float32
	MaxTokens    int
	RAGEnabled   bool
	RAGSnippets  int
}

// Agent представляет состояние ИИ-агента.
type Agent struct {
	Messages []Message
	MsgCount int
}

// APIRequest представляет тело запроса к OpenAI API.
type APIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
	Stream      bool      `json:"stream"`
	Tools       []Tool    `json:"tools"`
	ToolChoice  string    `json:"tool_choice"`
}

// Tool представляет инструмент, доступный модели.
type Tool struct {
	Type     string             `json:"type"`
	Function FunctionDefinition `json:"function"`
}

// FunctionDefinition описывает функцию.
type FunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

// APIResponse представляет тело ответа от OpenAI API.
type APIResponse struct {
	Choices []Choice `json:"choices"`
}

// Choice представляет один из вариантов ответа.
type Choice struct {
	Message ResponseMessage `json:"message"`
}

// ResponseMessage представляет сообщение в ответе API.
type ResponseMessage struct {
	Role      string     `json:"role"`
	Content   *string    `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// ToolCall представляет вызов инструмента, запрошенный моделью.
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// FunctionCall представляет конкретный вызов функции.
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// CommandArgs используется для разбора аргументов команды.
type CommandArgs struct {
	Command string `json:"command"`
}