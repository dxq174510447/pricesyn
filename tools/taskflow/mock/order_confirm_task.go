package mock

import (
	"context"
	"fmt"
	"pricesyn/tools/taskflow"
)

type OrderConfirmTask struct {
}

func (v OrderConfirmTask) Name() string {
	return "orderConfirm"
}

func (v OrderConfirmTask) Execute(ctx context.Context, flow *taskflow.Flow, result interface{}, argument map[string]string, param ...interface{}) (interface{}, error) {
	fmt.Println("execute orderConfirm")
	return nil, nil
}

var _ taskflow.Task = (*OrderConfirmTask)(nil)
