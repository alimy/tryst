// Copyright 2023 Michael Li <alimy@niubiu.com>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package pool

import (
	"testing"
	"time"
)

func TestGoroutinePoolOpt(t *testing.T) {
	expectOpt := gorotinePoolOpt{
		minWorker:          10,
		maxRequestInCh:     100,
		maxRequestInTempCh: 100,
		maxIdleTime:        60 * time.Second,
	}
	opt := gorotinePoolOpt{}
	for _, optFn := range []Option{
		WithMinWorker(10),
		WithMaxRequestBuf(100),
		WithMaxRequestTempBuf(100),
		WithMaxIdelTime(60 * time.Second),
	} {
		optFn(&opt)
	}
	if opt != expectOpt {
		t.Errorf("want %+v but got %+v", expectOpt, opt)
	}
}
