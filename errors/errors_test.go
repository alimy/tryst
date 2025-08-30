package errors

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"
)

type myError struct{}

type mirError struct{}

type hiError struct{}

type wrapError struct {
	innerErr error
	msg      string
}

var sink *myError

func (*myError) Error() string {
	return "my error"
}

func (*mirError) Error() string {
	return "mir error"
}

func (hiError) Error() string {
	return "hi error"
}

func (e *wrapError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.innerErr)
}

func (e *wrapError) Unwrap() error {
	return e.innerErr
}

func wrapErrors(msg string, errs ...error) error {
	var innerErr error
	if len(errs) == 1 {
		innerErr = errs[0]
	} else {
		innerErr = errors.Join(errs...)
	}
	return &wrapError{
		msg:      msg,
		innerErr: innerErr,
	}
}

func testError() error {
	return wrapErrors("wrap error", &myError{}, &mirError{}, hiError{})
}

func asA2[E error](err error) (E, bool) {
	var e E
	return e, errors.As(err, &e)
}

type poser struct {
	msg string
	f   func(error) bool
}

var poserPathErr = &fs.PathError{Op: "poser"}

func (p *poser) Error() string     { return p.msg }
func (p *poser) Is(err error) bool { return p.f(err) }
func (p *poser) As(err any) bool {
	switch x := err.(type) {
	case **poser:
		*x = p
	case *errorT:
		*x = errorT{"poser"}
	case **fs.PathError:
		*x = poserPathErr
	default:
		return false
	}
	return true
}

type errorT struct{ s string }

func (e errorT) Error() string { return fmt.Sprintf("errorT(%s)", e.s) }

type wrapped struct {
	msg string
	err error
}

func (e wrapped) Error() string { return e.msg }
func (e wrapped) Unwrap() error { return e.err }

type multiErr []error

func (m multiErr) Error() string   { return "multiError" }
func (m multiErr) Unwrap() []error { return []error(m) }

func TestAsA(t *testing.T) {
	var errT errorT
	var errP *fs.PathError
	var p *poser
	_, errF := os.Open("non-existing")
	poserErr := &poser{"oh no", nil}

	testCases := []struct {
		err    error
		target any
		match  bool
		want   any // value of target on match
	}{{
		nil,
		errP,
		false,
		nil,
	}, {
		wrapped{"pitied the fool", errorT{"T"}},
		errT,
		true,
		errorT{"T"},
	}, {
		errF,
		errP,
		true,
		errF,
	}, {
		errorT{},
		errP,
		false,
		nil,
	}, {
		wrapped{"wrapped", nil},
		errT,
		false,
		nil,
	}, {
		&poser{"error", nil},
		errT,
		true,
		errorT{"poser"},
	}, {
		&poser{"path", nil},
		errP,
		true,
		poserPathErr,
	}, {
		poserErr,
		p,
		true,
		poserErr,
	}, {
		multiErr{},
		errT,
		false,
		nil,
	}, {
		multiErr{errors.New("a"), errorT{"T"}},
		errT,
		true,
		errorT{"T"},
	}, {
		multiErr{errorT{"T"}, errors.New("a")},
		errT,
		true,
		errorT{"T"},
	}, {
		multiErr{errorT{"a"}, errorT{"b"}},
		errT,
		true,
		errorT{"a"},
	}, {
		multiErr{multiErr{errors.New("a"), errorT{"a"}}, errorT{"b"}},
		errT,
		true,
		errorT{"a"},
	}, {
		multiErr{nil},
		errT,
		false,
		nil,
	}}
	for i, tc := range testCases {
		name := fmt.Sprintf("%d:As(Errorf(..., %v), %v)", i, tc.err, tc.target)
		t.Run(name, func(t *testing.T) {
			var (
				got   any
				match bool
			)
			switch tc.target.(type) {
			case errorT:
				got, match = AsA[errorT](tc.err)
			case *fs.PathError:
				got, match = AsA[*fs.PathError](tc.err)
			case *poser:
				got, match = AsA[*poser](tc.err)
			}
			if match != tc.match {
				t.Fatalf("match: got %v; want %v", match, tc.match)
			}
			if !match {
				return
			}
			if got != tc.want {
				t.Fatalf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

func BenchmarkAsA2(b *testing.B) {
	err, ok := testError(), false

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if sink, ok = asA2[*myError](err); !ok {
			b.Fatal("AsA2 failed")
		}
	}
}

func BenchmarkAsA(b *testing.B) {
	err, ok := testError(), false

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if sink, ok = AsA[*myError](err); !ok {
			b.Fatal("AsA failed")
		}
	}
}
