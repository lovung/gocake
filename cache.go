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
	shards []*store
	hasher Hasher
}

// NewCache creates a new cache with the given capacity.
func NewCache(size int, hasher Hasher) *cache {
	c := &cache{}
	c.shards = make([]*store, shardCount)
	for i := 0; i < shardCount; i++ {
		c.shards[i] = newStore()
	}
	if hasher == nil {
		c.hasher = DefaultHasher()
	}
	return c
}

func (c *cache) Get(key string) interface{} {
	ukey := c.hasher.Sum64(key)
	shard := c.shards[ukey%shardCount]
	if v, ok := shard.get(ukey); ok {
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
}

func (c *cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	ukey := c.hasher.Sum64(key)
	shard := c.shards[ukey%shardCount]
	shard.set(ukey, value, ttl)
}

// func (c *cache) SetMany(mapItems map[string]interface{})
// func (c *cache) SetManyWithTTL(mapItems map[string]interface{}, ttl time.Duration)
// func (c *cache) Delete(key string)
// func (c *cache) DeleteMany(key []string)
// func (c *cache) Clear()
