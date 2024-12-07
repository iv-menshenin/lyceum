package sparseset

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSparseSet(t *testing.T) {
	sp := New[int64, string]()

	sp.Set(10, "foo")
	sp.Set(15, "bar")
	sp.Set(25, "baz")
	sp.Set(23, "qux")
	sp.Set(34, "quux")
	sp.Set(56, "corge")
	sp.Set(78, "grault")

	sp.Delete(23)

	sp.Set(1, "garply")
	sp.Set(44, "waldo")
	sp.Set(51, "fred")

	sp.Delete(10)

	sp.Set(0, "plugh")
	sp.Set(88, "xyzzy")
	sp.Set(91, "thud")

	sp.Delete(88)

	require.Equal(t, "fred", *sp.Get(51))
	require.Equal(t, "baz", *sp.Get(25))
	require.Equal(t, "corge", *sp.Get(56))
	require.Equal(t, "thud", *sp.Get(91))

	sp.Delete(91)

	require.Nil(t, sp.Get(10))
	require.Nil(t, sp.Get(23))
	require.Nil(t, sp.Get(88))
	require.Nil(t, sp.Get(91))
}

func TestSparseSetMass(t *testing.T) {
	sp := New[int64, string]()

	const count = 1000000
	for n := 0; n < count; n += 2 {
		sp.Set(int64(n), strconv.Itoa(n))
		sp.Set(int64(n+1), strconv.Itoa(n+1))
	}

	// test
	for n := 0; n < count; n++ {
		val := sp.Get(int64(n))
		if expected := strconv.Itoa(n); *val != expected {
			t.Errorf("expected %q, got %q", expected, *val)
		}
	}

	sp.Set(count-1, "last")

	// delete all thirds
	for n := 3; n < count; n += 3 {
		if n == count-1 {
			continue
		}
		sp.Delete(int64(n))
	}

	// replace all fifths
	for n := 5; n < count; n += 5 {
		if n == count-1 {
			continue
		}
		sp.Set(int64(n), "foo"+strconv.Itoa(n))
	}

	// test again
	for n := 0; n < count; n++ {
		val := sp.Get(int64(n))
		switch {
		case n == count-1:
			if expected := "last"; *val != expected {
				t.Errorf("expected %q, got %q", expected, *val)
			}
		case n > 0 && n%5 == 0:
			if expected := "foo" + strconv.Itoa(n); *val != expected {
				t.Errorf("expected %q, got %q", expected, *val)
			}
		case n > 0 && n%3 == 0:
			if val != nil {
				t.Errorf("expected nil, got %q", *val)
			}
		default:
			if expected := strconv.Itoa(n); *val != expected {
				t.Errorf("expected %q, got %q", expected, *val)
			}
		}
	}
}

func TestSparseSetEach(t *testing.T) {
	t.Run("100_000", func(t *testing.T) {
		type some struct {
			X, Y, Z int
		}
		sp := New[int, some]()
		for i := 0; i < 100_000; i++ {
			sp.Set(i, some{})
		}

		started := time.Now()
		var y int
		sp.Each(func(_ int, val *some) bool {
			val.X = 1
			val.Y = y
			val.Z = -y
			y++
			return true
		})
		t.Logf("done 100k with %v", time.Since(started))

		// check
		for i := 0; i < 100_000; i++ {
			v := sp.Get(i)
			require.NotNil(t, v)
			require.Equal(t, 1, v.X)
			require.Equal(t, i, v.Y)
			require.Equal(t, -i, v.Z)
		}
	})
	t.Run("1_000_000", func(t *testing.T) {
		type some struct {
			X, Y, Z int
		}
		sp := New[int, some]()
		for i := 0; i < 1_000_000; i++ {
			sp.Set(i, some{})
		}

		started := time.Now()
		var y int
		sp.Each(func(_ int, val *some) bool {
			val.X = 1
			val.Y = y
			val.Z = -y
			y++
			return true
		})
		t.Logf("done 1m with %v", time.Since(started))

		// check
		for i := 0; i < 1_000_000; i++ {
			v := sp.Get(i)
			require.NotNil(t, v)
			require.Equal(t, 1, v.X)
			require.Equal(t, i, v.Y)
			require.Equal(t, -i, v.Z)
		}
	})
	t.Run("16_000_000", func(t *testing.T) {
		type some struct {
			X, Y, Z int
		}
		sp := New[int, some]()
		for i := 0; i < 16_000_000; i++ {
			sp.Set(i, some{})
		}

		started := time.Now()
		var y int
		sp.Each(func(_ int, val *some) bool {
			val.X = 1
			val.Y = y
			val.Z = -y
			y++
			return true
		})
		t.Logf("done 16m with %v", time.Since(started))
	})
}

func BenchmarkSparseSet(b *testing.B) {
	b.Run("insert", func(b *testing.B) {
		sp := New[int, string]()
		for i := 0; i < 1000000; i++ {
			sp.Set(i, fmt.Sprintf("inited-%d", i))
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sp.Set(i, fmt.Sprintf("inserted-%d", i))
		}
	})

	b.Run("delete", func(b *testing.B) {
		sp := New[int, string]()
		for i := 0; i < b.N; i++ {
			sp.Set(i, strconv.Itoa(i))
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sp.Delete(i)
		}
	})
}

func BenchmarkSparseSetDelete(b *testing.B) {
	sp := New[int, string]()
	for i := 0; i < b.N; i++ {
		sp.Set(i, fmt.Sprintf("inited-%d", i))
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sp.Delete(i)
	}
}

func BenchmarkSparseSetEach(b *testing.B) {
	b.Run("by_count", func(b *testing.B) {
		sp := New[int, uint64]()
		for i := 0; i < b.N; i++ {
			sp.Set(i, uint64(i))
		}
		b.ReportAllocs()
		b.ResetTimer()
		sp.Each(func(_ int, val *uint64) bool {
			*val = *val + 12
			return true
		})
	})
	b.Run("static_size_1_000_000", func(b *testing.B) {
		sp := New[int, uint64]()
		for i := 0; i < 1_000_000; i++ {
			sp.Set(i, uint64(i))
		}
		b.ReportAllocs()
		b.ResetTimer()
		var left = b.N
		for left > 0 {
			sp.Each(func(_ int, val *uint64) bool {
				left--
				*val = *val + 12
				return left > 0
			})
		}
	})
}

func BenchmarkSparseSetEach18m(b *testing.B) {
	sp := New[int, uint64]()
	for i := 0; i < 18_000_000; i++ {
		sp.Set(i, uint64(i))
	}
	b.ResetTimer()
	b.Run("static_size", func(b *testing.B) {
		b.ReportAllocs()
		var left = b.N
		for left > 0 {
			sp.Each(func(_ int, val *uint64) bool {
				left--
				*val = *val + 12
				return left > 0
			})
		}
	})
}
