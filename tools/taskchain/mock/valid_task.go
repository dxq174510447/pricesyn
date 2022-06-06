package mock

import (
	"context"
	"fmt"
)

type ValidTask struct {
}

func (v ValidTask) Name() string {
	return "valid"
}

func (v ValidTask) Execute(ctx context.Context, result interface{}, argument map[string]string, param ...interface{}) (interface{}, error) {
	fmt.Println("execute valid")
	return nil, nil
}
