package main

// ANSI color codes for terminal output
const (
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
	CompressionTemp      = 0.15
	CompressionMaxTokens = 1500
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
