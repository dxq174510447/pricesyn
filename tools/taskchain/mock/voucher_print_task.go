package mock

import (
	"context"
	"fmt"
)

type VoucherPrintTask struct {
}

func (v VoucherPrintTask) Name() string {
	return "voucherPrint"
}

func (v VoucherPrintTask) Execute(ctx context.Context, result interface{}, argument map[string]string, param ...interface{}) (interface{}, error) {
	fmt.Println("execute voucherPrint")
	return nil, nil
}
