package priority_mq

import (
	"context"
	"sync"
)

/*
TimeOutClear
定时清理 权重就是时间搓
*/
type TimeOutClear struct {
	mq       *PriorityMq
	initLock sync.Once
}

func (t *TimeOutClear) init(ctx context.Context) error {
	t.initLock.Do(func() {

	})
	return nil
}

func (t *TimeOutClear) AddItem(ctx context.Context, key string, priority int64) {

}
