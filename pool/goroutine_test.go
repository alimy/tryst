// Copyright 2023 Michael Li <alimy@gility.net>. All rights reserved.
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
		maxTickCount:       60,
		tickWaitTime:       time.Second,
	}
	opt := gorotinePoolOpt{}
	for _, optFn := range []Option{
		MinWorkerOpt(10),
		MaxRequestBufOpt(100),
		MaxRequestTempBufOpt(100),
		MaxTickCountOpt(60),
		TickWaitTimeOpt(time.Second),
	} {
		optFn(&opt)
	}
	if opt != expectOpt {
		t.Errorf("want %+v but got %+v", expectOpt, opt)
	}
}
