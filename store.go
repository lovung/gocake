package gocake

import (
	"sync"
)

// store saves the cache data
type store struct {
	lock sync.RWMutex
	data map[uint64]interface{}
}

func newStore() *store {
	return &store{
		lock: sync.RWMutex{},
		data: make(map[uint64]interface{}),
	}
}

func (s *store) get(key uint64) (interface{}, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

func (s *store) set(key uint64, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[key] = value
}

func (s *store) del(key uint64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// TODO: improve GC overhead for this operation
	delete(s.data, key)
}

func (s *store) clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	// TODO: improve GC overhead for this operation
	s.data = make(map[uint64]interface{})
}

func (s *store) getMany(mapItems map[uint64]interface{}) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for key, value := range mapItems {
		if _, ok := s.data[key]; !ok {
			mapItems[key] = nil
		}
		value = s.data[key]
		mapItems[key] = value
	}
}

func (s *store) setMany(mapItems map[uint64]interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for key, value := range mapItems {
		s.data[key] = value
	}
}
