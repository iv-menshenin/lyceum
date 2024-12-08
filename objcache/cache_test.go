package objcache

import (
	"strconv"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestCacheReusable(t *testing.T) {
	t.Run("10", func(t *testing.T) {
		testNoAllocations(t, 10)
	})
	t.Run("1_000", func(t *testing.T) {
		testNoAllocations(t, 1_000)
	})
	t.Run("1_000_000", func(t *testing.T) {
		testNoAllocations(t, 1_000_000)
	})
}

func testNoAllocations(t *testing.T, count int) {
	t.Helper()

	type T struct {
		n int
		x int
	}
	var c = &Cache[T]{}
	for n := 0; n < count; n++ {
		_ = c.Get()
	}
	// heap leaking, do not remove
	var v *T

	a := testing.AllocsPerRun(10, func() {
		c.Clear()
		for n := 0; n < count; n++ {
			v = c.Get()
			v.n = 0
		}
	})
	require.Equal(t, float64(0), a)
}

func TestCache(t *testing.T) {
	t.Run("sizeRestrictions", func(t *testing.T) {
		require.Less(t, uint32(unsafe.Sizeof(Cache[int]{})), uint32(1024))
	})
	t.Run("the_same_address", func(t *testing.T) {
		var c Cache[int]
		v1 := c.Get()
		c.Clear()
		v2 := c.Get()
		// v1 and v2 must be the same address
		*v2 = 2004
		require.Equal(t, 2004, *v1)
		*v1 = 1003
		require.Equal(t, 1003, *v2)
	})
	t.Run("check_no_repetative", func(t *testing.T) {
		const count = 10_000
		var c Cache[string]
		var m = make(map[string]struct{}, count)
		var v = make([]*string, count)
		for n := 0; n < count; n++ {
			v[n] = c.Get()
			*v[n] = strconv.Itoa(n)
		}
		for _, vs := range v {
			_, ok := m[*vs]
			require.False(t, ok, "memory conflict on %q", *vs)
			m[*vs] = struct{}{}
		}
	})
}

func BenchmarkCache(b *testing.B) {
	type T struct {
		n int
	}
	// heap leaking, do not remove
	var (
		v1000    *T
		v10000   *T
		v1000000 *T
	)
	b.Run("1000", func(b *testing.B) {
		b.ReportAllocs()
		var c Cache[T]
		for n := 0; n < b.N; n++ {
			v1000 = c.Get()
			v1000.n = 1
			if (1+n)%1000 == 0 {
				c.Clear()
			}
		}
	})
	b.Run("10_000", func(b *testing.B) {
		b.ReportAllocs()
		var c Cache[T]
		for n := 0; n < b.N; n++ {
			v10000 = c.Get()
			v10000.n = 1
			if (1+n)%10000 == 0 {
				c.Clear()
			}
		}
	})
	b.Run("1000_000", func(b *testing.B) {
		b.ReportAllocs()
		var c Cache[T]
		for n := 0; n < b.N; n++ {
			v1000000 = c.Get()
			v1000000.n = 1
			if (1+n)%1000_000 == 0 {
				c.Clear()
			}
		}
	})
}

func BenchmarkCacheStatic(b *testing.B) {
	type T struct {
		n int
	}
	// heap leaking, do not remove
	var (
		v1000    *T
		v10000   *T
		v1000000 *T
	)
	b.Run("1000", func(b *testing.B) {
		b.ReportAllocs()
		var c Cache[T]
		for n := 0; n < b.N; n++ {
			for m := 0; m < 1000; m++ {
				v1000 = c.Get()
				v1000.n = 1
			}
			c.Clear()
		}
	})
	b.Run("10_000", func(b *testing.B) {
		b.ReportAllocs()
		var c Cache[T]
		for n := 0; n < b.N; n++ {
			for m := 0; m < 10_000; m++ {
				v10000 = c.Get()
				v10000.n = 1
			}
			c.Clear()
		}
	})
	b.Run("1000_000", func(b *testing.B) {
		b.ReportAllocs()
		var c Cache[T]
		for n := 0; n < b.N; n++ {
			for m := 0; m < 1_000_000; m++ {
				v1000000 = c.Get()
				v1000000.n = 1
			}
			c.Clear()
		}
	})
}
