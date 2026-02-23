package ai

import "context"

type Agent interface {
	Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)
}
