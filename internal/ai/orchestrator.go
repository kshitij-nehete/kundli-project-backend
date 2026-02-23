package ai

import "context"

type Orchestrator struct {
	agents []Agent
}

func NewOrchestrator(agents []Agent) *Orchestrator {
	return &Orchestrator{agents: agents}
}

func (o *Orchestrator) Run(
	ctx context.Context,
	input map[string]interface{},
) (map[string]interface{}, error) {

	result := input
	var err error

	for _, agent := range o.agents {
		result, err = agent.Execute(ctx, result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
