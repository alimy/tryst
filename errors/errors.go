package errors

// AsA finds the first error in err's tree that has the type E, and if one is found, returns that error value and true.
// Otherwise it returns the zero value of E and false.
// Note that this implement is refrence [#51945])(https://github.com/golang/go/issues/51945)
// and [#56949](https://github.com/golang/go/issues/56949).
func AsA[E error](err error) (e E, ok bool) {
	var r *E
	for err != nil {
		if e, ok = err.(E); ok {
			return
		}
		if x, ok := err.(interface{ As(any) bool }); ok {
			if r == nil {
				r = new(E)
			}
			if x.As(r) {
				return *r, true
			}
		}
		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap()
		case interface{ Unwrap() []error }:
			for _, err := range x.Unwrap() {
				if err == nil {
					continue
				}
				if e, ok = AsA[E](err); ok {
					return
				}
			}
			return
		default:
			return
		}
	}
	return
}
