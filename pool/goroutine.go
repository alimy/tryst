// Copyright 2023 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package pool

import (
	"time"
)

// ResponseFn[T, R] response handle function
type ResponseFn[T, R any] func(req T, resp R, err error)

// GoFn[T, R] request handle function
type GoFn[T, R any] func(req T) (R, error)

// GoroutinePool[T, R] goroutine pool interface
type GoroutinePool[T, R any] interface {
	Do(T, ResponseFn[T, R])
}

// GorotinePoolOptFn groutine pool option help function used to create groutine pool instance
type GorotinePoolOptFn = func(opt *gorotinePoolOpt)

// MinWorkerOpt set min worker
func MinWorkerOpt(num int) GorotinePoolOptFn {
	return func(opt *gorotinePoolOpt) {
		opt.MinWorker = num
	}
}

// MaxRequestBuffer set max request buffer size
func MaxRequestBuffer(num int) GorotinePoolOptFn {
	return func(opt *gorotinePoolOpt) {
		opt.MaxRequestInCh = num
	}
}

// MaxRequestTempBuffer set max request temp buffer size
func MaxRequestTempBuffer(num int) GorotinePoolOptFn {
	return func(opt *gorotinePoolOpt) {
		opt.MaxRequestInTempCh = num
	}
}

// MaxTickCount set max tick count
func MaxTickCount(num int) GorotinePoolOptFn {
	return func(opt *gorotinePoolOpt) {
		opt.MaxTickCount = num
	}
}

// TickWaitTime set tick wait time
func TickWaitTime(duration time.Duration) GorotinePoolOptFn {
	return func(opt *gorotinePoolOpt) {
		opt.TickWaitTime = duration
	}
}

// NewGoroutinePool[T, R] create a new GoroutinePool[T, R] instance
func NewGoroutinePool[T, R any](fn GoFn[T, R], opts ...GorotinePoolOptFn) GoroutinePool[T, R] {
	opt := &gorotinePoolOpt{
		MinWorker:          10,
		MaxRequestInCh:     100,
		MaxRequestInTempCh: 100,
		MaxTickCount:       60,
		TickWaitTime:       time.Second,
	}
	for _, optFn := range opts {
		optFn(opt)
	}
	p := &wormPool[T, R]{
		requestCh:     make(chan *requestItem[T, R], opt.MaxRequestInCh),
		requestTempCh: make(chan *requestItem[T, R], opt.MaxRequestInTempCh),
		maxTickCount:  opt.MaxTickCount,
		tickWaitTime:  opt.TickWaitTime,
		goFn:          fn,
	}
	p.startDoWork()
	return p
}

// grotinePoolOpt gorotine pool option used to create gorotine pool instance
type gorotinePoolOpt struct {
	MinWorker          int
	MaxRequestInCh     int
	MaxRequestInTempCh int
	MaxTickCount       int
	TickWaitTime       time.Duration
}

type requestItem[T, R any] struct {
	req    T
	respFn ResponseFn[T, R]
}

type wormPool[T, R any] struct {
	requestCh     chan *requestItem[T, R] // 正式工 缓存通道
	requestTempCh chan *requestItem[T, R] // 临时工 缓存通道
	minWorker     int                     // 最少正式工数
	maxTickCount  int
	tickWaitTime  time.Duration
	goFn          GoFn[T, R]
}

func (p *wormPool[T, R]) Do(req T, fn ResponseFn[T, R]) {
	item := &requestItem[T, R]{req, fn}
	select {
	case p.requestCh <- item:
		// send request item by requestCh chan
	default:
		select {
		case p.requestTempCh <- item:
			// send request item by requestTempCh chan"
		default:
			go func() {
				p.do(item)
				// watch requestTempCh to continue do work if needed.
				// cancel loop if no item had watched in s.maxTickCount * s.tickWaitTime.
				for count := 0; count < p.maxTickCount; count++ {
					select {
					case item := <-p.requestTempCh:
						// reset count to continue do work
						count = 0
						p.do(item)
					default:
						// sleeping to wait request item pass over to do work
						time.Sleep(p.tickWaitTime)
					}
				}
			}()
		}
	}
}

func (p *wormPool[T, R]) do(item *requestItem[T, R]) {
	resp, err := p.goFn(item.req)
	item.respFn(item.req, resp, err)
}

// startDoWork start do work
func (p *wormPool[T, R]) startDoWork() {
	for numWorker := p.minWorker; numWorker > 0; numWorker-- {
		go p.goDo()
	}
}

func (p *wormPool[T, R]) goDo() {
	for item := range p.requestCh {
		p.do(item)
	}
}
