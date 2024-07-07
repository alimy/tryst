// Copyright 2024 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package cyclist

import (
	"github.com/alimy/tryst/lets"
)

// Cyclist[T] A Cyclist is an element of a circular list like container/ring in standard library but based on slice.
type Cyclist[T any] struct {
	slice    []T
	capacity int
	size     int
	begin    int
	end      int
}

// Prev returns the previous n prev element.
func (l *Cyclist[T]) Prev(n int) []T {
	res := make([]T, 0, n)
	idx, size := l.end, l.size
	n %= (l.size + 1)
	for i := 0; i < n; i++ {
		res = append(res, l.slice[idx])
		idx--
		if idx < 0 {
			idx = l.capacity - 1
		}
		size--
	}
	return res
}

// Next returns the next n cyclist element.
func (l *Cyclist[T]) Next(n int) []T {
	res := make([]T, 0, n)
	idx, size := l.begin, l.size
	n %= (l.size + 1)
	for i := 0; i < n; i++ {
		res = append(res, l.slice[idx])
		idx++
		idx %= l.capacity
		size--
	}
	return res
}

func (l *Cyclist[T]) Put(s ...T) {
	for _, v := range s {
		l.slice[l.end] = v
		l.end++
		l.end %= l.capacity
		if l.size != l.capacity {
			l.size++
		} else {
			l.begin++
			l.begin %= l.begin
		}
	}
}

// Move moves n % l.Len() elements backward (n < 0) or forward (n >= 0) in the cyclist and returns that ring element.
func (l *Cyclist[T]) Move(n int) (res []T) {
	res = make([]T, 0, n)
	if n > 0 {
		n %= (l.size + 1)
		for i := 0; i < n; i++ {
			res = append(res, l.slice[l.begin])
			l.begin++
			l.begin %= l.capacity
			l.size--
		}
	} else if n < 0 {
		n %= (l.size + 1)
		for i := n; i < 0; i++ {
			res = append(res, l.slice[l.end])
			l.end--
			if l.end < 0 {
				l.end = l.capacity - 1
			}
			l.size--
		}
	}
	return
}

// Do calls function f on each element of the cyclist, in forward order.
func (l *Cyclist[T]) Do(f func(T)) {
	idx := l.begin
	for i := 0; i < l.size; i++ {
		f(l.slice[idx])
		idx++
		idx %= l.capacity
	}
}

// As convert cyclist to slice. In reverse mode if reverse is true.
func (l *Cyclist[T]) As(reverse ...bool) (res []T) {
	if inReverse := lets.Val(false, reverse...); inReverse {
		res = l.Prev(l.size)
	} else {
		res = l.Next(l.size)
	}
	return
}

// Len computes the number of elements in cyclist l. It executes in time proportional to the number of elements.
func (l *Cyclist[T]) Len() int {
	return l.size
}

// Capacity computes the capacity of cyclist l.
func (l *Cyclist[T]) Capacity() int {
	return l.capacity
}

// New[T]  creates a cyclist of n(n>0) elements.
func New[T any](n int) *Cyclist[T] {
	if n <= 0 {
		n = 1
	}
	return &Cyclist[T]{
		slice:    make([]T, n),
		capacity: n,
	}
}
