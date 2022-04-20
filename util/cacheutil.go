package util

import (
	"context"
	"github.com/hashicorp/golang-lru"
	"math"
	"sync"
	"time"
)

type cacheUtil struct {
	initLock  sync.Once
	cachePool *lru.TwoQueueCache
}

func (c *cacheUtil) init(ctx context.Context) error {
	var err error
	c.initLock.Do(func() {
		c.cachePool, err = lru.New2Q(500000)
	})
	return err
}

func (c *cacheUtil) Set(ctx context.Context, key string, value interface{}, timeOutSecond int64) {
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
