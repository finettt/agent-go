package main

// filterToolsByPolicy applies agent-specific tool policy and operation mode filtering to the base tool list
func filterToolsByPolicy(baseTools []Tool, agentDef *AgentDefinition, operationMode OperationMode) []Tool {
	// First, filter by operation mode
	tools := filterToolsByOperationMode(baseTools, operationMode)

	// Then apply agent-specific policy if defined
	return filterToolsByAgentPolicy(tools, agentDef)
}

// filterToolsByOperationMode filters tools based on the current operation mode (Plan vs Build)
func filterToolsByOperationMode(baseTools []Tool, operationMode OperationMode) []Tool {
	// Build sets of mode-specific tools for fast lookup
	buildModeSet := make(map[string]bool)
	for _, tool := range BuildModeTools {
		buildModeSet[tool] = true
	}

	planModeSet := make(map[string]bool)
	for _, tool := range PlanModeTools {
		planModeSet[tool] = true
	}

	filtered := make([]Tool, 0, len(baseTools))
	for _, tool := range baseTools {
		toolName := tool.Function.Name

		// In Build mode: include all tools except Plan-only tools
		if operationMode == Build {
			if !planModeSet[toolName] {
				filtered = append(filtered, tool)
			}
			continue
		}

		// In Plan mode: include all tools except Build-only tools
		if operationMode == Plan {
			if !buildModeSet[toolName] {
				filtered = append(filtered, tool)
			}
			continue
		}

		// Default: include the tool
		filtered = append(filtered, tool)
	}

	return filtered
}

// filterToolsByAgentPolicy applies agent-specific tool policy (whitelist/blacklist)
func filterToolsByAgentPolicy(baseTools []Tool, agentDef *AgentDefinition) []Tool {
	// No agent definition or no policy defined: return all tools
	if agentDef == nil || (len(agentDef.AllowedTools) == 0 && len(agentDef.DeniedTools) == 0) {
		return baseTools
	}

	// Whitelist mode (takes precedence)
	if len(agentDef.AllowedTools) > 0 {
		allowedSet := make(map[string]bool)
		for _, tool := range agentDef.AllowedTools {
			allowedSet[tool] = true
		}

		filtered := make([]Tool, 0, len(baseTools))
		for _, tool := range baseTools {
			if allowedSet[tool.Function.Name] {
				filtered = append(filtered, tool)
			}
		}
		return filtered
	}

	// Blacklist mode
	if len(agentDef.DeniedTools) > 0 {
		deniedSet := make(map[string]bool)
		for _, tool := range agentDef.DeniedTools {
			deniedSet[tool] = true
		}

		filtered := make([]Tool, 0, len(baseTools))
		for _, tool := range baseTools {
			if !deniedSet[tool.Function.Name] {
				filtered = append(filtered, tool)
			}
		}
		return filtered
	}

	return baseTools
}
