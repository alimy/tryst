// Copyright 2023 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package tryst

// Range traverse source by step to handle items.
// Forward range if step > 0 and backward range if step < 0.
// Just handle all source items if step == 0.
func Range[T any](source []T, step int, fn func(slice ...T) error) error {
	switch {
	case step > 0:
		return rangeForward(source, step, fn)
	case step < 0:
		return rangeBackward(source, -step, fn)
	default:
		return fn(source...)
	}
}

func rangeForward[T any](source []T, step int, fn func(slice ...T) error) (err error) {
	last := len(source) - step
	end := 0
	for i := 0; i < last; i += step {
		end += step
		if err = fn(source[i:end]...); err != nil {
			return
		}
	}
	err = fn(source[end:]...)
	return
}

func rangeBackward[T any](source []T, step int, fn func(slice ...T) error) (err error) {
	end := len(source)
	for i := end; i >= step; i -= step {
		if err = fn(source[i-step : end]...); err != nil {
			return
		}
		end -= step
	}
	err = fn(source[:end]...)
	return
}
