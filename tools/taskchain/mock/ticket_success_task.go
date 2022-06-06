package mock

import (
	"context"
	"fmt"
)

type TicketSuccessTask struct {
}

func (v TicketSuccessTask) Name() string {
	return "ticketSuccess"
}

func (v TicketSuccessTask) Execute(ctx context.Context, result interface{}, argument map[string]string, param ...interface{}) (interface{}, error) {
	fmt.Println("execute ticketSuccess")
	return nil, nil
}
