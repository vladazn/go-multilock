package multilock

import (
	"sync"
	"sync/atomic"
)

type RefCounter struct {
	Counter int64
	Lock    sync.Mutex
}

type Multilock struct {
	M  map[string]*RefCounter
	Mx sync.Mutex
}

func (s *Multilock) Lock(key string) {
	s.Mx.Lock()
	locker, ok := s.M[key]
	if !ok {
		newLocker := &RefCounter{
			Counter: 1,
		}
		newLocker.Lock.Lock()
		s.M[key] = newLocker
		s.Mx.Unlock()
	} else {
		s.Mx.Unlock()
		atomic.AddInt64(&locker.Counter, 1)
		locker.Lock.Lock()
	}
}

func (s *Multilock) Unlock(key string) {

	s.Mx.Lock()
	defer s.Mx.Unlock()
	locker, ok := s.M[key]
	if !ok {
		return
	} else {
		atomic.AddInt64(&locker.Counter, -1)
		locker.Lock.Unlock()
		if locker.Counter <= 0 {
			delete(s.M, key)
		}
	}
}

func NewMultipleLock() *Multilock {
	return &Multilock{
		M: make(map[string]*RefCounter),
	}
}

