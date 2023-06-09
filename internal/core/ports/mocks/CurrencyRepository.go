// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/isaias-dgr/currency-tracker/internal/core/domain"
	mock "github.com/stretchr/testify/mock"
)

// CurrencyRepository is an autogenerated mock type for the CurrencyRepository type
type CurrencyRepository struct {
	mock.Mock
}

// GetByCode provides a mock function with given fields: _a0, _a1, _a2
func (_m *CurrencyRepository) GetByCode(_a0 context.Context, _a1 string, _a2 domain.Filter) ([]*domain.CurrencyRepository, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 []*domain.CurrencyRepository
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.Filter) ([]*domain.CurrencyRepository, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, domain.Filter) []*domain.CurrencyRepository); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.CurrencyRepository)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, domain.Filter) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertBulk provides a mock function with given fields: ctx, ta
func (_m *CurrencyRepository) InsertBulk(ctx context.Context, ta domain.Currencies) error {
	ret := _m.Called(ctx, ta)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Currencies) error); ok {
		r0 = rf(ctx, ta)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCurrencyRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewCurrencyRepository creates a new instance of CurrencyRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCurrencyRepository(t mockConstructorTestingTNewCurrencyRepository) *CurrencyRepository {
	mock := &CurrencyRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
