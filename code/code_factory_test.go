package code

import (
	"context"
	"fmt"
	"pricesyn/util"
	"testing"
)

func TestCodeFactory_GetColumns(t *testing.T) {
	factory := &CodeFactory{}

	ctx := context.Background()
	col, err := factory.GetColumns(ctx, "flight_contact", "flight_order_db")
	if err != nil {
		panic(err)
	}
	fmt.Println(util.JsonUtil.To2PrettyString(col))

}

func TestCodeFactory_Generate(t *testing.T) {
	factory := &CodeFactory{}

	title := util.StringUtil.FieldName("flight_order_db")
	fmt.Println(title)
	ctx := context.Background()
	err := factory.Generate(ctx, "flight_contact", "flight_order_db")
	if err != nil {
		panic(err)
	}

}
