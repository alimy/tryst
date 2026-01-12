// Copyright 2024 Michael Li <alimy@niubiu.com>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package utils

import "unsafe"

// String convert bytes to string
func String(data []byte) (res string) {
	if size := len(data); size > 0 {
		res = unsafe.String(unsafe.SliceData(data), size)
	}
	return
}

// Bytes convert string to []byte
func Bytes(data string) (res []byte) {
	if size := len(data); size > 0 {
		res = unsafe.Slice(unsafe.StringData(data), size)
	} else {
		res = []byte{}
	}
	return
}

// EqualBytes compare b1 == b2
func EqualBytes(b1, b2 []byte) bool {
	b1Size, b2Size := len(b1), len(b2)
	if b1Size != b2Size {
		return false
	}
	return unsafe.String(unsafe.SliceData(b1), b1Size) == unsafe.String(unsafe.SliceData(b2), b2Size)
}

// CompareBytes compare b1 and b2, 1 if b1>b2, 0 if b1==b2, -1 if b1<b2
func CompareBytes(b1, b2 []byte) (res int) {
	s1, s2 := unsafe.String(unsafe.SliceData(b1), len(b1)), unsafe.String(unsafe.SliceData(b2), len(b2))
	switch {
	case s1 > s2:
		res = 1
	case s1 < s2:
		res = -1
	default:
		res = 0
	}
	return
}
