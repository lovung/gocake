package gocake

import (
	"errors"
)

const shardCount = 64

var (
	ErrCacheMiss = errors.New("cache miss")
)

type cache struct {
	shards []*store
	lfu    *lfu
	size   int
	hasher Hasher
}

// NewCache creates a new cache with the given capacity.
func NewCache(size int, hasher Hasher) *cache {
	c := &cache{}
	c.lfu = newLFU()
	c.size = size
	c.shards = make([]*store, shardCount)
	for i := 0; i < shardCount; i++ {
		c.shards[i] = newStore()
	}
	if hasher == nil {
		c.hasher = DefaultHasher()
	}
	return c
}

func (c *cache) Get(key string) (interface{}, error) {
	ukey := c.hasher.Sum64(key)
	shard := c.shards[ukey%uint64(shardCount)]
	if v, ok := shard.get(ukey); ok {
		c.lfu.touch(ukey)
		return v, nil
	}
	return nil, ErrCacheMiss
}

// func (c *cache) GetMany(mapItems map[string]interface{}) error {
// 	return nil
// }

func (c *cache) Set(key string, value interface{}) error {
	ukey := c.hasher.Sum64(key)
	shard := c.shards[ukey%uint64(shardCount)]
	shard.set(ukey, value)
	c.lfu.touch(ukey)
	return nil
}

// func (c *cache) SetWithTTL(key string, value interface{}, ttl time.Duration) error
// func (c *cache) SetMany(mapItems map[string]interface{}) error
// func (c *cache) SetManyWithTTL(mapItems map[string]interface{}, ttl time.Duration) error
// func (c *cache) Delete(key string) error
// func (c *cache) DeleteMany(key []string) error
// func (c *cache) Clear() error
