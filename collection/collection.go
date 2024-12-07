package collection

type (
	// Collection is a special data structure designed to avoid excessive memory allocation when adding
	// a large number of values to a slice in situations where you do not know in advance the amount of data
	// that needs to be stored.
	//
	// The collection guarantees correct sorting as long as you do not remove a value from the middle.
	// When a value is removed, all subsequent values are not shifted; instead, the current value is replaced
	// with the last one, and the length is reduced by one.
	//
	// The Push and Get methods return a reference to an object that can be modified. Note that the reference
	// is guaranteed to be valid only until the first call to methods that delete values, such as Delete or even Pop.
	// Avoid storing the reference for a long time.
	//
	// If you use a power of two as the bucket size, lightweight bit-shifting and bit-masking operations will be applied
	// for calculating read/write addresses, significantly improving performance
	Collection[T any] struct {
		len     int
		bsz     int
		bShift  int
		xMask   int
		buckets []*bucket[T]
	}
	bucket[T any] struct {
		data []T
	}
)

const defaultBucketSz = 1024

// New creates a new Collection with the specified bucket size. If the size is zero, the default value will be used.
//
// If you use a power of two as the bucket size, lightweight bit-shifting and bit-masking operations will be applied
// for calculating read/write addresses, significantly improving performance
func New[T any]() *Collection[T] {
	var c Collection[T]
	c.initBucketSize(defaultBucketSz)
	return &c
}

func (c *Collection[T]) initBucketSize(bsz int) {
	if bsz == 0 {
		bsz = defaultBucketSz
	}
	c.bsz = bsz
	// we can make all faster if bsz is power of two
	var bits, mask int
	for bsz > 0 {
		if bsz == 1 {
			break
		}
		if bsz&1 > 0 {
			// wasted
			return
		}
		bsz >>= 1
		bits += 1
		mask = (mask << 1) + 1
	}
	c.bShift = bits
	c.xMask = mask
}

func (c *Collection[T]) Len() int {
	return c.len
}

// Push adds a new value to the end of the Collection and returns a reference to it.
func (c *Collection[T]) Push(val T) *T {
	if c.bsz == 0 {
		c.initBucketSize(defaultBucketSz)
	}
	id := c.len
	var xId, bId int
	if c.xMask > 0 {
		bId = id >> c.bShift
		xId = id & c.xMask
	} else {
		xId = id % c.bsz
		bId = id / c.bsz
	}
	c.len++
	if len(c.buckets) <= bId {
		c.extendBuckets()
	}
	c.buckets[bId].data[xId] = val
	return &c.buckets[bId].data[xId]
}

func (c *Collection[T]) extendBuckets() {
	c.buckets = append(c.buckets, &bucket[T]{
		data: make([]T, c.bsz),
	})
}

// Get allows you to get a reference to an object located in a Collection.
//
// Avoid storing the link outside the Collection for long periods of time.
func (c *Collection[T]) Get(id int) *T {
	if id >= c.len {
		return nil
	}
	var xId, bId int
	if c.xMask > 0 {
		bId = id >> c.bShift
		xId = id & c.xMask
	} else {
		xId = id % c.bsz
		bId = id / c.bsz
	}
	return &c.buckets[bId].data[xId]
}

// Delete deletes an object by its index from the collection.
//
// Note that to improve performance, there is a side effect: the deleted object is replaced by the last object,
// not the next in line. This avoids large data movement when deleting values from the beginning.
//
// So if you need to delete several values from n to m, it is safe to do it only in the index decreasing direction,
// i.e. from m to n.
func (c *Collection[T]) Delete(id int) {
	if id >= c.len {
		panic("out of bounds")
	}
	var xId, bId int
	if c.xMask > 0 {
		bId = id >> c.bShift
		xId = id & c.xMask
	} else {
		xId = id % c.bsz
		bId = id / c.bsz
	}
	lId := c.len - 1
	var xxId, xbId int
	if c.xMask > 0 {
		xbId = lId >> c.bShift
		xxId = lId & c.xMask
	} else {
		xxId = lId % c.bsz
		xbId = lId / c.bsz
	}
	if bId != xbId || xId != xxId {
		// swap
		c.buckets[bId].data[xId] = c.buckets[xbId].data[xxId]
	}
	c.len--

	// clear cell
	var empty T
	c.buckets[xbId].data[xxId] = empty
}

// Pop selects the last item in the collection and returns a copy of it. The original item is deleted.
func (c *Collection[T]) Pop() T {
	if c.len < 1 {
		panic("out of bounds")
	}
	id := c.len - 1
	var xId, bId int
	if c.xMask > 0 {
		bId = id >> c.bShift
		xId = id & c.xMask
	} else {
		xId = id % c.bsz
		bId = id / c.bsz
	}
	c.len--
	b := c.buckets[bId]
	val := b.data[xId]
	// clean cell
	var empty T
	b.data[xId] = empty
	return val
}

// Prune clears unoccupied space. It can be used after a large number of calls to Delete or Pop method.
func (c *Collection[T]) Prune() {
	if c.bsz == 0 {
		c.initBucketSize(defaultBucketSz)
	}
	bId := c.len / c.bsz
	for n := bId + 1; n < len(c.buckets); n++ {
		c.buckets[n] = nil
	}
	c.buckets = c.buckets[:bId]
}
