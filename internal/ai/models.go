package ai

type AgentConfig struct {
	Name         string  `yaml:"name"`
	SystemPrompt string  `yaml:"system_prompt"`
	Temperature  float64 `yaml:"temperature"`
}

type AgentsFile struct {
	Agents []AgentConfig `yaml:"agents"`
}
