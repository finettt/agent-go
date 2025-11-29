package main

// ANSI color codes for terminal output
const (
	ColorRed       = "\033[31m"
	ColorGreen     = "\033[32m"
	ColorBlue      = "\033[34m"
	ColorYellow    = "\033[33m"
	ColorPurple    = "\033[35m"
	ColorReset     = "\033[0m"
	ColorHighlight = "\033[38;2;255;147;251m" // #FF93FB - Highlighting
	ColorMain      = "\033[38;2;255;255;255m" // #FFF    - Main color
	ColorMeta      = "\033[38;2;170;170;170m" // #AAA    - Unnecessary things
)

// Streaming constants
const (
	StreamDataPrefix = "data: "
	StreamDoneMarker = "[DONE]"
)

// Context compression settings
const (
	CompressionTemp      = 0.3
	CompressionMaxTokens = 500
)

// Sub-agent limits
const (
	MaxSubAgentIterations = 50
)

// Default configuration values
const (
	DefaultTemp                  = 0.1
	DefaultMaxTokens             = 1000
	DefaultAPIURL                = "https://api.openai.com"
	DefaultModel                 = "gpt-3.5-turbo"
	DefaultRAGSnippets           = 5
	DefaultAutoCompressThreshold = 20
	DefaultModelContextLength    = 131072
)

// Valid todo statuses
var ValidTodoStatuses = map[string]bool{
	"pending":     true,
	"in-progress": true,
	"completed":   true,
}
