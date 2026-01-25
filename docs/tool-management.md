# Tool Management

## Overview

Agent-Go supports per-agent tool control, allowing you to restrict which tools specific agents can access. This enables creating specialized agents with limited capabilities for security, safety, or role-specific purposes.

## Tool Policy Types

### No Policy (Default)
When no tool policy is specified, the agent has access to all available tools based on the current operation mode (Plan or Build).

### Whitelist Mode (Allowed Tools)
When `allowed_tools` is set, **only** the listed tools are available to the agent. This is the most restrictive mode and useful for creating specialized agents.

Example:
```json
{
  "name": "code-reviewer",
  "allowed_tools": [
    "use_mcp_tool",
    "create_note",
    "update_note",
    "get_todo_list"
  ]
}
```

### Blacklist Mode (Denied Tools)
When `denied_tools` is set (and `allowed_tools` is empty), the agent has access to all tools **except** the listed ones.

Example:
```json
{
  "name": "executor",
  "denied_tools": [
    "spawn_agent",
    "create_agent_definition"
  ]
}
```

### Precedence Rules
- If both `allowed_tools` and `denied_tools` are set, `allowed_tools` takes precedence (whitelist mode)
- A warning is displayed when both are configured

## Using Tool Policies

### Via Agent Studio (create_agent_definition tool)

When creating an agent through the LLM, you can specify tool policies:

```json
{
  "name": "security-auditor",
  "description": "Security-focused agent that cannot execute code",
  "system_prompt": "You are a security auditor. Review code for vulnerabilities and best practices.",
  "allowed_tools": [
    "use_mcp_tool",
    "create_note",
    "update_note",
    "create_todo",
    "update_todo",
    "get_todo_list"
  ]
}
```

### Via Manual Configuration

Create or edit agent definition files in `~/.config/agent-go/agents/`:

**Example: Code Reviewer (Whitelist)**
```json
{
  "name": "code-reviewer",
  "description": "Reviews code without execution capabilities",
  "system_prompt": "You are a code reviewer. Analyze code quality, patterns, and suggest improvements.",
  "allowed_tools": [
    "use_mcp_tool",
    "create_note",
    "update_note",
    "get_todo_list",
    "create_todo",
    "update_todo"
  ],
  "created_at": "2026-01-25T06:00:00Z",
  "updated_at": "2026-01-25T06:00:00Z"
}
```

**Example: Executor Agent (Blacklist)**
```json
{
  "name": "executor",
  "description": "Can execute commands but cannot spawn subagents",
  "system_prompt": "You execute commands directly. Work efficiently and report results.",
  "denied_tools": [
    "spawn_agent",
    "create_agent_definition"
  ],
  "created_at": "2026-01-25T06:00:00Z",
  "updated_at": "2026-01-25T06:00:00Z"
}
```

## Available Tools Reference

### Core Tools (Always Available)
- `create_todo` - Create todo items
- `update_todo` - Update todo status
- `get_todo_list` - View todo list
- `get_current_task` - Get current task
- `clear_todo` - Clear all todos
- `create_note` - Create persistent notes
- `update_note` - Update note content
- `delete_note` - Delete notes
- `name_session` - Name the current session
- `use_mcp_tool` - Access MCP server tools

### Build Mode Tools
- `execute_command` - Execute shell commands
- `kill_background_command` - Kill background processes
- `get_background_logs` - View background logs
- `list_background_commands` - List running background commands
- `create_checkpoint` - Create state checkpoint
- `list_checkpoints` - List checkpoints

### Plan Mode Tools
- `suggest_plan` - Suggest a plan for approval
- `create_agent_definition` - Create new agent definitions

### Advanced Tools
- `spawn_agent` - Spawn sub-agents (if enabled)

## Commands

### View Agent Tool Policy
```bash
/agent view <name>
```

This displays the agent's configuration including its tool policy.

### Activate Agent with Tool Policy
```bash
/agent use <name>
```

When you activate an agent with a tool policy, subsequent API requests will only include the allowed tools.

### Clear Active Agent
```bash
/agent clear
```

This restores the default tool set.

## Use Cases

### 1. Security-Focused Agent
Prevent code execution while allowing analysis:
```json
{
  "allowed_tools": [
    "use_mcp_tool",
    "create_note",
    "update_note"
  ]
}
```

### 2. Execution-Only Agent
Allow execution but prevent meta-operations:
```json
{
  "denied_tools": [
    "spawn_agent",
    "create_agent_definition"
  ]
}
```

### 3. Read-Only Reviewer
Only allow information gathering:
```json
{
  "allowed_tools": [
    "get_todo_list",
    "get_current_task",
    "use_mcp_tool"
  ]
}
```

### 4. Planning Agent
Focus on planning without execution:
```json
{
  "allowed_tools": [
    "create_todo",
    "update_todo",
    "get_todo_list",
    "create_note",
    "suggest_plan"
  ]
}
```

## Subagent Tool Policies

When spawning a subagent with a specific agent definition:
```
spawn_agent({
  "task": "Review security of auth module",
  "agent": "security-auditor"
})
```

The subagent will inherit the tool policy from the specified agent definition, independent of the parent agent's policy.

## Session Persistence

Tool policies are tracked per session:
- When you activate an agent with `/agent use`, the agent name is stored in the session
- When you restore a session, the agent definition (and its tool policy) is automatically reactivated
- The session file includes the `agent_def_name` field

## Best Practices

1. **Principle of Least Privilege**: Use whitelist mode for maximum control
2. **Clear Naming**: Name agents to reflect their capabilities (e.g., "read-only-reviewer")
3. **Document Restrictions**: Include tool limitations in the system prompt
4. **Test Policies**: Verify the agent has necessary tools before deployment
5. **Review Regularly**: Audit agent definitions periodically

## Implementation Details

- Tool filtering occurs at the API request boundary in [`sendAPIRequest()`](../src/api.go:12)
- The filtering function is [`filterToolsByPolicy()`](../src/policy.go:4)
- Agent definitions are stored in `~/.config/agent-go/agents/`
- Tool policies are optional and backward compatible

## Troubleshooting

### Agent Cannot Access Required Tool
- Check the agent definition with `/agent view <name>`
- Verify the tool name is correct (case-sensitive)
- Ensure `allowed_tools` includes the tool if using whitelist mode

### Both Policies Set
If you see a warning about both `allowed_tools` and `denied_tools` being set:
- Edit the agent definition file
- Remove one of the two fields
- The system uses `allowed_tools` (whitelist) when both are present

### Policy Not Applied
- Ensure you've activated the agent with `/agent use <name>`
- Check that the session has the correct `agent_def_name`
- Verify the agent definition file exists and is valid JSON