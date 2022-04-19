package util

import (
	"container/heap"
	"context"
	"github.com/hashicorp/golang-lru"
	"math"
	"pricesyn/algorithm/priority_mq"
	"pricesyn/tools/locker"
	"sync"
	"time"
)

type cacheUtil struct {
	initLock  sync.Once
	cachePool *lru.TwoQueueCache
	clearLock locker.TryMutex
	keyHeap   heap.Interface
}

func (c *cacheUtil) init(ctx context.Context) error {
	var err error
	c.initLock.Do(func() {
		c.cachePool, err = lru.New2Q(500000)
		c.keyHeap = new(priority_mq.PriorityQueue)
	})
	return err
}

func (c *cacheUtil) clearTimeOut1(ctx context.Context) {

}

func (c *cacheUtil) clearTimeOut(ctx context.Context) error {
	lock := c.clearLock.TryLock()
	if !lock {
		return nil
	}
	go FuncUtil.HandlePanic(context.Background(), func() {
		defer c.clearLock.Unlock()
		c.clearTimeOut1(context.Background())
	})
	return nil
}

func (c *cacheUtil) Cache(ctx context.Context, key string, value interface{}, timeOutSecond int64) {
	c.init(ctx)
	var expireTime int64 = 0
	if timeOutSecond > 0 {
		expireTime = time.Now().Unix() + timeOutSecond
	}
	node := &cacheNode{
		expireTime: expireTime,
		value:      value,
	}
	c.cachePool.Add(key, node)

	priority := expireTime
	if priority <= 0 {
		priority = math.MaxInt64
	}
	item := &priority_mq.ItemNode{
		Value:    key,
		Priority: priority,
	}
	heap.Push(c.keyHeap, item)
	heap.Pop()
}

func (c *cacheUtil) Get(ctx context.Context, key string) (bool, interface{}, error) {
	c.init(ctx)
	node, exist := c.cachePool.Get(key)
	if !exist {
		return false, nil, nil
	}
	cn := node.(*cacheNode)
	if cn.expireTime <= 0 {
		return true, cn.value, nil
	}
	if cn.expireTime > time.Now().Unix() {
		c.cachePool.Remove(key)
		return false, nil, nil
	}
	return true, cn.value, nil
}

type cacheNode struct {
	expireTime int64
	value      interface{}
}

var CacheUtil cacheUtil = cacheUtil{}
