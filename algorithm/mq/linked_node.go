package mq

type LinkedNode struct {
	Item interface{}
	Next *LinkedNode
}

func (n *LinkedNode) SetNext(next *LinkedNode) {
	n.Next = next
}

func NewNode(item interface{}) *LinkedNode {
	return &LinkedNode{
		Item: item,
	}
}