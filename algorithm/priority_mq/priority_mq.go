package priority_mq

import (
	"container/heap"
	"fmt"
	"strings"
)

// An Item is something we manage in a priority queue.
type ItemNode struct {
	Value    interface{} // The value of the item; arbitrary.
	Priority int64       // The priority of the item in the queue.
	Index    int         // The index of the item in the heap.
}

// A PriorityMq implements heap.Interface and holds Items.
type PriorityMq []*ItemNode

func (pq PriorityMq) Len() int { return len(pq) }

func (pq PriorityMq) String() string {
	if len(pq) == 0 {
		return ""
	}
	var result []string
	for _, i := range pq {
		result = append(result, fmt.Sprintf("%v|%d", i.Value, i.Priority))
	}
	return strings.Join(result, ",")
}

func (pq PriorityMq) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityMq) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityMq) Push(x interface{}) {
	n := len(*pq)
	item := x.(*ItemNode)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityMq) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq PriorityMq) GetTopPriority() int64 {
	n := len(pq)
	if n <= 0 {
		return -1
	}
	return pq[0].Priority
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityMq) update(item *ItemNode, value string, priority int64) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.Index)
}
