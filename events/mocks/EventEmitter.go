// Code generated by mockery v2.49.0. DO NOT EDIT.

package mocks

import (
	events "golangsevillabar/events"

	mock "github.com/stretchr/testify/mock"
)

// EventEmitter is an autogenerated mock type for the EventEmitter type
type EventEmitter struct {
	mock.Mock
}

// EmitEvent provides a mock function with given fields: event
func (_m *EventEmitter) EmitEvent(event events.Event) error {
	ret := _m.Called(event)

	if len(ret) == 0 {
		panic("no return value specified for EmitEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(events.Event) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewEventEmitter creates a new instance of EventEmitter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventEmitter(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventEmitter {
	mock := &EventEmitter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
