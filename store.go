package gocake

import (
	"sync"
	"time"
)

// Assume 1 item ~ 1KB
// maxStoreCount = 100MB / 1KB = 100,000
const maxStoreCount = 10000

// store saves the cache data
type store struct {
	lock  sync.RWMutex
	data  map[uint64]storeItem
	count int
	lfu   *lfu
}

type storeItem struct {
	value     interface{}
	expiredAt int64
}

func newStore() *store {
	return &store{
		lock: sync.RWMutex{},
		data: make(map[uint64]storeItem),
		lfu:  newLFU(),
	}
}

func (s *store) get(key uint64) (interface{}, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	item, ok := s.data[key]
	if !ok {
		return nil, false
	}
	if item.expiredAt > 0 && item.expiredAt < time.Now().UnixNano() {
		s.lfu.del(key)
		s.count--
		delete(s.data, key)
		return nil, false
	}
	s.lfu.touch(key)
	return item.value, true
}

func (s *store) set(key uint64, value interface{}, ttl time.Duration) {
	s.lock.Lock()
	defer s.lock.Unlock()
	var expiredAt int64
	if ttl != 0 {
		expiredAt = time.Now().Add(ttl).UnixNano()
	}
	item, ok := s.data[key]
	if !ok {
		// not exist
		if s.count > maxStoreCount {
			// use LFU to evict item
			delKeys := s.lfu.clean(1)
			for _, key := range delKeys {
				delete(s.data, key)
				s.count--
			}
		}
		s.count++
		s.lfu.touch(key)
		s.data[key] = storeItem{
			value:     value,
			expiredAt: expiredAt,
		}
		return
	}
	if item.expiredAt > 0 && item.expiredAt < time.Now().UnixNano() {
		// already expired: reset the freq count
		s.lfu.del(key)
	}
	s.lfu.touch(key)
	item.value = value
	item.expiredAt = expiredAt
}

func (s *store) del(key uint64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// TODO: improve GC overhead for this operation
	if _, ok := s.data[key]; ok {
		s.lfu.del(key)
		s.count--
		delete(s.data, key)
	}
}

func (s *store) clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	// TODO: improve GC overhead for this operation
	for k := range s.data {
		s.lfu.del(k)
		s.count--
		delete(s.data, k)
	}
	s.data = make(map[uint64]storeItem)
}

// Not be used
func (s *store) getMany(mapItems map[uint64]interface{}) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	now := time.Now().Unix()
	for key := range mapItems {
		item, ok := s.data[key]
		if !ok {
			mapItems[key] = nil
			continue
		}
		if item.expiredAt > 0 && item.expiredAt < now {
			s.lfu.del(key)
			s.count--
			mapItems[key] = nil
			continue
		}
		s.lfu.touch(key)
		mapItems[key] = s.data[key].value
	}
}

func (s *store) setMany(mapItems map[uint64]interface{}, ttl time.Duration) {
	s.lock.Lock()
	defer s.lock.Unlock()
	var expiredAt int64
	if ttl != 0 {
		expiredAt = time.Now().Add(ttl).UnixNano()
	}
	for key, value := range mapItems {
		item, ok := s.data[key]
		if !ok {
			s.count++
			s.data[key] = storeItem{
				value:     value,
				expiredAt: expiredAt,
			}
			return
		}
		if item.expiredAt > 0 && item.expiredAt < time.Now().UnixNano() {
			// already expired: reset the freq count
			s.lfu.del(key)
		}
		item.value = value
		item.expiredAt = expiredAt
	}

	if overSize := s.count - maxStoreCount; overSize > 0 {
		// use LFU to evict items
		delKeys := s.lfu.clean(overSize)
		for _, key := range delKeys {
			delete(s.data, key)
			s.count--
		}
	}
	for key := range mapItems {
		s.lfu.touch(key)
	}
}
