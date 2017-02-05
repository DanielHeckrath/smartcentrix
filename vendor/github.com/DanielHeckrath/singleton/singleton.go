package singleton

import (
	"sync"
	"sync/atomic"
)

// Singleton is an object that will try to call a factory function until it succedes
type Singleton struct {
	entity interface{}
	m      sync.Mutex
	done   uint32
}

// Get calls the function f if and only if previous calls to f have returned an error.
// In other words, given
// 	var singleton Singleton
// if singleton.Get(f) is called multiple times, f will be invoke only if previous calls to f
// have returned an error, even if f has a different value in each invocation.
//
// Get is intended for initialization that can possibly fail.
//
// Because no call to Get returns until the one call to f returns, if f causes
// Get to be called, it will deadlock.
func (s *Singleton) Get(f func() (interface{}, error)) (interface{}, error) {
	if atomic.LoadUint32(&s.done) == 1 {
		return s.entity, nil
	}
	// Slow-path.
	s.m.Lock()
	defer s.m.Unlock()
	if s.done == 0 {
		entity, err := f()
		if err != nil {
			return nil, err
		}
		defer atomic.StoreUint32(&s.done, 1)
		s.entity = entity
	}
	return s.entity, nil
}
