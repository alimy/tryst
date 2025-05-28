// Copyright 2024 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

//go:build !jsoniter

package json

import "unsafe"

func MarshalToString(v any) (res string, err error) {
	var out []byte
	out, err = Marshal(v)
	if size := len(out); size > 0 {
		res = unsafe.String(unsafe.SliceData(out), size)
	}
	return
}

func UnmarshalFromString(data string, v any) error {
	var contents []byte
	if size := len(data); size > 0 {
		contents = unsafe.Slice(unsafe.StringData(data), size)
	}
	return Unmarshal(contents, v)
}
