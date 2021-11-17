package gocake

import (
	"fmt"
	"testing"
	"time"

	"github.com/dgraph-io/ristretto"
)

func toKey(i int) string {
	return fmt.Sprintf("item:%d", i)
}

func BenchmarkRistretto(b *testing.B) {
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cache.Set(toKey(i), toKey(i), 1)
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value, ok := cache.Get(toKey(i))
			if ok {
				_ = value
			}
		}
	})

}

func BenchmarkGocake(b *testing.B) {
	cache := NewCache(1e7, nil)

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cache.Set(toKey(i), toKey(i))
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value := cache.Get(toKey(i))
			if value != nil {
				_ = value
			}
		}
	})
}

func TestExpire(t *testing.T) {
	cache := NewCache(10, nil)
	cache.SetWithTTL(toKey(2), toKey(2), time.Millisecond)
	t.Run(
		"expire", func(t *testing.T) {
			time.Sleep(2 * time.Millisecond)
			if cache.Get(toKey(2)) != nil {
				t.Error("Expire failed")
			}
		},
	)
}
