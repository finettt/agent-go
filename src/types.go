package main

import (
	"bytes"
	"os"
	"time"
)

type Message struct {
	Role             string     `json:"role"`
	Content          *string    `json:"content,omitempty"`
	ReasoningContent *string    `json:"reasoning_content,omitempty"`
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID       string     `json:"tool_call_id,omitempty"`
}

type Config struct {
	APIURL                string               `json:"api_url"`
	Model                 string               `json:"model"`
	MiniModel             string               `json:"mini_model"`
	APIKey                string               `json:"api_key"`
	RAGPath               string               `json:"rag_path"`
	Temp                  float32              `json:"temp"`
	MaxTokens             int                  `json:"max_tokens"`
	RAGEnabled            bool                 `json:"rag_enabled"`
	RAGSnippets           int                  `json:"rag_snippets"`
	AutoCompress          bool                 `json:"auto_compress"`
	AutoCompressThreshold int                  `json:"auto_compress_threshold"`
	ModelContextLength    int                  `json:"model_context_length"`
	SubagentsEnabled      bool                 `json:"subagents_enabled"`
	SubAgentVerboseMode   int                  `json:"subagent_verbose_mode"`
	ExecutionMode         ExecuteMode          `json:"execution_mode"`
	OperationMode         OperationMode        `json:"operation_mode"`
	MCPs                  map[string]MCPServer `json:"mcp_servers"`
	Skills                []Skill              `json:"skills"`
	UsageVerboseMode      int                  `json:"usage_verbose_mode"`
}

const (
	UsageSilent   = 1
	UsageBasic    = 2
	UsageDetailed = 3
)

// ExecutionMode controls command confirmation behavior
// This is INDEPENDENT of agent selection
type ExecuteMode string

const (
	// Ask is the default execution mode, which asks the user for confirmation before executing a command.
	Ask ExecuteMode = "ask"
	// YOLO is the execution mode that executes a command without asking for confirmation.
	YOLO ExecuteMode = "yolo"
)

// DEPRECATED: OperationMode is deprecated. Use agent definitions (plan.json/build.json) instead.
// This type remains for backward compatibility and will be removed in v2.0.0.
type OperationMode string

const (
	// Build is the default operation mode, which allows command execution.
	// DEPRECATED: Use the 'build' agent instead.
	Build OperationMode = "build"
	// Plan is the operation mode that blocks all command execution and focuses on planning.
	// DEPRECATED: Use the 'plan' agent instead.
	Plan OperationMode = "plan"
)

// Agent represents an AI agent with its properties and message history
type Agent struct {
	ID           string    // Unique identifier for the agent
	Messages     []Message // List of messages in the conversation
	AgentDefName string    `json:"agent_def_name,omitempty"` // Name of the agent definition in use
}

type APIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
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
	Command    string `json:"command"`
	Background bool   `json:"background,omitempty"`
}

type SubAgentTask struct {
	Task  string `json:"task"`
	Agent string `json:"agent,omitempty"` // Optional task-specific agent name (e.g., "build", "plan", or a custom agent). Defaults to "build".
	Model string `json:"model,omitempty"` // Optional model selection ("main" or "mini")
}

type UseMCPToolArgs struct {
	ServerName string                 `json:"server_name"`
	ToolName   string                 `json:"tool_name"`
	Arguments  map[string]interface{} `json:"arguments"`
}

// MCPServer defines the configuration for a single MCP server
type MCPServer struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

// Skill represents a custom tool backed by a script
type Skill struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"` // JSON Schema for parameters
	Command     string                 `json:"command"`    // The script/command to execute
}

type KillBackgroundCommandArgs struct {
	PID int `json:"pid"`
}

type GetBackgroundLogsArgs struct {
	PID int `json:"pid"`
}

type SwitchOperationModeArgs struct {
	Mode string `json:"mode"`
}

// Terminal session management types

type TerminalSession struct {
	ID           string        `json:"id"`
	PID          int           `json:"pid"`
	PTY          *os.File      `json:"-"`
	StdoutBuffer *bytes.Buffer `json:"-"`
	InputChan    chan string   `json:"-"`
	DoneChan     chan struct{} `json:"-"`
	StartTime    time.Time     `json:"start_time"`
}

type TerminalInputArgs struct {
	SessionID string `json:"session_id"`
	Input     string `json:"input"`
}

type TerminalReadArgs struct {
	SessionID string `json:"session_id"`
	Bytes     int    `json:"bytes,omitempty"`
	ReadAll   bool   `json:"read_all,omitempty"`
}

type TerminalCloseArgs struct {
	SessionID string `json:"session_id"`
}

// KeyMappings maps human-readable keys to terminal escape sequences
var KeyMappings = map[string]string{
	// Control codes (Ctrl+A through Ctrl+Z)
	"Ctrl+A": "\x01",
	"Ctrl+B": "\x02",
	"Ctrl+C": "\x03",
	"Ctrl+D": "\x04",
	"Ctrl+E": "\x05",
	"Ctrl+F": "\x06",
	"Ctrl+G": "\x07",
	"Ctrl+H": "\x08",
	"Ctrl+I": "\x09",
	"Ctrl+J": "\x0a",
	"Ctrl+K": "\x0b",
	"Ctrl+L": "\x0c",
	"Ctrl+M": "\x0d",
	"Ctrl+N": "\x0e",
	"Ctrl+O": "\x0f",
	"Ctrl+P": "\x10",
	"Ctrl+Q": "\x11",
	"Ctrl+R": "\x12",
	"Ctrl+S": "\x13",
	"Ctrl+T": "\x14",
	"Ctrl+U": "\x15",
	"Ctrl+V": "\x16",
	"Ctrl+W": "\x17",
	"Ctrl+X": "\x18",
	"Ctrl+Y": "\x19",
	"Ctrl+Z": "\x1a",

	// Special keys
	"Enter":     "\n",
	"Tab":       "\t",
	"Escape":    "\x1b",
	"Backspace": "\x7f",

	// Arrow keys
	"ArrowUp":    "\x1b[A",
	"ArrowDown":  "\x1b[B",
	"ArrowRight": "\x1b[C",
	"ArrowLeft":  "\x1b[D",

	// Function keys
	"F1":  "\x1b[11~",
	"F2":  "\x1b[12~",
	"F3":  "\x1b[13~",
	"F4":  "\x1b[14~",
	"F5":  "\x1b[15~",
	"F6":  "\x1b[17~",
	"F7":  "\x1b[18~",
	"F8":  "\x1b[19~",
	"F9":  "\x1b[20~",
	"F10": "\x1b[21~",
	"F11": "\x1b[23~",
	"F12": "\x1b[24~",

	// Page/Home/End
	"PageUp":   "\x1b[5~",
	"PageDown": "\x1b[6~",
	"Home":     "\x1b[H",
	"End":      "\x1b[F",
}
