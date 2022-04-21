package cache

import (
	"container/heap"
	"context"
	lru "github.com/hashicorp/golang-lru"
	"pricesyn/algorithm/priority_mq"
	"pricesyn/util"
	"sync"
	"time"
)

/*
CacheTimeOutClear
定时清理 权重就是时间搓
*/
type CacheTimeOutClear struct {
	mq        *priority_mq.PriorityMq
	keyMap    map[string]*priority_mq.ItemNode
	CachePool *lru.TwoQueueCache
	initLock  sync.Once
	lock      sync.Mutex
}

func (t *CacheTimeOutClear) init(ctx context.Context) error {
	t.initLock.Do(func() {
		t.mq = new(priority_mq.PriorityMq)
		t.keyMap = make(map[string]*priority_mq.ItemNode)

		go util.FuncUtil.HandlePanic(context.Background(), func() {
			t.Clearing(context.Background())
		})()

	})
	return nil
}

func (t *CacheTimeOutClear) GetTopPriority(ctx context.Context) int64 {
	t.init(ctx)
	return t.mq.GetTopPriority()
}

func (t *CacheTimeOutClear) Pop(ctx context.Context) *priority_mq.ItemNode {
	t.init(ctx)
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.mq.Len() <= 0 {
		return nil
	}
	n := heap.Pop(t.mq)
	result := n.(*priority_mq.ItemNode)
	key := result.Value.(string)
	delete(t.keyMap, key)
	return result
}

func (t *CacheTimeOutClear) Push(ctx context.Context, key string, priority int64) {
	t.init(ctx)
	t.lock.Lock()
	defer t.lock.Unlock()
	if n1, ok := t.keyMap[key]; ok {
		n1.Priority = priority
		heap.Fix(t.mq, n1.Index)
		return
	}

	node := &priority_mq.ItemNode{
		Value:    key,
		Priority: priority,
	}
	heap.Push(t.mq, node)
	t.keyMap[key] = node
}

func (t *CacheTimeOutClear) Remove(ctx context.Context, key string) {
	t.init(ctx)
	t.lock.Lock()
	defer t.lock.Unlock()
	if _, ok := t.keyMap[key]; ok {
		return
	}
	heap.Remove(t.mq, t.keyMap[key].Index)
	delete(t.keyMap, key)
}

func (t *CacheTimeOutClear) Clearing(ctx context.Context) error {
	t1 := time.NewTicker(time.Second * 5)
	for _ = range t1.C {
		util.FuncUtil.HandlePanic(ctx, func() {
			for i := 0; i < 5; i++ {
				m := t.mq.GetTopPriority()
				if m <= 0 || m > time.Now().Unix() {
					break
				}
				t.clearTopKey(ctx)
			}
		})()
	}
	return nil
}
func (t *CacheTimeOutClear) clearTopKey(ctx context.Context) {
	target := t.Pop(ctx)
	if target == nil {
		return
	}
	key := target.Value.(string)
	t.CachePool.Remove(key)
}
