package main

// ANSI color codes for terminal output - these are now variables that can be disabled
var (
	ColorRed           = "\033[31m"
	ColorGreen         = "\033[32m"
	ColorBlue          = "\033[34m"
	ColorYellow        = "\033[33m"
	ColorPurple        = "\033[35m"
	ColorCyan          = "\033[36m"
	ColorWhite         = "\033[37m"
	ColorBlack         = "\033[30m"
	ColorHighlight     = "\033[38;2;255;147;251m" // #FF93FB - Highlighting
	ColorMain          = "\033[38;2;255;255;255m" // #FFF    - Main color
	ColorMeta          = "\033[38;2;170;170;170m" // #AAA    - Unnecessary things
	StyleBold          = "\033[1m"
	StyleDim           = "\033[2m"
	StyleItalic        = "\033[3m"
	StyleUnderline     = "\033[4m"
	StyleBlink         = "\033[5m"
	StyleReverse       = "\033[7m"
	StyleHidden        = "\033[8m"
	StyleStrikethrough = "\033[9m"
	ColorReset         = "\033[0m"
)

// Context compression settings
const (
	CompressionTemp                      = 0.15
	CompressionMaxTokens                 = 1500
	CompressionExcludeSystemPrompts      = true // Always exclude system prompts from compression
	CompressionPreserveAgentInstructions = true // Preserve agent-specific instructions
)

// Sub-agent limits
const (
	MaxSubAgentIterations = 50
)

// Tool loop detection
const (
	// MaxRepeatedToolCalls is the maximum number of times the same tool call can be repeated
	// before the agent is stopped and asked to try a different approach
	MaxRepeatedToolCalls = 3

	// ToolLoopStopMessage is the message sent to the model when it gets stuck in a tool loop
	ToolLoopStopMessage = "STOP: You appear to be stuck in a loop, repeatedly calling the same tool. Please step back, analyze what you've tried so far, and try a completely different approach to solve this task. If you cannot complete the task, explain what's blocking you."
)

// Default configuration values
const (
	DefaultTemp                  = 0.1
	DefaultMaxTokens             = 32768
	DefaultAPIURL                = "https://api.openai.com"
	DefaultModel                 = "gpt-3.5-turbo"
	DefaultMiniModel             = "gpt-4o-mini"
	DefaultRAGSnippets           = 5
	DefaultAutoCompressThreshold = 20
	DefaultModelContextLength    = 262144
)

// Valid todo statuses
var ValidTodoStatuses = map[string]bool{
	"pending":     true,
	"in-progress": true,
	"completed":   true,
}
