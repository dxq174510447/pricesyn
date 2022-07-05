package mock

import (
	"context"
	"flag"
	"fmt"
	"pricesyn/tools/taskflow"
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
`

const Eg2 = `
name: ticket
version: 2
stage:
  - name: valid
  - name: ticketing
  - name: orderConfirm
    failure:
    - ticketFailure
    - ticketFailure
  - name: _stop
  - name: voucherPrint
  - name: ticketSuccess
`

var pwd = flag.String("pwd", "", "Input Your pwd")
var ChainName = "ticket"

func FactoryInit(ctx context.Context) (*taskflow.TaskflowFactory, error) {
	factory := &taskflow.TaskflowFactory{}
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

var serviceId string = "aaaa1"

func TestTaskChainFactory_Begin(t *testing.T) {
	flag.Parse()

	ctx := context.Background()
	factory, err := FactoryInit(ctx)
	if err != nil {
		t.Fatalf("%v", err)
	}

	param := make(map[string]string)
	param["serviceId"] = serviceId
	result, err1 := factory.Begin(ctx, ChainName, serviceId, param)
	//result, err1 := factory.BeginWithVersion(ctx,ChainName,1,serviceId,param)

	if err1 != nil {
		t.Fatalf("%v", err1)
	}
	if result == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(util.JsonUtil.To2PrettyString(result))
	}
}

func TestTaskChainFactory_Start(t *testing.T) {
	flag.Parse()

	ctx := context.Background()
	factory, err := FactoryInit(ctx)
	if err != nil {
		t.Fatalf("%v", err)
	}

	param := make(map[string]string)
	param["serviceId"] = serviceId

	//result, err1 := factory.StartTask(ctx, ChainName, serviceId, param)
	//result, err1 := factory.ReStartTaskWithStageName(ctx, ChainName, serviceId, "ticketing", param)
	result, err1 := factory.StartTaskWithStageName(ctx, ChainName, serviceId, "voucherPrint", param)
	if err1 != nil {
		t.Fatalf("%v", err1)
	}
	if result == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(util.JsonUtil.To2PrettyString(result))
	}
}
