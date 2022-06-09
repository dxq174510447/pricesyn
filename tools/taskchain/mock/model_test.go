package mock

import (
	"context"
	"flag"
	"fmt"
	"pricesyn/tools/taskchain"
	"pricesyn/util"
	"testing"
)

const Eg1 = `
name: ticket
version: 1
stage:
  - name: valid
  - name: ticketing
  - name: orderConfirm
  - name: voucherPrint
  - name: ticketSuccess
failure:
  - name: ticketFailure
`

const Eg2 = `
name: ticket
version: 2
stage:
  - name: valid
  - name: ticketing
  - name: orderConfirm
  - name: _stop
  - name: voucherPrint
  - name: ticketSuccess
failure:
  - name: ticketFailure
`

var pwd = flag.String("pwd", "", "Input Your pwd")

func FactoryInit(ctx context.Context) (*taskchain.TaskChainFactory, error) {
	factory := &taskchain.TaskChainFactory{}
	err := factory.ParseYaml(ctx, Eg1)
	if err != nil {
		return nil, err
	}
	err = factory.ParseYaml(ctx, Eg2)
	if err != nil {
		return nil, err
	}
	factory.RegisterTask(ctx, &ValidTask{})
	factory.RegisterTask(ctx, &TicketingTask{})
	factory.RegisterTask(ctx, &OrderConfirmTask{})
	factory.RegisterTask(ctx, &VoucherPrintTask{})
	factory.RegisterTask(ctx, &TicketSuccessTask{})

	factory.RegisterException(ctx, &TicketFailureTask{})

	factory.RegisterPersistenceService(ctx, &DbTaskChainService{
		pwd: *pwd,
	})
	return factory, nil
}

func TestTaskChainFactory_ParseYaml(t *testing.T) {
	flag.Parse()

	ctx := context.Background()
	factory, err := FactoryInit(ctx)
	if err != nil {
		t.Fatalf("%v", err)
	}
	param := make(map[string]string)
	result, err1 := factory.StartByChainId(ctx, "ticket", "service14", param)
	if err1 != nil {
		t.Fatalf("%v", err1)
	}
	if result == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(util.JsonUtil.To2PrettyString(result))
	}
}

func TestTaskChainFactory_Eg2(t *testing.T) {
	flag.Parse()

	ctx := context.Background()
	factory, err := FactoryInit(ctx)
	if err != nil {
		t.Fatalf("%v", err)
	}
	param := make(map[string]string)
	result, err1 := factory.StartByChainId(ctx, "ticket", "service14", param)
	if err1 != nil {
		t.Fatalf("%v", err1)
	}
	if result == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(util.JsonUtil.To2PrettyString(result))
	}
}
