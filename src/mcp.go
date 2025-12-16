package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// mcpManager manages MCP client sessions
type mcpManager struct {
	mu              sync.Mutex
	clients         map[string]*mcp.Client
	sessions        map[string]*mcp.ClientSession
	implementations map[string]*mcp.Implementation
}

var globalMCP = newMCPManager()

func newMCPManager() *mcpManager {
	return &mcpManager{
		clients:         make(map[string]*mcp.Client),
		sessions:        make(map[string]*mcp.ClientSession),
		implementations: make(map[string]*mcp.Implementation),
	}
}

// ensureMCP ensures a session to a configured MCP server
func (m *mcpManager) ensureMCP(serverName string) (*mcp.ClientSession, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if sess, ok := m.sessions[serverName]; ok && sess != nil {
		return sess, nil
	}

	mcpServer, ok := config.MCPs[serverName]
	if !ok {
		return nil, fmt.Errorf("mcp server not found in config: %s", serverName)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "agent-go", Version: "v0.1.0"}, nil)
	// Launch the MCP server using the configured command
	cmdParts := strings.Fields(mcpServer.Command)
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	transport := &mcp.CommandTransport{Command: cmd}
	ctx := context.Background()

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mcp server '%s': %w", serverName, err)
	}

	m.clients[serverName] = client
	m.sessions[serverName] = session

	return session, nil
}

// useMCPTool calls a tool on a specified MCP server
func useMCPTool(serverName, toolName string, arguments map[string]interface{}) (string, error) {
	session, err := globalMCP.ensureMCP(serverName)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	params := &mcp.CallToolParams{
		Name:      toolName,
		Arguments: arguments,
	}

	res, err := session.CallTool(ctx, params)
	if err != nil {
		return "", fmt.Errorf("mcp tool call failed: %w", err)
	}
	if res.IsError {
		return "", fmt.Errorf("mcp tool returned an error")
	}

	var out string
	for _, c := range res.Content {
		if t, ok := c.(*mcp.TextContent); ok {
			out += t.Text
		}
	}
	if out == "" {
		out = "(no text output from MCP tool)"
	}

	return out, nil
}

// getMCPToolInfo connects to all configured MCP servers and returns a summary of their tools.
func getMCPToolInfo() string {
	if config.MCPs == nil || len(config.MCPs) == 0 {
		return ""
	}

	var info strings.Builder
	info.WriteString("\n\nThe following MCP servers are available:\n")

	for name := range config.MCPs {
		session, err := globalMCP.ensureMCP(name)
		if err != nil {
			info.WriteString(fmt.Sprintf("- Server Name: '%s' (connection failed: %v)\n", name, err))
			continue
		}

		info.WriteString(fmt.Sprintf("- Server Name: '%s'\n", name))
		tools, err := session.ListTools(context.Background(), nil)
		if err != nil {
			info.WriteString(fmt.Sprintf("  (Failed to list tools: %v)\n", err))
			continue
		}

		if len(tools.Tools) > 0 {
			info.WriteString("  Tools:\n")
			for _, tool := range tools.Tools {
				info.WriteString(fmt.Sprintf("    - %s: %s\n", tool.Name, tool.Description))
			}
		} else {
			info.WriteString("  (No tools reported)\n")
		}
	}
	info.WriteString("You can use tools from these servers with the `use_mcp_tool` function, specifying the `server_name` from the list above.\n")
	return info.String()
}
