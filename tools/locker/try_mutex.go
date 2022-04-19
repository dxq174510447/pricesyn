package locker

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const mutexLocked = 1 << iota // mutex is locked
/**
在原有的基础上 加上try_lock功能
*/
type TryMutex struct {
	sync.Mutex
}

func (m *TryMutex) TryLock() bool {
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked)
}
