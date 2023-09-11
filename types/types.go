// Copyright 2023 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package types

// Empty empty alias type
type Empty = struct{}

// Set[T] a simple set alias to map
type Set[T comparable] map[T]struct{}

// Fn empty argument func alias type
type Fn = func()
