package mq

import (
	"fmt"
	"pricesyn/util"
	"reflect"
	"testing"
	"time"
)

func TestNewLinkedBlockingQueue(t *testing.T) {

	queue := NewLinkedBlockingQueue()
	var i int = 1
	go func() {
		for {
			time.Sleep(time.Second * 1)
			queue.Offer(i)
			i++
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * 1)
			queue.Offer(i)
			i++
		}
	}()

	go func() {
		for {
			m := queue.Poll(10)
			if m == nil || reflect.ValueOf(m).IsZero() {
				continue
			} else {
				fmt.Println(m)
			}
		}
	}()

	time.Sleep(time.Minute * 5)
}

func TestMm(t *testing.T) {
	var source []string = []string{"a", "b"}

	var target []string = make([]string, len(source)+1)
	copy(target, source)
	target[len(target)-1] = "c"
	fmt.Println(util.JsonUtil.To2String(target))
	source[0] = "d"
	fmt.Println(util.JsonUtil.To2String(target))

	source = make([]string, 2, 3)
	source[0] = "a"
	source[1] = "b"

	source1 := append(source, "c")
	fmt.Println(util.JsonUtil.To2String(source1))
	source[0] = "d"
	fmt.Println(util.JsonUtil.To2String(source1))
}
