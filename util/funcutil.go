package util

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
)

type funcUtil struct {
}

// HandlePanic 封装对 panic 的处理, 在开协程时比较有用
func (fun *funcUtil) HandlePanic(ctx context.Context, f func()) func() {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				stack := strings.Join(strings.Split(string(debug.Stack()), "\n")[2:], "\n")
				fmt.Printf("%s", stack)
			}
		}()
		f()
	}
}

// HandlePanicV2 增加参数传递 避免开协程时的闭包问题
func (fun *funcUtil) HandlePanicV2(ctx context.Context, f func(interface{})) func(interface{}) {
	return func(arg interface{}) {
		defer func() {
			if err := recover(); err != nil {
				stack := strings.Join(strings.Split(string(debug.Stack()), "\n")[2:], "\n")
				fmt.Printf("%s", stack)
			}
		}()
		f(arg)
	}
}

var FuncUtil funcUtil = funcUtil{}
