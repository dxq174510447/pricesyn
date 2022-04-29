package util

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestMachineUtil_New(t *testing.T) {
	MachineUtil.NewMachine()
}

func TestMathUtil_Twrailmon(t *testing.T) {
	m := make(map[string]int)
	m["a"] = 1
	fmt.Println(m["a"])
	//MathUtil.Twrailmon(970, 25)
}

func HeiheTest(ctx context.Context, content string) (bool, error) {
	fmt.Println("heihe")
	return true, nil
}

func BeiheTest(ctx context.Context, content string) (bool, error) {
	fmt.Println("baihe")
	return true, nil
}

func AllTest(ctx context.Context, content string, fn func(ctx context.Context, content string) (bool, error)) (bool, error) {
	return fn(ctx, content)
}

func TestPointer_New(t *testing.T) {
	AllTest(context.Background(), "haha", HeiheTest)
	AllTest(context.Background(), "haha", BeiheTest)

	var aa []string = []string{"a", "b"}
	mk := make(map[int]int)
	for i := 0; i < 10000; i++ {
		mk[i] = i
	}
	for i := range aa {
		j := i
		go func() {
			//time.Sleep(time.Second * 5)
			//fmt.Println(j)
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("%v", err)
				}
			}()

			for k, v := range mk {
				fmt.Println(j, k, v)
			}
		}()

	}
	time.Sleep(time.Minute * 10)
}
