package mock

import (
	"context"
	"fmt"
	"pricesyn/tools/taskflow"
)

type ValidTask struct {
}

func (v ValidTask) Name() string {
	return "valid"
}

func (v ValidTask) Execute(ctx context.Context, flow *taskflow.Flow, result interface{}, argument map[string]string, param ...interface{}) (interface{}, error) {
	fmt.Println("execute valid")
	return nil, nil
}

var _ taskflow.Task = (*ValidTask)(nil)
