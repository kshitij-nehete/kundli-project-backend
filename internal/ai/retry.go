package ai

import "context"

func WithRetry(
	ctx context.Context,
	attempts int,
	fn func() (map[string]interface{}, error),
) (map[string]interface{}, error) {

	var err error
	var result map[string]interface{}

	for i := 0; i < attempts; i++ {
		result, err = fn()
		if err == nil {
			return result, nil
		}
	}

	return nil, err
}
