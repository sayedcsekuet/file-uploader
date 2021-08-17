// Code generated by mockery 2.7.4. DO NOT EDIT.

package mocks

import (
	scanner "file-uploader/src/services/scanner"
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// Scanner is an autogenerated mock type for the Scanner type
type Scanner struct {
	mock.Mock
}

// Address provides a mock function with given fields:
func (_m *Scanner) Address() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// HasVirus provides a mock function with given fields: reader
func (_m *Scanner) HasVirus(reader io.Reader) (bool, error) {
	ret := _m.Called(reader)

	var r0 bool
	if rf, ok := ret.Get(0).(func(io.Reader) bool); ok {
		r0 = rf(reader)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(io.Reader) error); ok {
		r1 = rf(reader)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Ping provides a mock function with given fields:
func (_m *Scanner) Ping() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Scan provides a mock function with given fields: reader
func (_m *Scanner) Scan(reader io.Reader) (*scanner.Result, error) {
	ret := _m.Called(reader)

	var r0 *scanner.Result
	if rf, ok := ret.Get(0).(func(io.Reader) *scanner.Result); ok {
		r0 = rf(reader)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*scanner.Result)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(io.Reader) error); ok {
		r1 = rf(reader)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetAddress provides a mock function with given fields: address
func (_m *Scanner) SetAddress(address string) {
	_m.Called(address)
}

// Version provides a mock function with given fields:
func (_m *Scanner) Version() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
