package util

import (
	"context"
	"fmt"
	"testing"
)

func TestHttpUtil_PostBody(t *testing.T) {
	resp, err := HttpUtil.Get(context.Background(), nil, "http://10.2.10.12:80", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp))
}
