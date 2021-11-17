package gocake

import (
	"github.com/zeebo/xxh3"
)

type Hasher interface {
	Sum64(data string) uint64
}

func DefaultHasher() Hasher {
	return &defaultHasher{}
}

type defaultHasher struct{}

func (h *defaultHasher) Sum64(data string) uint64 {
	return xxh3.HashString(data)
}
