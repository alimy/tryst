// Copyright 2025 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package sync

import (
	"sync"
)

type onceSys[T any] struct {
	once   *sync.Once
	object T
	newFn  func() T
}

type onceSys2[T1, T2 any] struct {
	once    *sync.Once
	object1 T1
	object2 T2
	newFn   func() (T1, T2)
}

func (s *onceSys[T]) Val() T {
	s.once.Do(s.initial)
	return s.object
}

func (s *onceSys[T]) initial() {
	s.object = s.newFn()
}

func (s *onceSys2[T1, T2]) Val() (T1, T2) {
	s.once.Do(s.initial)
	return s.object1, s.object2
}

func (s *onceSys2[T1, T2]) initial() {
	s.object1, s.object2 = s.newFn()
}

func newOnceSys[T any](newFn func() T) *onceSys[T] {
	return &onceSys[T]{
		once:  &sync.Once{},
		newFn: newFn,
	}
}

func newOnceSys2[T1, T2 any](newFn func() (T1, T2)) *onceSys2[T1, T2] {
	return &onceSys2[T1, T2]{
		once:  &sync.Once{},
		newFn: newFn,
	}
}
