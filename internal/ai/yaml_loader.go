package ai

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadAgentsConfig(path string) (*AgentsFile, error) {

	data, err := os.ReadFile(path)
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
