package gocake

import (
	"errors"
	"time"
)

const shardCount = 64

// All errors of the library.
var (
	ErrCacheMiss = errors.New("cache miss")
)

type cache struct {
	shards     []*store
	lfu        *lfu
	size       int
	hasher     Hasher
	expireChan <-chan uint64
}

// NewCache creates a new cache with the given capacity.
func NewCache(size int, hasher Hasher) *cache {
	c := &cache{}
	c.lfu = newLFU()
	c.size = size
	c.shards = make([]*store, shardCount)
	expireChan := make(chan uint64, 1024)
	c.expireChan = expireChan
	for i := 0; i < shardCount; i++ {
		c.shards[i] = newStore(expireChan)
	}
	if hasher == nil {
		c.hasher = DefaultHasher()
	}
	go c.cleanExpired()
	return c
}

func (c *cache) Get(key string) interface{} {
	ukey := c.hasher.Sum64(key)
	shard := c.shards[ukey%shardCount]
	if v, ok := shard.get(ukey); ok {
		c.lfu.touch(ukey)
		return v
	}
	return nil
}

// func (c *cache) GetMany(mapItems map[string]interface{}) {
// 	return nil
// }

func (c *cache) Set(key string, value interface{}) {
	ukey := c.hasher.Sum64(key)
	shard := c.shards[ukey%shardCount]
	shard.set(ukey, value, 0)
	c.lfu.touch(ukey)
}

func (c *cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	ukey := c.hasher.Sum64(key)
	shard := c.shards[ukey%shardCount]
	shard.set(ukey, value, ttl)
	c.lfu.touch(ukey)
}

// func (c *cache) SetMany(mapItems map[string]interface{})
// func (c *cache) SetManyWithTTL(mapItems map[string]interface{}, ttl time.Duration)
// func (c *cache) Delete(key string)
// func (c *cache) DeleteMany(key []string)
// func (c *cache) Clear()

func (c *cache) cleanExpired() {
	for e := range c.expireChan {
		c.lfu.del(e)
	}
}
