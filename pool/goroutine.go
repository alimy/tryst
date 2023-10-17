// Copyright 2023 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package pool

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

// WorkerHook hook worker status
type WorkerHook interface {
	OnJoin(count int)
	OnLeave(count int)
}

// ResponseFn[T, R] response handle function
type ResponseFn[T, R any] func(req T, resp R, err error)

// DoFn[T, R] request handle function
type DoFn[T, R any] func(req T) (R, error)

// GoroutinePool[T, R] goroutine pool interface
type GoroutinePool[T, R any] interface {
	Start()
	Stop()
	Do(T, ResponseFn[T, R])
}

// RespFn[T] response handle function
type RespFn[T any] func(req T, err error)

// RunFn[T] request handle function
type RunFn[T any] func(req T) error

// GoroutinePool2[T] goroutine pool interface
type GoroutinePool2[T any] interface {
	Start()
	Stop()
	Run(T, RespFn[T])
}

// Option groutine pool option help function used to create groutine pool instance
type Option = func(opt *gorotinePoolOpt)

// grotinePoolOpt gorotine pool option used to create gorotine pool instance
type gorotinePoolOpt struct {
	minWorker          int
	maxRequestInCh     int
	maxRequestInTempCh int
	maxTickCount       int
	tickWaitTime       time.Duration
	workerHook         WorkerHook
}

type requestItem[T, R any] struct {
	req    T
	respFn ResponseFn[T, R]
}

type requestItem2[T any] struct {
	req    T
	respFn RespFn[T]
}

type wormPool[T, R any] struct {
	ctx                context.Context
	isStarted          atomic.Bool
	requestCh          chan *requestItem[T, R] // 正式工 缓存通道
	requestTempCh      chan *requestItem[T, R] // 临时工 缓存通道
	maxRequestInCh     int
	maxRequestInTempCh int
	minWorker          int // 最少正式工数
	maxTickCount       int
	tempWorkerCount    int
	tickWaitTime       time.Duration
	doFn               DoFn[T, R]
	cancelFn           context.CancelFunc
	workerHook         WorkerHook
}

type wormPool2[T any] struct {
	ctx                context.Context
	isStarted          atomic.Bool
	requestCh          chan *requestItem2[T] // 正式工 缓存通道
	requestTempCh      chan *requestItem2[T] // 临时工 缓存通道
	maxRequestInCh     int
	maxRequestInTempCh int
	minWorker          int // 最少正式工数
	maxTickCount       int
	tempWorkerCount    int
	tickWaitTime       time.Duration
	runFn              RunFn[T]
	cancelFn           context.CancelFunc
	workerHook         WorkerHook
}

func (p *wormPool[T, R]) Do(req T, fn ResponseFn[T, R]) {
	item := &requestItem[T, R]{req, fn}
	select {
	case p.requestCh <- item:
		// send request item by requestCh chan
	case <-p.ctx.Done():
		// do nothing
	default:
		select {
		case p.requestTempCh <- item:
			// send request item by requestTempCh chan"
		default:
			go func() {
				// update temp worker count and run worker hook
				p.tempWorkerCount++
				if p.workerHook != nil {
					p.workerHook.OnJoin(p.tempWorkerCount)
				}
				defer func() {
					p.tempWorkerCount--
					if p.workerHook != nil {
						p.workerHook.OnLeave(p.tempWorkerCount)
					}
				}()
				// handle the request
				p.do(item)
				// watch requestTempCh to continue do work if needed.
				// cancel loop if no item had watched in s.maxTickCount * s.tickWaitTime.
			For:
				for count := 0; count < p.maxTickCount; count++ {
					select {
					case item := <-p.requestTempCh:
						// reset count to continue do work
						count = 0
						p.do(item)
					case <-p.ctx.Done():
						break For
					default:
						// sleeping to wait request item pass over to do work
						time.Sleep(p.tickWaitTime)
					}
				}
			}()
		}
	}
}

func (p *wormPool2[T]) Run(req T, fn RespFn[T]) {
	item := &requestItem2[T]{req, fn}
	select {
	case p.requestCh <- item:
		// send request item by requestCh chan
	case <-p.ctx.Done():
		// do nothing
	default:
		select {
		case p.requestTempCh <- item:
			// send request item by requestTempCh chan"
		default:
			go func() {
				// update temp worker count and run worker hook
				p.tempWorkerCount++
				if p.workerHook != nil {
					p.workerHook.OnJoin(p.tempWorkerCount)
				}
				defer func() {
					p.tempWorkerCount--
					if p.workerHook != nil {
						p.workerHook.OnLeave(p.tempWorkerCount)
					}
				}()
				// handle the request
				p.run(item)
				// watch requestTempCh to continue do work if needed.
				// cancel loop if no item had watched in s.maxTickCount * s.tickWaitTime.
			For:
				for count := 0; count < p.maxTickCount; count++ {
					select {
					case item := <-p.requestTempCh:
						// reset count to continue do work
						count = 0
						p.run(item)
					case <-p.ctx.Done():
						break For
					default:
						// sleeping to wait request item pass over to do work
						time.Sleep(p.tickWaitTime)
					}
				}
			}()
		}
	}
}

func (p *wormPool[T, R]) Start() {
	if !p.isStarted.Swap(true) {
		p.ctx, p.cancelFn = context.WithCancel(context.Background())
		p.requestCh = make(chan *requestItem[T, R], p.maxRequestInCh)
		p.requestTempCh = make(chan *requestItem[T, R], p.maxRequestInTempCh)
		for numWorker := p.minWorker; numWorker > 0; numWorker-- {
			go p.goDo()
		}
	}
}

func (p *wormPool2[T]) Start() {
	if !p.isStarted.Swap(true) {
		p.ctx, p.cancelFn = context.WithCancel(context.Background())
		p.requestCh = make(chan *requestItem2[T], p.maxRequestInCh)
		p.requestTempCh = make(chan *requestItem2[T], p.maxRequestInTempCh)
		for numWorker := p.minWorker; numWorker > 0; numWorker-- {
			go p.goRun()
		}
	}
}

func (p *wormPool[T, R]) Stop() {
	if p.isStarted.Swap(false) {
		p.cancelFn()
		close(p.requestCh)
		close(p.requestTempCh)
	}
}

func (p *wormPool2[T]) Stop() {
	if p.isStarted.Swap(false) {
		p.cancelFn()
		close(p.requestCh)
		close(p.requestTempCh)
	}
}

func (p *wormPool[T, R]) do(item *requestItem[T, R]) {
	if item != nil {
		resp, err := p.doFn(item.req)
		item.respFn(item.req, resp, err)
		defer func() {
			if err := recover(); err != nil {
				item.respFn(item.req, resp, fmt.Errorf("do fn occurs panic: %s", err))
			}
		}()
	}
}

func (p *wormPool2[T]) run(item *requestItem2[T]) {
	if item != nil {
		item.respFn(item.req, p.runFn(item.req))
		defer func() {
			if err := recover(); err != nil {
				item.respFn(item.req, fmt.Errorf("run fn occurs panic: %s", err))
			}
		}()
	}
}

func (p *wormPool[T, R]) goDo() {
For:
	for {
		select {
		case item := <-p.requestCh:
			p.do(item)
		case <-p.ctx.Done():
			break For
		}
	}
}

func (p *wormPool2[T]) goRun() {
For:
	for {
		select {
		case item := <-p.requestCh:
			p.run(item)
		case <-p.ctx.Done():
			break For
		}
	}
}

// MinWorkerOpt set min worker
func MinWorkerOpt(num int) Option {
	return func(opt *gorotinePoolOpt) {
		opt.minWorker = num
	}
}

// MaxRequestBufOpt set max request buffer size
func MaxRequestBufOpt(num int) Option {
	return func(opt *gorotinePoolOpt) {
		opt.maxRequestInCh = num
	}
}

// MaxRequestTempBufOpt set max request temp buffer size
func MaxRequestTempBufOpt(num int) Option {
	return func(opt *gorotinePoolOpt) {
		opt.maxRequestInTempCh = num
	}
}

// MaxTickCountOpt set max tick count
func MaxTickCountOpt(num int) Option {
	return func(opt *gorotinePoolOpt) {
		opt.maxTickCount = num
	}
}

// TickWaitTimeOpt set tick wait time
func TickWaitTimeOpt(duration time.Duration) Option {
	return func(opt *gorotinePoolOpt) {
		opt.tickWaitTime = duration
	}
}

// WorkerHookOpt set wroker hook
func WorkerHookOpt(h WorkerHook) Option {
	return func(opt *gorotinePoolOpt) {
		opt.workerHook = h
	}
}

// NewGoroutinePool[T, R] create a new GoroutinePool[T, R] instance
func NewGoroutinePool[T, R any](fn DoFn[T, R], opts ...Option) GoroutinePool[T, R] {
	opt := &gorotinePoolOpt{
		minWorker:          10,
		maxRequestInCh:     100,
		maxRequestInTempCh: 100,
		maxTickCount:       60,
		tickWaitTime:       time.Second,
	}
	for _, optFn := range opts {
		optFn(opt)
	}
	p := &wormPool[T, R]{
		maxRequestInCh:     opt.maxRequestInCh,
		maxRequestInTempCh: opt.maxRequestInTempCh,
		minWorker:          opt.minWorker,
		maxTickCount:       opt.maxTickCount,
		tickWaitTime:       opt.tickWaitTime,
		workerHook:         opt.workerHook,
		doFn:               fn,
	}
	p.Start()
	return p
}

// NewGoroutinePool2[T] create a new GoroutinePool[T, R] instance
func NewGoroutinePool2[T any](fn RunFn[T], opts ...Option) GoroutinePool2[T] {
	opt := &gorotinePoolOpt{
		minWorker:          10,
		maxRequestInCh:     100,
		maxRequestInTempCh: 100,
		maxTickCount:       60,
		tickWaitTime:       time.Second,
	}
	for _, optFn := range opts {
		optFn(opt)
	}
	p := &wormPool2[T]{
		maxRequestInCh:     opt.maxRequestInCh,
		maxRequestInTempCh: opt.maxRequestInTempCh,
		minWorker:          opt.minWorker,
		maxTickCount:       opt.maxTickCount,
		tickWaitTime:       opt.tickWaitTime,
		workerHook:         opt.workerHook,
		runFn:              fn,
	}
	p.Start()
	return p
}
