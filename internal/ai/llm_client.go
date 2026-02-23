package ai

import "context"

type LLMClient interface {
	Generate(ctx context.Context, systemPrompt string, input map[string]interface{}, temperature float64) (map[string]interface{}, error)
}

type StubLLMClient struct{}

func (s *StubLLMClient) Generate(
	ctx context.Context,
	systemPrompt string,
	input map[string]interface{},
	temperature float64,
) (map[string]interface{}, error) {

	// Simulated response
	input["planetary_data"] = map[string]interface{}{
		"sun":  "Libra",
		"moon": "Aries",
	}

	input["analysis"] = map[string]interface{}{
		"summary": "You are analytical and balanced.",
	}

	input["confidence_score"] = 85

	return input, nil
}
