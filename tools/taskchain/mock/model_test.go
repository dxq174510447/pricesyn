package mock

import (
	"context"
	"fmt"
	"pricesyn/tools/taskchain"
	"pricesyn/util"
	"testing"
)

const Eg1 = `
name: ticket
version: 1
stage:
  - valid
  - ticketing
  - orderConfirm
  - voucherPrint
  - ticketSuccess
failure:
  - ticketFailure
`

func FactoryInit(ctx context.Context) (*taskchain.TaskChainFactory, error) {
	factory := &taskchain.TaskChainFactory{}
	err := factory.ParseYaml(ctx, Eg1)
	if err != nil {
		return nil, err
	}
	factory.RegisterTask(ctx, &ValidTask{})
	factory.RegisterTask(ctx, &TicketingTask{})
	factory.RegisterTask(ctx, &OrderConfirmTask{})
	factory.RegisterTask(ctx, &VoucherPrintTask{})
	factory.RegisterTask(ctx, &TicketSuccessTask{})

	factory.RegisterException(ctx, &TicketFailureTask{})
	return factory, nil
}

func TestTaskChainFactory_ParseYaml(t *testing.T) {
	ctx := context.Background()
	factory, err := FactoryInit(ctx)
	if err != nil {
		t.Fatalf("%v", err)
	}
	param := make(map[string]string)
	result, err1 := factory.StartByChainId(ctx, "ticket", "service1", param)
	if err1 != nil {
		t.Fatalf("%v", err1)
	}
	if result == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(util.JsonUtil.To2PrettyString(result))
	}
}
