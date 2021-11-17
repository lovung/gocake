package gocake

import (
	"sync"
	"time"
)

// store saves the cache data
type store struct {
	lock       sync.RWMutex
	data       map[uint64]storeItem
	expireChan chan uint64
}

type storeItem struct {
	value     interface{}
	expiredAt int64
}

func newStore(expireChan chan uint64) *store {
	return &store{
		lock:       sync.RWMutex{},
		data:       make(map[uint64]storeItem),
		expireChan: expireChan,
	}
}

func (s *store) get(key uint64) (interface{}, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	item, ok := s.data[key]
	if item.expiredAt > 0 && item.expiredAt < time.Now().UnixNano() {
		s.expireChan <- key
		delete(s.data, key)
		return nil, false
	}
	return item.value, ok
}

func (s *store) set(key uint64, value interface{}, ttl time.Duration) {
	s.lock.Lock()
	defer s.lock.Unlock()
	item, ok := s.data[key]
	var expiredAt int64
	if ttl != 0 {
		expiredAt = time.Now().Add(ttl).UnixNano()
	}
	if !ok {
		// not exist
		s.data[key] = storeItem{
			value:     value,
			expiredAt: expiredAt,
		}
		return
	}
	// already expired: don't need to remove freq of this key
	item.value = value
	item.expiredAt = expiredAt

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
	s.data = make(map[uint64]storeItem)
}

func (s *store) getMany(mapItems map[uint64]interface{}) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	now := time.Now().Unix()
	for key, _ := range mapItems {
		var item storeItem
		var ok bool
		if item, ok = s.data[key]; !ok {
			mapItems[key] = nil
			continue
		}
		if item.expiredAt > 0 && item.expiredAt < now {
			s.expireChan <- key
			mapItems[key] = nil
			continue
		}
		mapItems[key] = s.data[key].value
	}
}

// func (s *store) setMany(mapItems map[uint64]interface{}) {
// 	s.lock.Lock()
// 	defer s.lock.Unlock()
// 	for key, value := range mapItems {
// 		s.data[key] = storeItem{
// 			value:     value,
// 			expiredAt: 0,
// 		}
// 	}
// }
