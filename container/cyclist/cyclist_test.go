// Copyright 2024 Michael Li <alimy@niubiu.com>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package cyclist_test

import (
	"testing"

	"github.com/alimy/tryst/container/cyclist"
)

func TestNew(t *testing.T) {
	for _, d := range []struct {
		input  int
		expect int
	}{
		{2, 2},
		{3, 3},
		{1, 1},
		{0, 1},
		{-1, 1},
		{-2, 1},
	} {
		l := cyclist.New[int](d.input)
		if l.Len() != 0 || l.Capacity() != d.expect {
			t.Errorf("create instance by New(%d) expect len 0 and cacpcity %d but got %d/%d", d.input, d.expect, l.Len(), l.Capacity())
		}
	}
}

func TestAs(t *testing.T) {
	for _, d := range []struct {
		input   []int
		size    int
		reverse bool
		expect  []int
	}{
		{
			input:   []int{},
			size:    6,
			reverse: false,
			expect:  []int{},
		},
		{
			input:   []int{1},
			size:    6,
			reverse: true,
			expect:  []int{1},
		},
		{
			input:   []int{1},
			size:    6,
			reverse: true,
			expect:  []int{1},
		},
		{
			input:   []int{1, 2, 3, 4, 5, 6},
			size:    6,
			reverse: false,
			expect:  []int{1, 2, 3, 4, 5, 6},
		},
		{
			input:   []int{1, 2, 3, 4, 5, 6},
			size:    5,
			reverse: true,
			expect:  []int{6, 5, 4, 3, 2},
		},
		{
			input:   []int{1, 2, 3, 4, 5, 6},
			size:    5,
			reverse: false,
			expect:  []int{2, 3, 4, 5, 6},
		},
		{
			input:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			size:    5,
			reverse: false,
			expect:  []int{11, 12, 13, 14, 15},
		},
		{
			input:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			size:    5,
			reverse: true,
			expect:  []int{15, 14, 13, 12, 11},
		},
		{
			input:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			size:    4,
			reverse: false,
			expect:  []int{13, 14, 15, 16},
		},
		{
			input:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			size:    4,
			reverse: true,
			expect:  []int{16, 15, 14, 13},
		},
	} {
		l := cyclist.New[int](d.size)
		l.Put(d.input...)
		val := []int{}
		l.Do(func(i int) {
			val = append(val, i)
		})
		res := l.As(d.reverse)
		if !eqSlice(res, d.expect) {
			t.Errorf("cyclist(%d){%v} As(%t) expect %v but got %v", d.size, d.input, d.reverse, d.expect, res)
		}
		if !d.reverse {
			res = l.As()
			if !eqSlice(res, d.expect) {
				t.Errorf("cyclist(%d){%v} As(%t) expect %v but got %v", d.size, d.input, d.reverse, d.expect, res)
			}
		}
	}
}

func TestPrev(t *testing.T) {
	for _, d := range []struct {
		input  []int
		size   int
		n      int
		expect []int
	}{
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   6,
			n:      3,
			expect: []int{6, 5, 4},
		},
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   5,
			n:      3,
			expect: []int{6, 5, 4},
		},
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   5,
			n:      2,
			expect: []int{6, 5},
		},
		{
			input:  []int{1, 2, 3, 4, 5, 6, 7, 8},
			size:   4,
			n:      2,
			expect: []int{8, 7},
		},
	} {
		l := cyclist.New[int](d.size)
		l.Put(d.input...)
		res := l.Prev(d.n)
		if !eqSlice(res, d.expect) {
			t.Errorf("cyclist(%d){%v} Prev(%d) expect %v but got %v", d.size, d.input, d.n, d.expect, res)
		}
	}
}

func TestMove(t *testing.T) {
	for _, d := range []struct {
		input  []int
		size   int
		n      int
		expect []int
		val    []int
	}{
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   6,
			n:      3,
			expect: []int{1, 2, 3},
			val:    []int{4, 5, 6},
		},
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   5,
			n:      3,
			expect: []int{2, 3, 4},
			val:    []int{5, 6},
		},
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   3,
			n:      2,
			expect: []int{4, 5},
			val:    []int{6},
		},
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   4,
			n:      2,
			expect: []int{3, 4},
			val:    []int{5, 6},
		},
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   4,
			n:      -2,
			expect: []int{6, 5},
			val:    []int{3, 4},
		},
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   3,
			n:      2,
			expect: []int{4, 5},
			val:    []int{6},
		},
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   3,
			n:      -2,
			expect: []int{6, 5},
			val:    []int{4},
		},
	} {
		l := cyclist.New[int](d.size)
		l.Put(d.input...)
		res := l.Move(d.n)
		var val []int
		l.Do(func(i int) {
			val = append(val, i)
		})
		if !eqSlice(res, d.expect) || !eqSlice(val, d.val) {
			t.Errorf("cyclist(%d){%v} Move(%d) expect %v val %v but got %v val %v", d.size, d.input, d.n, d.expect, d.val, res, val)
		}
	}
}

func TestNext(t *testing.T) {
	for _, d := range []struct {
		input  []int
		size   int
		n      int
		expect []int
	}{
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   6,
			n:      3,
			expect: []int{1, 2, 3},
		},
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   5,
			n:      3,
			expect: []int{2, 3, 4},
		},
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   3,
			n:      2,
			expect: []int{4, 5},
		},
	} {
		l := cyclist.New[int](d.size)
		l.Put(d.input...)
		res := l.Next(d.n)
		if !eqSlice(res, d.expect) {
			t.Errorf("cyclist(%d){%v} Next(%d) expect %v but got %v", d.size, d.input, d.n, d.expect, res)
		}
	}
}

func TestDo(t *testing.T) {
	for _, d := range []struct {
		input  []int
		size   int
		expect []int
	}{
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   6,
			expect: []int{1, 2, 3, 4, 5, 6},
		},
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   5,
			expect: []int{2, 3, 4, 5, 6},
		},
		{
			input:  []int{1, 2, 3, 4, 5, 6},
			size:   3,
			expect: []int{4, 5, 6},
		},
	} {
		l := cyclist.New[int](d.size)
		l.Put(d.input...)
		var val []int
		l.Do(func(i int) {
			val = append(val, i)
		})
		if !eqSlice(val, d.expect) {
			t.Errorf("cyclist(%d){%v} expect %v but got %v", d.size, d.input, d.expect, val)
		}
	}
}

func eqSlice(s1, s2 []int) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
