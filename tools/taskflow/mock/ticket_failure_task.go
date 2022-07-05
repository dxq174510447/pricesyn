package mock

import (
	"context"
	"fmt"
	"pricesyn/tools/taskflow"
)

type TicketFailureTask struct {
}

func (v TicketFailureTask) Name() string {
	return "ticketFailure"
}

func (v TicketFailureTask) Callback(ctx context.Context, stageName string, err error, param ...interface{}) error {
	fmt.Println("execute ticketFailure")
	return nil
}

var _ taskflow.ExceptionTask = (*TicketFailureTask)(nil)
