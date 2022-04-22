package priority_mq

import (
	"container/heap"
	"fmt"
	"testing"
)

func TestPriorityMq_Push(t *testing.T) {

	mq := new(PriorityMq)

	heap.Push(mq, &ItemNode{
		Value:    "b",
		Priority: 2,
	})
	heap.Push(mq, &ItemNode{
		Value:    "d",
		Priority: 4,
	})
	heap.Push(mq, &ItemNode{
		Value:    "e",
		Priority: 5,
	})
	heap.Push(mq, &ItemNode{
		Value:    "a",
		Priority: 1,
	})
	heap.Push(mq, &ItemNode{
		Value:    "c",
		Priority: 3,
	})
	fmt.Printf("%v\n", mq)
	var result *ItemNode = (heap.Pop(mq)).(*ItemNode)
	for result != nil {
		fmt.Printf("%v -> %v\n", result.Value, mq)
		if mq.Len() == 0 {
			break
		}
		result = (heap.Pop(mq)).(*ItemNode)
	}
}
