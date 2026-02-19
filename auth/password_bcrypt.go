// Copyright 2024 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package auth

import (
	"golang.org/x/crypto/bcrypt"
)

var _ PasswordProvider = (*bcryptPasswordProvider)(nil)

type bcryptPasswordProvider struct {
	cost int
}

func (p *bcryptPasswordProvider) Generate(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, p.cost)
}

func (p *bcryptPasswordProvider) Compare(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func NewBcryptPasswordProvider(cost int) *bcryptPasswordProvider {
	return &bcryptPasswordProvider{
		cost: cost,
	}
}
