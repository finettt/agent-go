package main

// filterToolsByPolicy applies agent-specific tool policy to filter the base tool list
func filterToolsByPolicy(baseTools []Tool, agentDef *AgentDefinition) []Tool {
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
