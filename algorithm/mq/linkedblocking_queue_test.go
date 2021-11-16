package mq

import (
	"fmt"
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
			}else{
				fmt.Println(m)
			}
		}
	}()

	time.Sleep(time.Minute * 5)
}
