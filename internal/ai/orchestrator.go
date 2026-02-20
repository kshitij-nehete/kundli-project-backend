package ai

import "context"

type Agent interface {
	Execute(ctx context.Context, input interface{}) (interface{}, error)
}

type Orchestrator struct {
	agents []Agent
}

func NewOrchestrator(agents []Agent) *Orchestrator {
	return &Orchestrator{
		agents: agents,
	}
}

func (o *Orchestrator) Run(ctx context.Context, input interface{}) (interface{}, error) {
	var result interface{} = input
	var err error

	for _, agent := range o.agents {
		result, err = agent.Execute(ctx, result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
