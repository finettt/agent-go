package main
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// loadSkills loads skills from system, .agent-go project, and project directories
func loadSkills() ([]Skill, error) {
	var skills []Skill
	seenSkills := make(map[string]bool)

	// Helper to add unique skills
	addSkills := func(newSkills []Skill) {
		for _, s := range newSkills {
			if !seenSkills[s.Name] {
				skills = append(skills, s)
				seenSkills[s.Name] = true
			}
		}
	}

	// System-wide skills
	home, err := os.UserHomeDir()
	if err == nil {
		systemSkillsPath := filepath.Join(home, ".config", "agent-go", "skills")
		sysSkills, _ := loadSkillsFromDir(systemSkillsPath)
		addSkills(sysSkills)
	}

	// Project-wide skills in .agent-go/skills
	cwd, err := os.Getwd()
	if err == nil {
		agentGoSkillsPath := filepath.Join(cwd, ".agent-go", "skills")
		agentGoSkills, _ := loadSkillsFromDir(agentGoSkillsPath)
		addSkills(agentGoSkills)
	}

	// Project-wide skills in skills/ (legacy support)
	if err == nil {
		projectSkillsPath := filepath.Join(cwd, "skills")
		projSkills, _ := loadSkillsFromDir(projectSkillsPath)
		addSkills(projSkills)
	}

	return skills, nil
}

func loadSkillsFromDir(dir string) ([]Skill, error) {
	var skills []Skill
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Check for skill.json inside directory
			skillPath := filepath.Join(dir, entry.Name(), "skill.json")
			skill, err := loadSkillFromFile(skillPath)
			if err == nil {
				skills = append(skills, *skill)
			}
		} else if filepath.Ext(entry.Name()) == ".json" {
			skillPath := filepath.Join(dir, entry.Name())
			skill, err := loadSkillFromFile(skillPath)
			if err == nil {
				skills = append(skills, *skill)
			}
		} else if filepath.Ext(entry.Name()) == ".sh" {
			skillPath := filepath.Join(dir, entry.Name())
			skill, err := loadSkillFromScript(skillPath)
			if err == nil {
				skills = append(skills, *skill)
			}
		}
	}
	return skills, nil
}

func loadSkillFromScript(path string) (*Skill, error) {
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	
	// Default parameters for script skills
	params := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"args": map[string]interface{}{
				"type":        "string",
				"description": "Arguments to pass to the script",
			},
		},
	}

	return &Skill{
		Name:        name,
		Description: fmt.Sprintf("Executes the %s script", name),
		Command:     path,
		Parameters:  params,
	}, nil
}

func loadSkillFromFile(path string) (*Skill, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var skill Skill
	if err := json.Unmarshal(data, &skill); err != nil {
		return nil, err
	}

	// Basic validation
	if skill.Name == "" || skill.Command == "" {
		return nil, fmt.Errorf("skill name and command are required")
	}

	// Resolve command path if it's relative and exists locally
	// We assume that if the command is a file in the same directory as skill.json, we should use the absolute path
	dir := filepath.Dir(path)
	cmdPath := filepath.Join(dir, skill.Command)

	// Check if the command exists as a file relative to the skill definition
	if info, err := os.Stat(cmdPath); err == nil && !info.IsDir() {
		skill.Command = cmdPath
	}

	return &skill, nil
}