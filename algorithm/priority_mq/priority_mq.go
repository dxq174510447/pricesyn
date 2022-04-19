package priority_mq

import (
	"container/heap"
)

// An Item is something we manage in a priority queue.
type ItemNode struct {
	Value    interface{} // The value of the item; arbitrary.
	Priority int64       // The priority of the item in the queue.
	index    int         // The index of the item in the heap.
}

// A PriorityMq implements heap.Interface and holds Items.
type PriorityMq []*ItemNode

func (pq PriorityMq) Len() int { return len(pq) }

func (pq PriorityMq) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityMq) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityMq) Push(x interface{}) {
	n := len(*pq)
	item := x.(*ItemNode)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityMq) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityMq) update(item *ItemNode, value string, priority int64) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.index)
}

func PushNode(mq heap.Interface, value interface{}, priority int64) {

}
