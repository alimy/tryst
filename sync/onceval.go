// Copyright 2025 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package sync

import (
	"sync"
)

type (
	ValFn[T any]       func() T
	ValsFn[T1, T2 any] func() (T1, T2)
)

// OnceValFn returns a function that invokes f only once and returns the value
// returned by f. The returned function may be called concurrently.
func OnceValFn[T any](newFn func() T) ValFn[T] {
	res := newOnceSys(newFn)
	return res.Val
}

// OnceValsFn returns a function that invokes f only once and returns the value
// returned by f. The returned function may be called concurrently.
func OnceValsFn[T1, T2 any](newFn func() (T1, T2)) ValsFn[T1, T2] {
	res := newOnceSys2(newFn)
	return res.Val
}

// OnceVal returns a function that invokes f only once and returns the value
// returned by f. The returned function may be called concurrently.
//
// If f panics, the returned function will panic with the same value on every call.
func OnceVal[T any](newFn func() T) ValFn[T] {
	return sync.OnceValue(newFn)
}

// OnceVals returns a function that invokes f only once and returns the values
// returned by f. The returned function may be called concurrently.
//
// If f panics, the returned function will panic with the same value on every call.
func OnceVals[T1, T2 any](newFn func() (T1, T2)) ValsFn[T1, T2] {
	return sync.OnceValues(newFn)
}
