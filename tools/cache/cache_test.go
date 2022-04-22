package cache

import (
	"context"
	"fmt"
	"pricesyn/util"
	"strings"
	"testing"
	"time"
)

func TestCache_Set(t *testing.T) {
	ctx := context.Background()
	c := &Cache{}

	var keys []string = []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i",
	}
	for i, key := range keys {
		index := i
		k := key
		go util.FuncUtil.HandlePanic(ctx, func() {
			c.Set(ctx, k, "1", int64(index+1))
		})()
	}

	go util.FuncUtil.HandlePanic(ctx, func() {
		for {
			var hasKeys []string
			for _, key := range keys {
				if result, ok, _ := c.Get(ctx, key); ok {
					hasKeys = append(hasKeys, fmt.Sprintf("%s|%v", key, result))
				}
			}
			fmt.Println(strings.Join(hasKeys, ","))
			time.Sleep(time.Second * 1)
		}
	})()

	time.Sleep(time.Second * 100)
}
