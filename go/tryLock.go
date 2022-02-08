package tool

import "sync/atomic"

type TryLock struct {
	lock uint32
}

func (t *TryLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&t.lock, 0, 1)
}

func (t *TryLock) Unlock() {
	atomic.CompareAndSwapUint32(&t.lock, 1, 0)
}
