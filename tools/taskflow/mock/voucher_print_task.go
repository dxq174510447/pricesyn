package mock

import (
	"context"
	"fmt"
	"pricesyn/tools/taskflow"
)

type VoucherPrintTask struct {
}

func (v VoucherPrintTask) Name() string {
	return "voucherPrint"
}

func (v VoucherPrintTask) Execute(ctx context.Context, flow *taskflow.Flow, result interface{}, argument map[string]string, param ...interface{}) (interface{}, error) {
	fmt.Println("execute voucherPrint")
	return nil, nil
}

var _ taskflow.Task = (*VoucherPrintTask)(nil)
