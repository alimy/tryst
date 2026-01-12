// Copyright 2024 Michael Li <alimy@niubiu.com>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package auth

import (
	"crypto"
	"errors"
	"hash"

	"github.com/alimy/tryst/utils"
)

var _ HashPasswordProvider = (*hashPasswordProvider)(nil)

type hashPasswordProvider struct {
	hashFactor func() hash.Hash
}

func (p *hashPasswordProvider) Generate(password []byte, salt []byte) ([]byte, error) {
	hashFn := p.hashFactor()
	hashFn.Write(password)
	return hashFn.Sum(salt), nil
}

func (p *hashPasswordProvider) Compare(hashedPassword, password []byte, salt []byte) error {
	hashFn := p.hashFactor()
	hashFn.Write(password)
	if !utils.EqualBytes(hashedPassword, hashFn.Sum(salt)) {
		return errors.New("incorrect password")
	}
	return nil
}

func NewHashPasswordProvider(hashFactor func() hash.Hash) *hashPasswordProvider {
	return &hashPasswordProvider{
		hashFactor: hashFactor,
	}
}

func NewSha1PasswordProvider() *hashPasswordProvider {
	return &hashPasswordProvider{
		hashFactor: crypto.SHA1.New,
	}
}
