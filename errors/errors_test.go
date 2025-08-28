package errors

import (
	"errors"
	"testing"
)

type myError struct{}

var sink *myError

func (me *myError) Error() string {
	return "my error"
}

func asA2[E error](err error) (E, bool) {
	var e E
	return e, errors.As(err, &e)
}

func BenchmarkAsA2(b *testing.B) {
	var (
		err error = &myError{}
		ok  bool
	)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if sink, ok = asA2[*myError](err); !ok {
			b.Fatal("AsA2 failed")
		}
	}
}

func BenchmarkAsA(b *testing.B) {
	var (
		err error = &myError{}
		ok  bool
	)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if sink, ok = AsA[*myError](err); !ok {
			b.Fatal("AsA failed")
		}
	}
}
