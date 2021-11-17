package gocake

import "time"

// Cache is a fast in-memory cache with the optimised storage for avoiding GC overhead.
// It is thread-safe and optimze for maximum hit ratio.
type Cache interface {
	Get(key string) (interface{}, error)
	GetMany(mapItems map[string]interface{}) error
	Set(key string, value interface{}) error
	SetWithTTL(key string, value interface{}, ttl time.Duration) error
	SetMany(mapItems map[string]interface{}) error
	SetManyWithTTL(mapItems map[string]interface{}, ttl time.Duration) error
	Delete(key string) error
	DeleteMany(key []string) error
	Clear() error
}
