package ai

import "context"

type ConfigurableAgent struct {
	Config    AgentConfig
	LLMClient LLMClient
}

func (a *ConfigurableAgent) Execute(
	ctx context.Context,
	input map[string]interface{},
) (map[string]interface{}, error) {

	return a.LLMClient.Generate(
		ctx,
		a.Config.SystemPrompt,
		input,
		a.Config.Temperature,
	)
}
