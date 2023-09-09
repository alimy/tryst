// Copyright 2023 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package event

import (
	"github.com/alimy/tryst/pool"
)

type eventPool struct {
	pool   pool.GoroutinePool2[Event]
	respFn pool.RespFn[Event]
}

type eventPool2[T any] struct {
	pool   pool.GoroutinePool[Event2[T], T]
	respFn pool.ResponseFn[Event2[T], T]
}

func (p *eventPool) Start() {
	p.pool.Start()
}

func (p *eventPool2[T]) Start() {
	p.pool.Start()
}

func (p *eventPool) Stop() {
	p.pool.Stop()
}

func (p *eventPool2[T]) Stop() {
	p.pool.Stop()
}

func (p *eventPool) OnEvent(event Event) {
	p.pool.Run(event, p.respFn)
}

func (p *eventPool2[T]) OnEvent(event Event2[T]) {
	p.pool.Do(event, p.respFn)
}

// NewEventManager create new event manager instance
func NewEventManager(respFn pool.RespFn[Event], opts ...pool.Option) (res EventManager) {
	res = &eventPool{
		respFn: respFn,
		pool: pool.NewGoroutinePool2(func(event Event) (err error) {
			if err = event.Before(); err != nil {
				return
			}
			if err = event.Action(); err != nil {
				return
			}
			return event.After()
		}, opts...),
	}
	res.Start()
	return
}

// NewEventManager2[T] create new event manager2 instance
func NewEventManager2[T any](respFn pool.ResponseFn[Event2[T], T], opts ...pool.Option) (res EventManager2[T]) {
	res = &eventPool2[T]{
		respFn: respFn,
		pool: pool.NewGoroutinePool(func(event Event2[T]) (res T, err error) {
			if err = event.Before(); err != nil {
				return
			}
			if res, err = event.Handle(); err != nil {
				return
			}
			err = event.After()
			return
		}, opts...),
	}
	res.Start()
	return
}
