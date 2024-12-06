package collection

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFusionCollection(t *testing.T) {
	t.Parallel()
	t.Run("one_append", func(t *testing.T) {
		t.Parallel()

		var c = New[string]()
		v := c.Push("test")
		require.Equal(t, 1, c.Len(), "length checking")
		require.Equal(t, "test", *v, "value checking")
	})
	t.Run("push_few", func(t *testing.T) {
		t.Parallel()

		var c = New[string]()
		c.Push("101")
		c.Push("102")
		c.Push("103")
		require.Equal(t, "103", *c.Get(2))
		require.Equal(t, "102", *c.Get(1))
		require.Equal(t, "101", *c.Get(0))
	})
	t.Run("len", func(t *testing.T) {
		t.Parallel()

		var c = New[int64]()
		require.Equal(t, 0, c.Len())
		c.Push(101)
		require.Equal(t, 1, c.Len())
		c.Push(102)
		require.Equal(t, 2, c.Len())
		c.Push(103)
		require.Equal(t, 3, c.Len())
	})
	t.Run("force", func(t *testing.T) {
		t.Parallel()
		type Elem struct {
			i int
			s string
		}
		var c = New[Elem]()

		const elemCount = 100000

		for n := 0; n < elemCount; n++ {
			require.Equal(t, n, c.Len())
			c.Push(Elem{i: n, s: strconv.Itoa(n)})
		}

		// forward
		for n := 0; n < elemCount; n++ {
			e := c.Get(n)
			require.NotNil(t, e)
			require.Equal(t, Elem{i: n, s: strconv.Itoa(n)}, *e, "value checking")
		}

		// backward
		for n := elemCount; n > 0; n-- {
			e := c.Get(n - 1)
			require.NotNil(t, e)
			require.Equal(t, Elem{i: n - 1, s: strconv.Itoa(n - 1)}, *e, "value checking")
		}
	})
	t.Run("random-get", func(t *testing.T) {
		t.Parallel()

		var c = New[int]()
		_ = c.Push(67)
		_ = c.Push(13)
		_ = c.Push(54)
		_ = c.Push(2)
		_ = c.Push(1)
		_ = c.Push(42)
		_ = c.Push(43)

		require.Equal(t, *c.Get(4), 1)
		require.Equal(t, *c.Get(2), 54)
		require.Equal(t, *c.Get(5), 42)
		require.Equal(t, *c.Get(0), 67)
		require.Equal(t, *c.Get(6), 43)

		c.Push(2)
		require.Equal(t, *c.Get(7), 2)

		c.Push(67)
		require.Equal(t, *c.Get(2), 54)
	})
	t.Run("append_get_delete", func(t *testing.T) {
		t.Parallel()

		var c = New[int64]()
		c.Push(801)
		c.Push(802)
		c.Push(803)
		require.Equal(t, int64(803), *c.Get(2))
		require.Equal(t, int64(802), *c.Get(1))
		require.Equal(t, int64(801), *c.Get(0))
		c.Delete(1)
		require.Equal(t, int64(803), *c.Get(1))
		require.Equal(t, int64(801), *c.Get(0))
		c.Delete(0)
		c.Delete(0)
		require.Equal(t, 0, c.Len())

		c.Push(0)
		require.Equal(t, int64(0), *c.Get(0))
		require.Nil(t, c.Get(1))
		require.Nil(t, c.Get(2))
	})
}

func TestCollectionPushPop(t *testing.T) {
	t.Parallel()
	t.Run("push_pop", func(t *testing.T) {
		t.Parallel()
		var c = New[string]()

		c.Push("foo")
		require.Equal(t, 1, c.Len())
		c.Push("bar")
		require.Equal(t, 2, c.Len())
		require.Equal(t, "bar", c.Pop())
		require.Equal(t, 1, c.Len())
		require.Equal(t, "foo", c.Pop())
		require.Equal(t, 0, c.Len())

		c.Push("1")
		c.Push("3")
		require.Equal(t, "3", c.Pop())
		c.Push("2")
		c.Push("3")
		c.Push("4")
		require.Equal(t, 4, c.Len())
		require.Equal(t, "4", c.Pop())
		require.Equal(t, "3", c.Pop())
		require.Equal(t, "2", c.Pop())
		require.Equal(t, "1", c.Pop())
		require.Equal(t, 0, c.Len())
	})
	t.Run("force", func(t *testing.T) {
		t.Parallel()

		var c = New[int]()
		const elemCount = 1000000

		for n := 0; n < elemCount; n++ {
			c.Push(n)
		}
		for n := elemCount - 1; n >= 0; n-- {
			require.Equal(t, n, c.Pop())
		}
	})
}

func BenchmarkCollectionGet_10mln(b *testing.B) {
	type Elem struct {
		s          string
		a, b, c, d int64
		n          int
	}
	var c = New[Elem]()
	const count = 10_000_000
	for n := 0; n < count; n++ {
		c.Push(Elem{n: n})
	}
	b.Run("Last", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Get(count - 1)
		}
	})
	b.Run("First", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Get(1)
		}
	})
	b.Run("Mid", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Get(count / 2)
		}
	})
	b.Run("Forward", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Get(n % count)
		}
	})
	b.Run("Backward", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			_ = c.Get(count - (n % count) - 1)
		}
	})
	b.Run("Random", func(b *testing.B) {
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			switch n % 5 {
			case 0:
				// static
				_ = c.Get(count / 2)
			case 1:
				// next
				_ = c.Get(n % count)
			case 2:
				// backward
				_ = c.Get(count - (n % count) - 1)
			case 3:
				// fast-backward
				_ = c.Get(count - ((n * 127) % count) - 1)
			case 4:
				// fast-forward
				_ = c.Get((n * 127) % count)
			}
		}
	})
}

func BenchmarkCollectionPushPop(b *testing.B) {
	type Elem struct {
		s          string
		a, b, c, d int64
		n          int
	}
	b.Run("Push_All", func(b *testing.B) {
		b.ReportAllocs()
		var c = New[Elem]()
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n, s: "Push_All"})
		}
	})
	b.Run("Push_Len_All", func(b *testing.B) {
		b.ReportAllocs()
		var c = New[Elem]()
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n, s: "Push_Len_All"})
			_ = c.Len()
		}
	})
	b.Run("Push_Then_Pop", func(b *testing.B) {
		b.ReportAllocs()
		var c = New[Elem]()
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n, s: "Push_Then_Pop"})
			_ = c.Pop()
		}
	})
	b.Run("Push_Push_Then_Pop", func(b *testing.B) {
		b.ReportAllocs()
		var c = New[Elem]()
		for n := 0; n < b.N; n++ {
			c.Push(Elem{a: 0, n: n, s: "Push_Push_Then_Pop"})
			c.Push(Elem{a: 1, n: n, s: "Push_Push_Then_Pop"})
			_ = c.Pop()
		}
	})
	b.Run("Push_Push_Then_Pop_Len", func(b *testing.B) {
		b.ReportAllocs()
		var c Collection[Elem]
		for n := 0; n < b.N; n++ {
			c.Push(Elem{a: 0, n: n, s: "Push_Push_Then_Pop"})
			c.Push(Elem{a: 1, n: n, s: "Push_Push_Then_Pop"})
			_ = c.Pop()
			_ = c.Len()
		}
	})
	b.Run("Pop_All", func(b *testing.B) {
		b.ReportAllocs()
		var c = New[Elem]()
		for n := 0; n < b.N; n++ {
			c.Push(Elem{n: n, s: "Pop_All"})
		}
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			_ = c.Pop()
		}
	})
}

func BenchmarkCollectionInsert(b *testing.B) {
	type Elem struct {
		s          string
		a, b, c, d int64
		n          int
	}
	b.Run("Push_1K", func(b *testing.B) {
		b.ReportAllocs()
		var c = New[Elem]()
		for n := 0; n < b.N; n++ {
			for i := 0; i < 1_000; i++ {
				c.Push(Elem{n: n, a: int64(i), b: -1, s: "Push 1K"})
			}
		}
	})
	b.Run("Push_65K", func(b *testing.B) {
		b.ReportAllocs()
		var c = New[Elem]()
		for n := 0; n < b.N; n++ {
			for i := 0; i < 65_535; i++ {
				c.Push(Elem{n: n, a: int64(i), b: -1, s: "Push 65K"})
			}
		}
	})
	b.Run("Push_1M", func(b *testing.B) {
		b.ReportAllocs()
		var c = New[Elem]()
		for n := 0; n < b.N; n++ {
			for i := 0; i < 1_000_000; i++ {
				c.Push(Elem{n: n, a: int64(i), b: -1, s: "Push 1m"})
			}
		}
	})
	b.Run("Push_10M", func(b *testing.B) {
		b.ReportAllocs()
		var c = New[Elem]()
		for n := 0; n < b.N; n++ {
			for i := 0; i < 10_000_000; i++ {
				c.Push(Elem{n: n, a: int64(i), b: -1, s: "Push 10m"})
			}
		}
	})
}
