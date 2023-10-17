package skiplist

import (
	"math/rand"
)

const SKIPLIST_MAXLEVEL = 32
const SKIPLIST_BRANCH = 4

type skiplistLevel[T Interface[T]] struct {
	forward *Element[T]
	span    int
}

type Element[T Interface[T]] struct {
	Value    T
	backward *Element[T]
	level    []*skiplistLevel[T]
}

// Next returns the next skiplist element or nil.
func (e *Element[T]) Next() *Element[T] {
	return e.level[0].forward
}

// Prev returns the previous skiplist element of nil.
func (e *Element[T]) Prev() *Element[T] {
	return e.backward
}

// newElement returns an initialized element.
func newElement[T Interface[T]](level int, v T) *Element[T] {
	slLevels := make([]*skiplistLevel[T], level)
	for i := 0; i < level; i++ {
		slLevels[i] = new(skiplistLevel[T])
	}

	return &Element[T]{
		Value:    v,
		backward: nil,
		level:    slLevels,
	}
}

// randomLevel returns a random level.
func randomLevel() int {
	level := 1
	for (rand.Int31()&0xFFFF)%SKIPLIST_BRANCH == 0 {
		level += 1
	}

	if level < SKIPLIST_MAXLEVEL {
		return level
	} else {
		return SKIPLIST_MAXLEVEL
	}
}
