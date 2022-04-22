package cache

import (
	"context"
	"github.com/hashicorp/golang-lru"
	"sync"
	"time"
)

type Cache struct {
	initLock   sync.Once
	cachePool  *lru.TwoQueueCache
	timerClear *CacheTimeOutClear
}

func (c *Cache) init(ctx context.Context) error {
	var err error
	c.initLock.Do(func() {
		c.cachePool, err = lru.New2Q(500000)
		c.timerClear = &CacheTimeOutClear{
			CachePool: c.cachePool,
		}
	})
	return err
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, timeOutSecond int64) {
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

	if expireTime > 0 {
		c.timerClear.Push(ctx, key, expireTime)
	}
}

func (c *Cache) Delete(ctx context.Context, key string) {
	c.init(ctx)
	c.cachePool.Remove(key)
	c.timerClear.Remove(ctx, key)
}

func (c *Cache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	c.init(ctx)
	node, exist := c.cachePool.Get(key)
	if !exist {
		return nil, false, nil
	}
	cn := node.(*cacheNode)
	if cn.expireTime <= 0 {
		return cn.value, true, nil
	}
	if cn.expireTime < time.Now().Unix() {
		c.cachePool.Remove(key)
		return nil, false, nil
	}
	return cn.value, true, nil
}

type cacheNode struct {
	expireTime int64
	value      interface{}
}
