// Copyright 2023 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package pool

/* Difference between WorkPool and wpool(github.com/cloudwego/kitex/internal/wpool):
- wpool is a goroutine pool with high reuse rate. The old goroutine will block to wait for new tasks coming.
- simpleWorkPool sample as wpool behavior.
- bufferWorkPool like wpool but add new fn to buffer instead direct run fn.
*/

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/alimy/tryst/types"
)

// WorkPool work pool interface
type WorkPool interface {
	Go(fn types.Fn)
}

// NewWorkPool create a work pool instance
func NewWorkPool(maxIdle int, maxIdleTime time.Duration) WorkPool {
	return NewBufferWorkPool(maxIdle, maxIdleTime)
}

// NewBufferWorkPool create a work pool instance use buffer style
func NewBufferWorkPool(maxIdle int, maxIdleTime time.Duration) WorkPool {
	return &bufferWorkPool{
		maxIdle:        int32(maxIdle),
		maxIdleTime:    maxIdleTime,
		maxBufIdleTime: 5 * maxIdleTime,
		workCh:         make(chan types.Fn),
		workBufCh:      make(chan types.Fn, 1),
	}
}

// NewSimpleWorkPool create a work pool instance use simple style
func NewSimpleWorkPool(maxIdle int, maxIdleTime time.Duration) WorkPool {
	return &simpleWorkPool{
		maxIdle:     int32(maxIdle),
		maxIdleTime: maxIdleTime,
		workCh:      make(chan types.Fn),
	}
}

type simpleWorkPool struct {
	workCh chan types.Fn
	size   atomic.Int32
	// maxIdle is the number of the max idle workers in the pool.
	// if maxIdle too small, the pool works like a native 'go func()'.
	maxIdle     int32
	maxIdleTime time.Duration
}

type bufferWorkPool struct {
	workCh       chan types.Fn
	workBufCh    chan types.Fn
	inBufWorking atomic.Bool
	size         atomic.Int32
	// maxIdle is the number of the max idle workers in the pool.
	// if maxIdle is -1 or 0, the pool workers with idle is no limit.
	maxIdle        int32
	maxIdleTime    time.Duration
	maxBufIdleTime time.Duration
}

func (p *simpleWorkPool) Go(fn types.Fn) {
	if fn == nil {
		return
	}
	select {
	case p.workCh <- fn:
		// send fn by workCh chan
		return
	default:
	}
	go func() {
		p.size.Add(1)
		defer func() {
			if err := recover(); err != nil {
				log.Printf("do fn occurs panic: %s", err)
			}
			p.size.Add(-1)
		}()
		// run fn
		fn()

		if p.size.Load() > p.maxIdle {
			return
		}

		// watch workCh to continue do work if needed.
		idleTimer := time.NewTimer(p.maxIdleTime)
		for {
			select {
			case fn = <-p.workCh:
				fn()
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

func (p *bufferWorkPool) Go(fn types.Fn) {
	if fn == nil {
		return
	}
	select {
	case p.workCh <- fn:
		// send fn by workCh chan
		return
	default:
	}
	if p.maxIdle > 0 && p.size.Load() >= p.maxIdle {
		p.runBufferWorker()
		p.workBufCh <- fn
		return
	}
	go func() {
		p.size.Add(1)
		defer func() {
			if err := recover(); err != nil {
				log.Printf("do fn occurs panic: %s", err)
			}
			p.size.Add(-1)
		}()
		// run fn
		fn()
		// watch workCh to continue do work if needed.
		idleTimer := time.NewTimer(p.maxIdleTime)
		for {
			select {
			case fn = <-p.workCh:
				fn()
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

func (p *bufferWorkPool) runBufferWorker() {
	if !p.inBufWorking.Swap(true) {
		defer p.inBufWorking.Store(false)
		// do the buffer work
		go func() {
			var fnBuf []types.Fn
			// watch workCh to continue do work if needed.
			idleTimer := time.NewTimer(p.maxBufIdleTime)
			for {
				if latesIdx := len(fnBuf) - 1; latesIdx >= 0 {
					select {
					case p.workCh <- fnBuf[0]:
						fnBuf[0] = fnBuf[latesIdx]
						fnBuf = fnBuf[:latesIdx]
					case item := <-p.workBufCh:
						fnBuf = append(fnBuf, item)
					case <-idleTimer.C:
						if p.size.Load() == 0 && len(fnBuf) == 0 {
							return
						}
					}
				} else {
					select {
					case item := <-p.workBufCh:
						fnBuf = append(fnBuf, item)
					case <-idleTimer.C:
						if p.size.Load() == 0 && len(fnBuf) == 0 {
							return
						}
					}
				}
				if !idleTimer.Stop() {
					<-idleTimer.C
				}
				idleTimer.Reset(p.maxBufIdleTime)
			}
		}()
	}
}
