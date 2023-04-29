// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/isaias-dgr/currency-tracker/internal/core/domain"
	mock "github.com/stretchr/testify/mock"
)

// CurrencyAPIService is an autogenerated mock type for the CurrencyAPIService type
type CurrencyAPIService struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0
func (_m *CurrencyAPIService) Get(_a0 context.Context) (domain.Currencies, error) {
	ret := _m.Called(_a0)

	var r0 domain.Currencies
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (domain.Currencies, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) domain.Currencies); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.Currencies)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCurrencyAPIService interface {
	mock.TestingT
	Cleanup(func())
}

// NewCurrencyAPIService creates a new instance of CurrencyAPIService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCurrencyAPIService(t mockConstructorTestingTNewCurrencyAPIService) *CurrencyAPIService {
	mock := &CurrencyAPIService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
