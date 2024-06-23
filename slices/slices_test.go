// Copyright 2024 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package slices_test

import (
	"cmp"
	stdSlices "slices"
	"testing"

	"github.com/alimy/tryst/slices"
)

type aliasInt int

func (e aliasInt) Key() int {
	return int(e)
}

func TestDistinct(t *testing.T) {
	for _, d := range []struct {
		input  []aliasInt
		expect []aliasInt
	}{
		{
			input:  []aliasInt{1, 2, 3, 1, 2, 3, 4, 5, 2, 6},
			expect: []aliasInt{1, 2, 3, 4, 5, 6},
		},
		{
			input:  []aliasInt{0, 2, 3, 1, 2, 3, 4, 5, 2, 6, 1, 1, 3, 5, 6, 0, 1, 2, 3},
			expect: []aliasInt{0, 1, 2, 3, 4, 5, 6},
		},
	} {
		res := slices.Distinct(d.input)
		if !eqAliasIntSlice(res, d.expect) {
			t.Errorf("input:%v want:%v but got:%v", d.input, d.expect, res)
		}
	}
}

func TestDistinctFunc(t *testing.T) {
	for _, d := range []struct {
		input  []int
		expect []int
	}{
		{
			input:  []int{1, 2, 3, 1, 2, 3, 4, 5, 2, 6},
			expect: []int{1, 2, 3, 4, 5, 6},
		},
		{
			input:  []int{0, 2, 3, 1, 2, 3, 4, 5, 2, 6, 1, 1, 3, 5, 6, 0, 1, 2, 3},
			expect: []int{0, 1, 2, 3, 4, 5, 6},
		},
	} {
		res := slices.DistinctFunc(d.input, func(e int) int {
			return e
		})
		if !eqSlice(res, d.expect) {
			t.Errorf("input:%v want:%v but got:%v", d.input, d.expect, res)
		}
	}
}

func TestDistinctFn(t *testing.T) {
	for _, d := range []struct {
		input  []string
		expect []string
	}{
		{
			input:  []string{"1", "2", "3", "1", "2", "3", "4", "5", "2", "6"},
			expect: []string{"1", "2", "3", "4", "5", "6"},
		},
		{
			input:  []string{"abc", "bcd", "ccd", "abc", "ccd", "bcd", "a", "a", "b", "c"},
			expect: []string{"abc", "bcd", "ccd", "a", "b", "c"},
		},
	} {
		res := slices.DistinctFn(d.input, func(e string) string {
			return e
		})
		if !eqSlice(res, d.expect) {
			t.Errorf("input:%v want:%v but got:%v", d.input, d.expect, res)
		}
	}
}

func eqSlice[E cmp.Ordered](s1, s2 []E) bool {
	if len(s1) != len(s2) {
		return false
	}
	stdSlices.Sort(s1)
	stdSlices.Sort(s2)
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func eqAliasIntSlice(s1, s2 []aliasInt) bool {
	t1, t2 := toIntSlice(s1), toIntSlice(s2)
	return eqSlice(t1, t2)
}

func toIntSlice(s []aliasInt) []int {
	res := make([]int, 0, len(s))
	for _, e := range s {
		res = append(res, int(e))
	}
	return res
}
