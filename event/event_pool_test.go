// Copyright 2023 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package event

import (
	"sync/atomic"
	"testing"
	"time"
)

var totalCount, totalCount2 atomic.Int32

type fakeEvent struct {
	UnimplementedEvent
	count int32
}

type fakeEvent2 struct {
	UnimplementedEvent2
	count int32
}

func (e *fakeEvent) Name() string {
	return "fakeEvent"
}

func (e *fakeEvent) Action() error {
	totalCount.Add(e.count)
	return nil
}

func (e *fakeEvent2) Name() string {
	return "fakeEvent2"
}

func (e *fakeEvent2) Handle() (int32, error) {
	return totalCount2.Add(e.count), nil
}

func TestEventManager(t *testing.T) {
	em := NewEventManager(func(Event, error) {
		// do nothing
	})
	for i := 0; i < 100; i++ {
		evt := &fakeEvent{count: 1}
		em.OnEvent(evt)
	}
	em.Stop()
	time.Sleep(5 * time.Second)
	if count := totalCount.Load(); count != 100 {
		t.Errorf("expect total count equel 100 but got %d", count)
	}
	em.Start()
	for i := 0; i < 100; i++ {
		evt := &fakeEvent{count: 1}
		em.OnEvent(evt)
	}
	em.Stop()
	time.Sleep(5 * time.Second)
	if count := totalCount.Load(); count != 200 {
		t.Errorf("expect total count equel 200 but got %d", count)
	}
}

func TestEventManager2(t *testing.T) {
	em := NewEventManager2(func(Event2[int32], int32, error) {
		// do nothing
	})
	for i := 0; i < 100; i++ {
		evt := &fakeEvent2{count: 1}
		em.OnEvent(evt)
	}
	em.Stop()
	time.Sleep(5 * time.Second)
	if count := totalCount2.Load(); count != 100 {
		t.Errorf("expect total count equel 100 but got %d", count)
	}
	em.Start()
	for i := 0; i < 100; i++ {
		evt := &fakeEvent2{count: 1}
		em.OnEvent(evt)
	}
	em.Stop()
	time.Sleep(5 * time.Second)
	if count := totalCount2.Load(); count != 200 {
		t.Errorf("expect total count equel 200 but got %d", count)
	}
}
