// Copyright 2024 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package slices

import (
	"hash/maphash"
)

// DistinctFn return distinc elements of S
func DistinctFn[S ~[]E, E any](s S, key func(e E) string) (res S) {
	seed := maphash.MakeSeed()
	return DistinctFunc(s, func(e E) int {
		return int(maphash.String(seed, key(e)))
	})
}
