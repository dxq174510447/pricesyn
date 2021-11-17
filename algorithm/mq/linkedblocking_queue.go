package mq

import (
	"fmt"
	"pricesyn/util"
	"sync"
	"sync/atomic"
	"time"
)

type IdNoder interface {
	IdKey() string
}

type LinkedBlockingQueue struct {
	head     *LinkedNode
	last     *LinkedNode
	takeLock *sync.Mutex
	putLock  *sync.Mutex
	count    *int64
	syn      chan int
}

//Poll Retrieves and removes the head of this queue, or returns null if this queue is empty
func (l *LinkedBlockingQueue) Poll(sec int) interface{} {
	if *l.count == 0 && sec > 0 {
		// 要么有数据插入到队列 要么等待N second 之后 返回nil
		select {
		case <-l.syn:
			// 直接放行
		case <-time.After(time.Duration(sec) * time.Second):
			return nil
		}
	}
	l.takeLock.Lock()
	defer l.takeLock.Unlock()
	if *l.count == 0 {
		return nil
	}

	h := l.head
	first := h.Next

	// 队列已经空了
	if h.Item == nil && first == nil {
		return nil
	}

	l.head = first

	result := first.Item
	first.Item = nil

	h = nil

	c := atomic.AddInt64(l.count, -1)
	l.count = &c

	return result
}

func (l *LinkedBlockingQueue) GetLen() int64 {
	return *l.count
}

//Offer Inserts the specified element at the tail of this queue
func (l *LinkedBlockingQueue) Offer(ele interface{}) {
	l.putLock.Lock()
	defer l.putLock.Unlock()

	node := NewNode(ele)
	l.last.Next = node
	l.last = node

	c := atomic.AddInt64(l.count, 1)
	l.count = &c

	if c == 1 {
		select {
		case l.syn <- 1:
			// 通知等待的可以执行了
		case <-time.After(time.Duration(1) * time.Millisecond):
			// 保证l.syn不被锁住
		}
	}
}

func (l *LinkedBlockingQueue) PrintLink() {
	var link []interface{}
	var current *LinkedNode = l.head
	var index int = 0
	for {
		if current == nil {
			break
		}
		if index == 0 {
			link = append(link, "头部")
			current = current.Next
			index++
			continue
		}
		idNode := current.Item.(IdNoder)
		link = append(link, idNode.IdKey())
		current = current.Next
		index++
	}
	fmt.Println(util.JsonUtil.To2String(link))
}

func NewLinkedBlockingQueue() *LinkedBlockingQueue {
	node := NewNode(nil)
	m := int64(0)
	return &LinkedBlockingQueue{
		head:     node,
		last:     node,
		takeLock: &sync.Mutex{},
		putLock:  &sync.Mutex{},
		count:    &m,
		syn:      make(chan int),
	}
}
