// Copyright 2023 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package pool

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkPool(t *testing.T) {
	p := NewWorkPool(10, 100*time.Millisecond)
	var (
		sum  int32
		wg   sync.WaitGroup
		size int32 = 100
	)
	for i := int32(0); i < size; i++ {
		wg.Add(1)
		p.Go(func() {
			defer wg.Done()
			atomic.AddInt32(&sum, 1)
		})
	}
	wg.Wait()
	if v := atomic.LoadInt32(&sum); v != size {
		t.Errorf("want sum equel %d but got %d", size, v)
	}
}

func TestBufferWorkPool(t *testing.T) {
	p := NewBufferWorkPool(-1, 100*time.Millisecond)
	var (
		sum  int32
		wg   sync.WaitGroup
		size int32 = 100
	)
	for i := int32(0); i < size; i++ {
		wg.Add(1)
		p.Go(func() {
			defer wg.Done()
			atomic.AddInt32(&sum, 1)
		})
	}
	wg.Wait()
	if v := atomic.LoadInt32(&sum); v != size {
		t.Errorf("want sum equel %d but got %d", size, v)
	}
}

func TestSimpleWorkPool(t *testing.T) {
	p := NewSimpleWorkPool(10, 100*time.Millisecond)
	var (
		sum  int32
		wg   sync.WaitGroup
		size int32 = 100
	)
	for i := int32(0); i < size; i++ {
		wg.Add(1)
		p.Go(func() {
			defer wg.Done()
			atomic.AddInt32(&sum, 1)
		})
	}
	wg.Wait()
	if v := atomic.LoadInt32(&sum); v != size {
		t.Errorf("want sum equel %d but got %d", size, v)
	}
}
