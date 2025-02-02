// Code generated by mockery v2.49.0. DO NOT EDIT.

package mocks

import (
	events "cqrseventsourcingbar/events"

	ksuid "github.com/segmentio/ksuid"

	mock "github.com/stretchr/testify/mock"

	queries "cqrseventsourcingbar/queries"
)

// OpenTabQueries is an autogenerated mock type for the OpenTabQueries type
type OpenTabQueries struct {
	mock.Mock
}

// ActiveTableNumbers provides a mock function with given fields:
func (_m *OpenTabQueries) ActiveTableNumbers() []int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ActiveTableNumbers")
	}

	var r0 []int
	if rf, ok := ret.Get(0).(func() []int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	return r0
}

// HandleEvent provides a mock function with given fields: e
func (_m *OpenTabQueries) HandleEvent(e events.Event) error {
	ret := _m.Called(e)

	if len(ret) == 0 {
		panic("no return value specified for HandleEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(events.Event) error); ok {
		r0 = rf(e)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InvoiceForTable provides a mock function with given fields: table
func (_m *OpenTabQueries) InvoiceForTable(table int) (queries.TabInvoice, error) {
	ret := _m.Called(table)

	if len(ret) == 0 {
		panic("no return value specified for InvoiceForTable")
	}

	var r0 queries.TabInvoice
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (queries.TabInvoice, error)); ok {
		return rf(table)
	}
	if rf, ok := ret.Get(0).(func(int) queries.TabInvoice); ok {
		r0 = rf(table)
	} else {
		r0 = ret.Get(0).(queries.TabInvoice)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(table)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TabForTable provides a mock function with given fields: table
func (_m *OpenTabQueries) TabForTable(table int) (queries.TabStatus, error) {
	ret := _m.Called(table)

	if len(ret) == 0 {
		panic("no return value specified for TabForTable")
	}

	var r0 queries.TabStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (queries.TabStatus, error)); ok {
		return rf(table)
	}
	if rf, ok := ret.Get(0).(func(int) queries.TabStatus); ok {
		r0 = rf(table)
	} else {
		r0 = ret.Get(0).(queries.TabStatus)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(table)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TabIdForTable provides a mock function with given fields: table
func (_m *OpenTabQueries) TabIdForTable(table int) (ksuid.KSUID, error) {
	ret := _m.Called(table)

	if len(ret) == 0 {
		panic("no return value specified for TabIdForTable")
	}

	var r0 ksuid.KSUID
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (ksuid.KSUID, error)); ok {
		return rf(table)
	}
	if rf, ok := ret.Get(0).(func(int) ksuid.KSUID); ok {
		r0 = rf(table)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ksuid.KSUID)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(table)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TodoListForWaiter provides a mock function with given fields: waiter
func (_m *OpenTabQueries) TodoListForWaiter(waiter string) map[int][]queries.TabItem {
	ret := _m.Called(waiter)

	if len(ret) == 0 {
		panic("no return value specified for TodoListForWaiter")
	}

	var r0 map[int][]queries.TabItem
	if rf, ok := ret.Get(0).(func(string) map[int][]queries.TabItem); ok {
		r0 = rf(waiter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[int][]queries.TabItem)
		}
	}

	return r0
}

// NewOpenTabQueries creates a new instance of OpenTabQueries. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOpenTabQueries(t interface {
	mock.TestingT
	Cleanup(func())
}) *OpenTabQueries {
	mock := &OpenTabQueries{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
