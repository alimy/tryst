package tryst_test

import (
	"bytes"

	"github.com/alimy/tryst"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scan", Ordered, func() {
	type rangeBytes []struct {
		origin string
		fixed  string
	}
	var samples rangeBytes

	BeforeAll(func() {
		samples = rangeBytes{
			{
				origin: `SELECT * FROM @user WHERE username=?@_`,
				fixed:  `SELECT * FROM user WHERE username=?_`,
			},
			{
				origin: `SELECT * FROM @user WHERE username=?`,
				fixed:  `SELECT * FROM user WHERE username=?`,
			},
			{
				origin: `SELECT * FROM @@user WHERE 用户名=?`,
				fixed:  `SELECT * FROM @@user WHERE 用户名=?`,
			},
			{
				origin: `SELECT * FROM @@user, @@@@contact WHERE 用户名=?`,
				fixed:  `SELECT * FROM @@user, @@@@contact WHERE 用户名=?`,
			},
			{
				origin: `SELECT @@name, @@@@id FROM @@user, @@@@contact WHERE 用户名=?`,
				fixed:  `SELECT @@name, @@@@id FROM @@user, @@@@contact WHERE 用户名=?`,
			},
			{
				origin: `SELECT @name, @id FROM @user, @contact WHERE 用户名=?`,
				fixed:  `SELECT name, id FROM user, contact WHERE 用户名=?`,
			},
		}
	})

	It("Scan rune by step 2", func() {
		for _, t := range samples {
			buf := bytes.NewBuffer(make([]byte, 0, len(t.origin)))
			src := make([]rune, 0, len(t.origin))
			for _, c := range t.origin {
				src = append(src, c)
			}
			isPrevAt := false
			tryst.Scan(src, 2, func(slice ...rune) error {
				if len(slice) == 2 {
					if slice[0] == '@' && isPrevAt {
						buf.WriteRune('@')
						isPrevAt = false
					} else if slice[0] == '@' && !isPrevAt {
						if slice[1] == '@' {
							buf.WriteRune('@')
						}
						isPrevAt = true
					} else {
						buf.WriteRune(slice[0])
						isPrevAt = false
					}
				} else if len(slice) == 1 {
					if slice[0] == '@' && isPrevAt {
						buf.WriteRune('@')
					}
					buf.WriteRune(slice[0])
				}
				return nil
			})
			Expect(buf.String()).To(Equal(t.fixed))
		}
	})

	It("Scan rune by step -2", func() {
		for _, t := range samples {
			buf := bytes.NewBuffer(make([]byte, 0, len(t.origin)))
			src := make([]rune, 0, len(t.origin))
			for _, c := range t.origin {
				src = append(src, c)
			}
			// size := len(src)
			// for i := 0; i < size/2; i++ {
			// 	src[i], src[size-i-1] = src[size-i-1], src[i]
			// }
			isPrevAt := false
			tryst.Scan(src, -2, func(slice ...rune) error {
				if len(slice) == 2 {
					slice[0], slice[1] = slice[1], slice[0]
					if slice[0] == '@' && isPrevAt {
						buf.WriteRune('@')
						isPrevAt = false
					} else if slice[0] == '@' && !isPrevAt {
						if slice[1] == '@' {
							buf.WriteRune('@')
						}
						isPrevAt = true
					} else {
						buf.WriteRune(slice[0])
						isPrevAt = false
					}
				} else if len(slice) == 1 {
					if slice[0] == '@' && isPrevAt {
						buf.WriteRune('@')
					}
					buf.WriteRune(slice[0])
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
