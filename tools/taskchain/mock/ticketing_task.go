package mock

import (
	"context"
	"fmt"
)

type TicketingTask struct {
}

func (v TicketingTask) Name() string {
	return "ticketing"
}

func (v TicketingTask) Execute(ctx context.Context, result interface{}, argument map[string]string, param ...interface{}) (interface{}, error) {
	fmt.Println("execute ticketing")
	return nil, nil
}
