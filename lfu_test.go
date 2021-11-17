package gocake

import (
	"math/rand"
	"testing"
)

func TestLFU(t *testing.T) {
	lfuObj := newLFU()
	lfuObj.touch(1)
	lfuObj.touch(2)
	lfuObj.touch(1)

	if lfuObj.String() != "1: 2 - 2: 1 " {
		t.Errorf("lfuObj.String() = %s", lfuObj.String())
	}
	lfuObj.touch(3)
	if lfuObj.String() != "1: 3 2 - 2: 1 " && lfuObj.String() != "1: 2 3 - 2: 1 " {
		t.Errorf("lfuObj.String() = %s", lfuObj.String())
	}
	lfuObj.touch(1)
	if lfuObj.String() != "1: 3 2 - 3: 1 " && lfuObj.String() != "1: 2 3 - 3: 1 " {
		t.Errorf("lfuObj.String() = %s", lfuObj.String())
	}
	lfuObj.touch(2)
	if lfuObj.String() != "1: 3 - 2: 2 - 3: 1 " {
		t.Errorf("lfuObj.String() = %s", lfuObj.String())
	}
	if got := lfuObj.clean(1); got[0] != 3 && lfuObj.String() != "2: 2 - 3: 1 " {
		t.Errorf("lfuObj.String() = %s", lfuObj.String())
	}
	lfuObj.touch(4)
	if lfuObj.String() != "1: 4 - 2: 2 - 3: 1 " {
		t.Errorf("lfuObj.String() = %s", lfuObj.String())
	}
	lfuObj.clean(2)
	if lfuObj.String() != "3: 1 " {
		t.Errorf("lfuObj.String() = %s", lfuObj.String())
	}
}

func BenchmarkLFU(b *testing.B) {
	c := newLFU()

	b.Run("Touch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.touch(rand.Uint64())
		}
	})

	b.Run("Clean", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.clean(1)
		}
	})
}
