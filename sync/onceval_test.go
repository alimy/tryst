// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"sync"
	"testing"
	_ "unsafe"

	isync "github.com/alimy/tryst/sync"
)

func TestOnceVal(t *testing.T) {
	calls := 0
	f := isync.OnceVal(func() int {
		calls++
		return calls
	})
	allocs := testing.AllocsPerRun(10, func() { f() })
	value := f()
	if calls != 1 {
		t.Errorf("want calls==1, got %d", calls)
	}
	if value != 1 {
		t.Errorf("want value==1, got %d", value)
	}
	if allocs != 0 {
		t.Errorf("want 0 allocations per call, got %v", allocs)
	}
}

func TestOnceVals(t *testing.T) {
	calls := 0
	f := isync.OnceValsFn(func() (int, int) {
		calls++
		return calls, calls + 1
	})
	allocs := testing.AllocsPerRun(10, func() { f() })
	v1, v2 := f()
	if calls != 1 {
		t.Errorf("want calls==1, got %d", calls)
	}
	if v1 != 1 || v2 != 2 {
		t.Errorf("want v1==1 and v2==2, got %d and %d", v1, v2)
	}
	if allocs != 0 {
		t.Errorf("want 0 allocations per call, got %v", allocs)
	}
}

var (
	onceVal = isync.OnceValFn(func() int { return 42 })

	onceValOnce  sync.Once
	onceValValue int
)

func doOnceVal() int {
	onceValueOnce.Do(func() {
		onceValueValue = 42
	})
	return onceValueValue
}

func BenchmarkOnceValFn(b *testing.B) {
	// See BenchmarkOnceFunc
	b.Run("v=Once", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if want, got := 42, doOnceValue(); want != got {
				b.Fatalf("want %d, got %d", want, got)
			}
		}
	})
	b.Run("v=Global", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if want, got := 42, onceValue(); want != got {
				b.Fatalf("want %d, got %d", want, got)
			}
		}
	})
	b.Run("v=Local", func(b *testing.B) {
		b.ReportAllocs()
		onceValue := isync.OnceValFn(func() int { return 42 })
		for i := 0; i < b.N; i++ {
			if want, got := 42, onceValue(); want != got {
				b.Fatalf("want %d, got %d", want, got)
			}
		}
	})
}
