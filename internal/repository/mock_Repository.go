// Code generated by mockery v2.10.0. DO NOT EDIT.

package repository

import mock "github.com/stretchr/testify/mock"

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: name
func (_m *MockRepository) Delete(name string) error {
	ret := _m.Called(name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Destroy provides a mock function with given fields:
func (_m *MockRepository) Destroy() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: query
func (_m *MockRepository) Find(query *Query) ([]*Cat, error) {
	ret := _m.Called(query)

	var r0 []*Cat
	if rf, ok := ret.Get(0).(func(*Query) []*Cat); ok {
		r0 = rf(query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Cat)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*Query) error); ok {
		r1 = rf(query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: name
func (_m *MockRepository) Get(name string) (*Cat, error) {
	ret := _m.Called(name)

	var r0 *Cat
	if rf, ok := ret.Get(0).(func(string) *Cat); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Cat)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: cat
func (_m *MockRepository) Insert(cat *Cat) error {
	ret := _m.Called(cat)

	var r0 error
	if rf, ok := ret.Get(0).(func(*Cat) error); ok {
		r0 = rf(cat)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: catUpdate
func (_m *MockRepository) Update(catUpdate *CatUpdate) error {
	ret := _m.Called(catUpdate)

	var r0 error
	if rf, ok := ret.Get(0).(func(*CatUpdate) error); ok {
		r0 = rf(catUpdate)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}