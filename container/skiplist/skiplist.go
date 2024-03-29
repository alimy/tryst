package skiplist

type Interface[T any] interface {
	Less(other T) bool
}

type SkipList[T Interface[T]] struct {
	header *Element[T]
	tail   *Element[T]
	update []*Element[T]
	rank   []int
	length int
	level  int
}

// New returns an initialized skiplist.
func New[T Interface[T]]() *SkipList[T] {
	var v T
	return &SkipList[T]{
		header: newElement(SKIPLIST_MAXLEVEL, v),
		tail:   nil,
		update: make([]*Element[T], SKIPLIST_MAXLEVEL),
		rank:   make([]int, SKIPLIST_MAXLEVEL),
		length: 0,
		level:  1,
	}
}

// Init initializes or clears skiplist sl.
func (sl *SkipList[T]) Init() *SkipList[T] {
	var v T
	sl.header, sl.tail = newElement(SKIPLIST_MAXLEVEL, v), nil
	sl.update = make([]*Element[T], SKIPLIST_MAXLEVEL)
	sl.rank = make([]int, SKIPLIST_MAXLEVEL)
	sl.length = 0
	sl.level = 1
	return sl
}

// Front returns the first elements of skiplist sl or nil.
func (sl *SkipList[T]) Front() *Element[T] {
	return sl.header.level[0].forward
}

// Back returns the last elements of skiplist sl or nil.
func (sl *SkipList[T]) Back() *Element[T] {
	return sl.tail
}

// Len returns the numbler of elements of skiplist sl.
func (sl *SkipList[T]) Len() int {
	return sl.length
}

// Insert inserts v, increments sl.length, and returns a new element of wrap v.
func (sl *SkipList[T]) Insert(v T) *Element[T] {
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		// store rank that is crossed to reach the insert position
		if i == sl.level-1 {
			sl.rank[i] = 0
		} else {
			sl.rank[i] = sl.rank[i+1]
		}
		for x.level[i].forward != nil && x.level[i].forward.Value.Less(v) {
			sl.rank[i] += x.level[i].span
			x = x.level[i].forward
		}
		sl.update[i] = x
	}

	// ensure that the v is unique, the re-insertion of v should never happen since the
	// caller of sl.Insert() should test in the hash table if the element is already inside or not.
	level := randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			sl.rank[i] = 0
			sl.update[i] = sl.header
			sl.update[i].level[i].span = sl.length
		}
		sl.level = level
	}

	x = newElement(level, v)
	for i := 0; i < level; i++ {
		x.level[i].forward = sl.update[i].level[i].forward
		sl.update[i].level[i].forward = x

		// update span covered by update[i] as x is inserted here
		x.level[i].span = sl.update[i].level[i].span - sl.rank[0] + sl.rank[i]
		sl.update[i].level[i].span = sl.rank[0] - sl.rank[i] + 1
	}

	// increment span for untouched levels
	for i := level; i < sl.level; i++ {
		sl.update[i].level[i].span++
	}

	if sl.update[0] == sl.header {
		x.backward = nil
	} else {
		x.backward = sl.update[0]
	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		sl.tail = x
	}
	sl.length++

	return x
}

// deleteElement deletes e from its skiplist, and decrements sl.length.
func (sl *SkipList[T]) deleteElement(e *Element[T], update []*Element[T]) {
	for i := 0; i < sl.level; i++ {
		if update[i].level[i].forward == e {
			update[i].level[i].span += e.level[i].span - 1
			update[i].level[i].forward = e.level[i].forward
		} else {
			update[i].level[i].span -= 1
		}
	}

	if e.level[0].forward != nil {
		e.level[0].forward.backward = e.backward
	} else {
		sl.tail = e.backward
	}

	for sl.level > 1 && sl.header.level[sl.level-1].forward == nil {
		sl.level--
	}
	sl.length--
}

// Remove removes e from sl if e is an element of skiplist sl.
// It returns the element value e.Value.
func (sl *SkipList[T]) Remove(e *Element[T]) (res T) {
	x := sl.find(e.Value)                 // x.Value >= e.Value
	if x == e && !e.Value.Less(x.Value) { // e.Value >= x.Value
		sl.deleteElement(x, sl.update)
		res = x.Value
	}
	return
}

// Delete deletes an element e that e.Value == v, and returns e.Value or zero T value.
func (sl *SkipList[T]) Delete(v T) (res T) {
	x := sl.find(v)                   // x.Value >= v
	if x != nil && !v.Less(x.Value) { // v >= x.Value
		sl.deleteElement(x, sl.update)
		res = x.Value
	}
	return
}

// Find finds an element e that e.Value == v, and returns e or nil.
func (sl *SkipList[T]) Find(v T) *Element[T] {
	x := sl.find(v)                   // x.Value >= v
	if x != nil && !v.Less(x.Value) { // v >= x.Value
		return x
	}
	return nil
}

// find finds the first element e that e.Value >= v, and returns e or nil.
func (sl *SkipList[T]) find(v T) *Element[T] {
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.Value.Less(v) {
			x = x.level[i].forward
		}
		sl.update[i] = x
	}
	return x.level[0].forward
}

// GetRank finds the rank for an element e that e.Value == v,
// Returns 0 when the element cannot be found, rank otherwise.
// Note that the rank is 1-based due to the span of sl.header to the first element.
func (sl *SkipList[T]) GetRank(v T) int {
	x := sl.header
	rank := 0
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.Value.Less(v) {
			rank += x.level[i].span
			x = x.level[i].forward
		}
		if x.level[i].forward != nil && !x.level[i].forward.Value.Less(v) && !v.Less(x.level[i].forward.Value) {
			rank += x.level[i].span
			return rank
		}
	}
	return 0
}

// GetElementByRank finds an element by ites rank. The rank argument needs bo be 1-based.
// Note that is the first element e that GetRank(e.Value) == rank, and returns e or nil.
func (sl *SkipList[T]) GetElementByRank(rank int) *Element[T] {
	x := sl.header
	traversed := 0
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && traversed+x.level[i].span <= rank {
			traversed += x.level[i].span
			x = x.level[i].forward
		}
		if traversed == rank {
			return x
		}
	}
	return nil
}
