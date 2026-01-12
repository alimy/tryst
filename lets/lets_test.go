// Copyright 2023 Michael Li <alimy@niubiu.com>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package lets_test

import (
	"testing"

	"github.com/alimy/tryst/lets"
)

func TestIf(t *testing.T) {
	for _, tc := range []struct {
		condition bool
		trueVal   int
		falseVal  int
		result    int
	}{
		{true, 1, 2, 1},
		{false, 1, 2, 2},
	} {
		if res := lets.If(tc.condition, tc.trueVal, tc.falseVal); res != tc.result {
			t.Errorf("If(%t, %d, %d) want %d but got %d", tc.condition, tc.trueVal, tc.falseVal, tc.result, res)
		}
	}
}

func TestVal(t *testing.T) {
	for _, tc := range []struct {
		v int
		s []int
		r int
	}{
		{4, []int{5, 6}, 5},
		{4, []int{5}, 5},
		{5, []int{}, 5},
	} {
		if res := lets.Val(tc.v, tc.s...); res != tc.r {
			t.Errorf("give v:%d s:%+v want: %d but got %d", tc.v, tc.s, tc.r, res)
		}
	}
}
