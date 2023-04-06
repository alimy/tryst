package tryst_test

import (
	"bytes"

	"github.com/alimy/tryst"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Range", Ordered, func() {
	type rangeBytes []struct {
		origin string
		fixed  string
	}
	var samples rangeBytes

	BeforeAll(func() {
		samples = rangeBytes{
			{
				origin: `a`,
				fixed:  `a`,
			},
			{
				origin: ``,
				fixed:  ``,
			},
			{
				origin: `ab`,
				fixed:  `ab`,
			},
			{
				origin: `abc`,
				fixed:  `abc`,
			},
			{
				origin: `abcdef`,
				fixed:  `abcdef`,
			},
			{
				origin: `abcdefg`,
				fixed:  `abcdefg`,
			},
		}
	})

	It("Range rune by step 2", func() {
		for _, t := range samples {
			buf := bytes.NewBuffer(make([]byte, 0, len(t.origin)))
			src := make([]rune, 0, len(t.origin))
			for _, c := range t.origin {
				src = append(src, c)
			}
			tryst.Range(src, 2, func(slice ...rune) error {
				for _, c := range slice {
					buf.WriteRune(c)
				}
				return nil
			})
			Expect(buf.String()).To(Equal(t.fixed))
		}
	})

	It("Range rune by step -2", func() {
		for _, t := range samples {
			buf := bytes.NewBuffer(make([]byte, 0, len(t.origin)))
			src := make([]rune, 0, len(t.origin))
			for _, c := range t.origin {
				src = append(src, c)
			}
			tryst.Range(src, -2, func(slice ...rune) error {
				for i := len(slice) - 1; i >= 0; i-- {
					buf.WriteRune(slice[i])
				}
				return nil
			})
			other := bytes.NewBuffer(buf.Bytes())
			for i := len(src) - 1; i >= 0; i-- {
				src[i], _, _ = other.ReadRune()
			}
			buf = &bytes.Buffer{}
			for _, c := range src {
				buf.WriteRune(c)
			}
			Expect(buf.String()).To(Equal(t.fixed))
		}
	})
})
