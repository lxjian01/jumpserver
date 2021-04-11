package utils

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	LockedFlag   int32 = 1
	UnlockedFlag int32 = 0
)

type Mutex struct {
	in      sync.Mutex
	status  *int32
	opsLock sync.Mutex
	// 用于多线程控制采集任务数量不能超过2(1正在采集，2进行自旋，后续事件判读ops为2直接返回)
	ops int32
}

func NewMutex() *Mutex {
	status := UnlockedFlag
	return &Mutex{
		status: &status,
		ops:    0,
	}
}

func (m *Mutex) Lock() {
	m.in.Lock()
	atomic.AddInt32(m.status, LockedFlag)
}

func (m *Mutex) Unlock() {
	m.in.Unlock()
	atomic.AddInt32(m.status, -LockedFlag)
}

func (m *Mutex) TryLock() bool {
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.in)), UnlockedFlag, LockedFlag) {
		atomic.AddInt32(m.status, LockedFlag)
		return true
	}
	return false
}

func (m *Mutex) IsLocked() bool {
	if atomic.LoadInt32(m.status) == LockedFlag {
		return true
	}
	return false
}

func (m *Mutex) TryAcquire() bool {
	m.opsLock.Lock()
	if m.loadOps() >= 2 {
		m.opsLock.Unlock()
		return false
	} else if m.loadOps() == 1 {
		m.addOps(1)
		m.opsLock.Unlock()
		//最多自旋1000次,100秒
		for i := 0; i < 1000; i++ {
			//进行自旋
			if m.loadOps() == 1 {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		return true
	} else if m.loadOps() == 0 {
		m.addOps(1)
		m.opsLock.Unlock()
		return true
	} else {
		m.opsLock.Unlock()
		return true
	}
}

func (m *Mutex) Release() {
	m.addOps(-1)
}

func (m *Mutex) loadOps() int32 {
	m.in.Lock()
	defer m.in.Unlock()
	return atomic.LoadInt32(&m.ops)
}

func (m *Mutex) addOps(delta int32) {
	m.in.Lock()
	defer m.in.Unlock()
	atomic.AddInt32(&m.ops, delta)
}
