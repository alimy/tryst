// Copyright 2023 Michael Li <alimy@niubiu.com>. All rights reserved.
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
	OnJoin(count int32)
	OnLeave(count int32)
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

// ExecFn[T] request handle function
type ExecFn[T any] func(req T)

// GoroutinePool2[T] goroutine pool interface
type GoroutinePool2[T any] interface {
	Start()
	Stop()
	Run(T, RespFn[T])
}

// GoroutinePool3[T] goroutine pool interface
type GoroutinePool3[T any] interface {
	Start()
	Stop()
	Exec(T)
}

// Option groutine pool option help function used to create groutine pool instance
type Option = func(opt *gorotinePoolOpt)

// grotinePoolOpt gorotine pool option used to create gorotine pool instance
type gorotinePoolOpt struct {
	minWorker          int
	maxTempWorker      int
	maxRequestInCh     int
	maxRequestInTempCh int
	maxIdleTime        time.Duration
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
	requestBufCh       chan *requestItem[T, R] // 请求缓存通道
	maxRequestInCh     int
	maxRequestInTempCh int
	minWorker          int // 最少正式工数
	maxTempWorker      int // 最大临时工数，-1表示无限制
	maxIdleTime        time.Duration
	tempWorkerCount    atomic.Int32
	doFn               DoFn[T, R]
	cancelFn           context.CancelFunc
	workerHook         WorkerHook
}

type wormPool2[T any] struct {
	ctx                context.Context
	isStarted          atomic.Bool
	requestCh          chan *requestItem2[T] // 正式工 缓存通道
	requestTempCh      chan *requestItem2[T] // 临时工 缓存通道
	requestBufCh       chan *requestItem2[T] // 请求缓存通道
	maxRequestInCh     int
	maxRequestInTempCh int
	minWorker          int // 最少正式工数
	maxTempWorker      int // 最大临时工数，-1表示无限制
	maxIdleTime        time.Duration
	tempWorkerCount    atomic.Int32
	runFn              RunFn[T]
	cancelFn           context.CancelFunc
	workerHook         WorkerHook
}

type wormPool3[T any] struct {
	ctx                context.Context
	isStarted          atomic.Bool
	requestCh          chan T // 正式工 缓存通道
	requestTempCh      chan T // 临时工 缓存通道
	requestBufCh       chan T // 请求缓存通道
	maxRequestInCh     int
	maxRequestInTempCh int
	minWorker          int // 最少正式工数
	maxTempWorker      int // 最大临时工数，-1表示无限制
	maxIdleTime        time.Duration
	tempWorkerCount    atomic.Int32
	execFn             ExecFn[T]
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
			if p.maxTempWorker >= 0 && p.tempWorkerCount.Load() >= int32(p.maxTempWorker) {
				p.requestBufCh <- item
				break
			}
			go func() {
				// update temp worker count and run worker hook
				count := p.tempWorkerCount.Add(1)
				if p.workerHook != nil {
					p.workerHook.OnJoin(count)
				}
				defer func() {
					count = p.tempWorkerCount.Add(-1)
					if p.workerHook != nil {
						p.workerHook.OnLeave(count)
					}
				}()
				// handle the request
				p.do(item)
				// watch requestTempCh to continue do work if needed.
				idleTimer := time.NewTimer(p.maxIdleTime)
				for {
					select {
					case item = <-p.requestTempCh:
						p.do(item)
					case <-p.ctx.Done():
						// worker exits
						return
					case <-idleTimer.C:
						// worker exits
						return
					}
					if !idleTimer.Stop() {
						<-idleTimer.C
					}
					idleTimer.Reset(p.maxIdleTime)
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
			if p.maxTempWorker >= 0 && p.tempWorkerCount.Load() >= int32(p.maxTempWorker) {
				p.requestBufCh <- item
				break
			}
			go func() {
				// update temp worker count and run worker hook
				count := p.tempWorkerCount.Add(1)
				if p.workerHook != nil {
					p.workerHook.OnJoin(count)
				}
				defer func() {
					count = p.tempWorkerCount.Add(-1)
					if p.workerHook != nil {
						p.workerHook.OnLeave(count)
					}
				}()
				// handle the request
				p.run(item)
				// watch requestTempCh to continue do work if needed.
				idleTimer := time.NewTimer(p.maxIdleTime)
				for {
					select {
					case item = <-p.requestTempCh:
						p.run(item)
					case <-p.ctx.Done():
						// worker exits
						return
					case <-idleTimer.C:
						// worker exits
						return
					}
					if !idleTimer.Stop() {
						<-idleTimer.C
					}
					idleTimer.Reset(p.maxIdleTime)
				}
			}()
		}
	}
}

func (p *wormPool3[T]) Exec(item T) {
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
			if p.maxTempWorker >= 0 && p.tempWorkerCount.Load() >= int32(p.maxTempWorker) {
				p.requestBufCh <- item
				break
			}
			go func() {
				// update temp worker count and run worker hook
				count := p.tempWorkerCount.Add(1)
				if p.workerHook != nil {
					p.workerHook.OnJoin(count)
				}
				defer func() {
					count = p.tempWorkerCount.Add(-1)
					if p.workerHook != nil {
						p.workerHook.OnLeave(count)
					}
				}()
				// handle the request
				p.exec(item)
				// watch requestTempCh to continue do work if needed.
				idleTimer := time.NewTimer(p.maxIdleTime)
				for {
					select {
					case item = <-p.requestTempCh:
						p.exec(item)
					case <-p.ctx.Done():
						// worker exits
						return
					case <-idleTimer.C:
						// worker exits
						return
					}
					if !idleTimer.Stop() {
						<-idleTimer.C
					}
					idleTimer.Reset(p.maxIdleTime)
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
		if p.maxTempWorker >= 0 {
			p.requestBufCh = make(chan *requestItem[T, R], 1)
			go p.runBufferWorker()
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
		if p.maxTempWorker >= 0 {
			p.requestBufCh = make(chan *requestItem2[T], 1)
			go p.runBufferWorker()
		}
	}
}

func (p *wormPool3[T]) Start() {
	if !p.isStarted.Swap(true) {
		p.ctx, p.cancelFn = context.WithCancel(context.Background())
		p.requestCh = make(chan T, p.maxRequestInCh)
		p.requestTempCh = make(chan T, p.maxRequestInTempCh)
		for numWorker := p.minWorker; numWorker > 0; numWorker-- {
			go p.goExec()
		}
		if p.maxTempWorker >= 0 {
			p.requestBufCh = make(chan T, 1)
			go p.runBufferWorker()
		}
	}
}

func (p *wormPool[T, R]) runBufferWorker() {
	var reqBuf []*requestItem[T, R]
	for {
		if latesIdx := len(reqBuf) - 1; latesIdx >= 0 {
			select {
			case p.requestCh <- reqBuf[0]:
				reqBuf[0] = reqBuf[latesIdx]
				reqBuf = reqBuf[:latesIdx]
			case p.requestTempCh <- reqBuf[0]:
				reqBuf[0] = reqBuf[latesIdx]
				reqBuf = reqBuf[:latesIdx]
			case item := <-p.requestBufCh:
				reqBuf = append(reqBuf, item)
			case <-p.ctx.Done():
				return
			}
		} else {
			select {
			case item := <-p.requestBufCh:
				reqBuf = append(reqBuf, item)
			case <-p.ctx.Done():
				return
			}
		}
	}
}

func (p *wormPool2[T]) runBufferWorker() {
	var reqBuf []*requestItem2[T]
	for {
		if latesIdx := len(reqBuf) - 1; latesIdx >= 0 {
			select {
			case p.requestCh <- reqBuf[0]:
				reqBuf[0] = reqBuf[latesIdx]
				reqBuf = reqBuf[:latesIdx]
			case p.requestTempCh <- reqBuf[0]:
				reqBuf[0] = reqBuf[latesIdx]
				reqBuf = reqBuf[:latesIdx]
			case item := <-p.requestBufCh:
				reqBuf = append(reqBuf, item)
			case <-p.ctx.Done():
				return
			}
		} else {
			select {
			case item := <-p.requestBufCh:
				reqBuf = append(reqBuf, item)
			case <-p.ctx.Done():
				return
			}
		}
	}
}

func (p *wormPool3[T]) runBufferWorker() {
	var reqBuf []T
	for {
		if latesIdx := len(reqBuf) - 1; latesIdx >= 0 {
			select {
			case p.requestCh <- reqBuf[0]:
				reqBuf[0] = reqBuf[latesIdx]
				reqBuf = reqBuf[:latesIdx]
			case p.requestTempCh <- reqBuf[0]:
				reqBuf[0] = reqBuf[latesIdx]
				reqBuf = reqBuf[:latesIdx]
			case item := <-p.requestBufCh:
				reqBuf = append(reqBuf, item)
			case <-p.ctx.Done():
				return
			}
		} else {
			select {
			case item := <-p.requestBufCh:
				reqBuf = append(reqBuf, item)
			case <-p.ctx.Done():
				return
			}
		}
	}
}

func (p *wormPool[T, R]) Stop() {
	if p.isStarted.Swap(false) {
		p.cancelFn()
		close(p.requestCh)
		close(p.requestTempCh)
		if p.maxTempWorker >= 0 {
			close(p.requestBufCh)
		}
	}
}

func (p *wormPool2[T]) Stop() {
	if p.isStarted.Swap(false) {
		p.cancelFn()
		close(p.requestCh)
		close(p.requestTempCh)
		if p.maxTempWorker >= 0 {
			close(p.requestBufCh)
		}
	}
}

func (p *wormPool3[T]) Stop() {
	if p.isStarted.Swap(false) {
		p.cancelFn()
		close(p.requestCh)
		close(p.requestTempCh)
		if p.maxTempWorker >= 0 {
			close(p.requestBufCh)
		}
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

func (p *wormPool3[T]) exec(item T) {
	p.execFn(item)
	defer func() {
		if err := recover(); err != nil {
			// TODO: add log
			// do nothing
		}
	}()
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

func (p *wormPool3[T]) goExec() {
For:
	for {
		select {
		case item := <-p.requestCh:
			p.exec(item)
		case <-p.ctx.Done():
			break For
		}
	}
}

// WithMinWorker set min worker count
func WithMinWorker(num int) Option {
	return func(opt *gorotinePoolOpt) {
		opt.minWorker = num
	}
}

// WithMaxRequestBuf set max temp worker count
func WithMaxTempWorker(num int) Option {
	return func(opt *gorotinePoolOpt) {
		opt.maxTempWorker = num
	}
}

// WithMaxRequestBuf set max request buffer size
func WithMaxRequestBuf(num int) Option {
	return func(opt *gorotinePoolOpt) {
		opt.maxRequestInCh = num
	}
}

// WithMaxRequestTempBuf set max request temp buffer size
func WithMaxRequestTempBuf(num int) Option {
	return func(opt *gorotinePoolOpt) {
		opt.maxRequestInTempCh = num
	}
}

// WithMaxIdelTime set max idle time to custom a worker max wait tile to worker
func WithMaxIdelTime(d time.Duration) Option {
	return func(opt *gorotinePoolOpt) {
		opt.maxIdleTime = d
	}
}

// WithWorkerHookOpt set wroker hook
func WithWorkerHook(h WorkerHook) Option {
	return func(opt *gorotinePoolOpt) {
		opt.workerHook = h
	}
}

// NewGoroutinePool[T, R] create a new GoroutinePool[T, R] instance
func NewGoroutinePool[T, R any](fn DoFn[T, R], opts ...Option) GoroutinePool[T, R] {
	opt := &gorotinePoolOpt{
		minWorker:          10,
		maxTempWorker:      -1,
		maxRequestInCh:     100,
		maxRequestInTempCh: 100,
		maxIdleTime:        60 * time.Second,
	}
	for _, optFn := range opts {
		optFn(opt)
	}
	p := &wormPool[T, R]{
		maxRequestInCh:     opt.maxRequestInCh,
		maxRequestInTempCh: opt.maxRequestInTempCh,
		minWorker:          opt.minWorker,
		maxTempWorker:      opt.maxTempWorker,
		maxIdleTime:        opt.maxIdleTime,
		workerHook:         opt.workerHook,
		doFn:               fn,
	}
	p.Start()
	return p
}

// NewGoroutinePool2[T] create a new GoroutinePool2[T] instance
func NewGoroutinePool2[T any](fn RunFn[T], opts ...Option) GoroutinePool2[T] {
	opt := &gorotinePoolOpt{
		minWorker:          10,
		maxTempWorker:      -1,
		maxRequestInCh:     100,
		maxRequestInTempCh: 100,
		maxIdleTime:        60 * time.Second,
	}
	for _, optFn := range opts {
		optFn(opt)
	}
	p := &wormPool2[T]{
		maxRequestInCh:     opt.maxRequestInCh,
		maxRequestInTempCh: opt.maxRequestInTempCh,
		minWorker:          opt.minWorker,
		maxTempWorker:      opt.maxTempWorker,
		maxIdleTime:        opt.maxIdleTime,
		workerHook:         opt.workerHook,
		runFn:              fn,
	}
	p.Start()
	return p
}

// NewGoroutinePool3[T] create a new GoroutinePool3[T] instance
func NewGoroutinePool3[T any](fn ExecFn[T], opts ...Option) GoroutinePool3[T] {
	opt := &gorotinePoolOpt{
		minWorker:          10,
		maxTempWorker:      -1,
		maxRequestInCh:     100,
		maxRequestInTempCh: 100,
		maxIdleTime:        60 * time.Second,
	}
	for _, optFn := range opts {
		optFn(opt)
	}
	p := &wormPool3[T]{
		maxRequestInCh:     opt.maxRequestInCh,
		maxRequestInTempCh: opt.maxRequestInTempCh,
		minWorker:          opt.minWorker,
		maxTempWorker:      opt.maxTempWorker,
		maxIdleTime:        opt.maxIdleTime,
		workerHook:         opt.workerHook,
		execFn:             fn,
	}
	p.Start()
	return p
}
