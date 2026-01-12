// Copyright 2024 Michael Li <alimy@niubiu.com>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package auth

type PasswordProvider interface {
	Generate(password []byte) ([]byte, error)
	Compare(hashedPassword, password []byte) error
}

type HashPasswordProvider interface {
	Generate(password []byte, salt []byte) ([]byte, error)
	Compare(hashedPassword, password []byte, salt []byte) error
}
