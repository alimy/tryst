// Copyright 2023 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package cache

import (
	lru "github.com/hashicorp/golang-lru/v2"
)

// KeyPool[K] key pool used for cache keys
type KeyPool[K comparable] interface {
	Get(key K) string
}

type lruKeyPool[K comparable] struct {
	pool  *lru.Cache[K, string]
	newFn func(K) string
}

func (p *lruKeyPool[K]) Get(key K) string {
	res, ok := p.pool.Get(key)
	if ok {
		return res
	}
	res = p.newFn(key)
	p.pool.Add(key, res)
	return res
}

// NewKeyPool[K] create a new KeyPool[K] instance
func NewKeyPool[K comparable](size int, newFn func(key K) string) (KeyPool[K], error) {
	pool, err := lru.New[K, string](size)
	if err != nil {
		return nil, err
	}
	return &lruKeyPool[K]{
		pool:  pool,
		newFn: newFn,
	}, nil
}
