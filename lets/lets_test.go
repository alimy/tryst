// Copyright 2023 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package lets

import (
	"testing"
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
		if res := If(tc.condition, tc.trueVal, tc.falseVal); res != tc.result {
			t.Errorf("If(%t, %d, %d) want %d but got %d", tc.condition, tc.trueVal, tc.falseVal, tc.result, res)
		}
	}
}
