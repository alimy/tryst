// Copyright 2023 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package event

// Event event interface
type Event interface {
	Name() string
	Before() error
	Action() error
	After() error

	mustEmbedUnimplementedEvent()
}

// Event2[T] event2 interface
type Event2[T any] interface {
	Name() string
	Before() error
	Handle() (T, error)
	After() error

	mustEmbedUnimplementedEvent2()
}

// EventManger event manager
type EventManager interface {
	Start()
	Stop()
	OnEvent(event Event)
}

// EventManger2[T] event manager
type EventManager2[T any] interface {
	Start()
	Stop()
	OnEvent(event Event2[T])
}

// UnimplementedEvent unimplemented Event
type UnimplementedEvent struct{}

func (UnimplementedEvent) Name() string {
	return "UnimplementedEvent"
}

func (UnimplementedEvent) Before() error {
	// do nothing
	return nil
}

func (UnimplementedEvent) After() error {
	// do nothing
	return nil
}

func (UnimplementedEvent) mustEmbedUnimplementedEvent() {}

// UnimplementedEvent unimplemented Event2
type UnimplementedEvent2 struct{}

func (UnimplementedEvent2) Name() string {
	return "UnimplementedEvent2"
}

func (UnimplementedEvent2) Before() error {
	// do nothing
	return nil
}

func (UnimplementedEvent2) After() error {
	// do nothing
	return nil
}

func (UnimplementedEvent2) mustEmbedUnimplementedEvent2() {}
