// Copyright 2024 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package slices

import (
	"github.com/RoaringBitmap/roaring"
)

type keyable interface {
	Key() int
}

// Distinct return distinc elements of S
func Distinct[S ~[]E, E keyable](s S) S {
	return DistinctFunc(s, func(e E) int {
		return e.Key()
	})
}

// DistinctFunc return distinc elements of S
func DistinctFunc[S ~[]E, E any](s S, key func(e E) int) (res S) {
	res = make(S, 0, len(s))
	x, bitmap := 0, roaring.New()
	for _, e := range s {
		x = key(e)
		if !bitmap.ContainsInt(x) {
			res = append(res, e)
			bitmap.AddInt(x)
		}
	}
	return
}
