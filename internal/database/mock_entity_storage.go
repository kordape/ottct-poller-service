// Code generated by mockery. DO NOT EDIT.

package database

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockEntityStorage is an autogenerated mock type for the EntityStorage type
type MockEntityStorage struct {
	mock.Mock
}

// GetEntities provides a mock function with given fields: _a0
func (_m *MockEntityStorage) GetEntities(_a0 context.Context) ([]Entity, error) {
	ret := _m.Called(_a0)

	var r0 []Entity
	if rf, ok := ret.Get(0).(func(context.Context) []Entity); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Entity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type NewMockEntityStorageT interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockEntityStorage creates a new instance of MockEntityStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockEntityStorage(t NewMockEntityStorageT) *MockEntityStorage {
	mock := &MockEntityStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
