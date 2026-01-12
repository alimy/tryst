// Copyright 2023 Michael Li <alimy@niubiu.com>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

// package lets contain some help function for go develop

package lets

// If[T] simulate ternary conditional operator (condition ? trueVal : falseVal)
func If[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

// Val[T] return s[0] if s is give or else return v
func Val[T any](v T, s ...T) T {
	if len(s) > 0 {
		return s[0]
	}
	return v
}
