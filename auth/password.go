// Copyright 2024 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordProvider interface {
	Generate(password []byte) ([]byte, error)
	Compare(hashedPassword, password []byte) error
}

func NewBcryptPasswordProvider(cost int) PasswordProvider {
	return &bcryptPasswordProvider{
		cost: cost,
	}
}

type bcryptPasswordProvider struct {
	cost int
}

func (p *bcryptPasswordProvider) Generate(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, p.cost)
}

func (p *bcryptPasswordProvider) Compare(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
