package errors

import (
	"errors"
)

// AsA finds the first error in err's tree that has the type E, and if one is found, returns that error value and true.
// Otherwise it returns the zero value of E and false.
// Note that this implement is refrence [#51945])(https://github.com/golang/go/issues/51945)
// and [#56949](https://github.com/golang/go/issues/56949).
func AsA[E error](err error) (_ E, ok bool) {
	var e *E
	for err != nil {
		if e, ok := err.(E); ok {
			return e, true
		}
		if x, ok := err.(interface{ As(any) bool }); ok {
			if e == nil {
				e = new(E)
			}
			if x.As(e) {
				return *e, true
			}
		}
		err = errors.Unwrap(err)
	}
	return
}
