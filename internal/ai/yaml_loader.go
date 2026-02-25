package ai

import (
	"embed"

	"gopkg.in/yaml.v3"
)

//go:embed prompts/agents.yaml
var embeddedFiles embed.FS

func LoadAgentsConfig() (*AgentsFile, error) {

	data, err := embeddedFiles.ReadFile("prompts/agents.yaml")
	if err != nil {
		return nil, err
	}

	var config AgentsFile
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
