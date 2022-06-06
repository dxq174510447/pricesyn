package mock

import (
	"context"
	"fmt"
)

type OrderConfirmTask struct {
}

func (v OrderConfirmTask) Name() string {
	return "orderConfirm"
}

func (v OrderConfirmTask) Execute(ctx context.Context, result interface{}, argument map[string]string, param ...interface{}) (interface{}, error) {
	fmt.Println("execute orderConfirm")
	return nil, nil
}
