package mock

import (
	"context"
	"fmt"
	"pricesyn/tools/taskflow"
)

type TicketSuccessTask struct {
}

func (v TicketSuccessTask) Name() string {
	return "ticketSuccess"
}

func (v TicketSuccessTask) Execute(ctx context.Context, flow *taskflow.Flow, result interface{}, argument map[string]string, param ...interface{}) (interface{}, error) {
	fmt.Println("execute ticketSuccess")
	return nil, nil
}

var _ taskflow.Task = (*TicketSuccessTask)(nil)
