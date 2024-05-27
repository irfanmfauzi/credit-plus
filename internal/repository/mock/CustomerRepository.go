// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"
	repository "credit-plus/internal/repository"

	mock "github.com/stretchr/testify/mock"

	request "credit-plus/internal/model/request"
)

// CustomerRepository is an autogenerated mock type for the CustomerRepository type
type CustomerRepository struct {
	mock.Mock
}

// CheckLimitTransaction provides a mock function with given fields: ctx, tx, otrPrice, customerId
func (_m *CustomerRepository) CheckLimitTransaction(ctx context.Context, tx repository.TxProvider, otrPrice int, customerId int64) (bool, error) {
	ret := _m.Called(ctx, tx, otrPrice, customerId)

	if len(ret) == 0 {
		panic("no return value specified for CheckLimitTransaction")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.TxProvider, int, int64) (bool, error)); ok {
		return rf(ctx, tx, otrPrice, customerId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.TxProvider, int, int64) bool); ok {
		r0 = rf(ctx, tx, otrPrice, customerId)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.TxProvider, int, int64) error); ok {
		r1 = rf(ctx, tx, otrPrice, customerId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertContract provides a mock function with given fields: ctx, tx, req
func (_m *CustomerRepository) InsertContract(ctx context.Context, tx repository.TxProvider, req request.CreateContactRequest) (int64, error) {
	ret := _m.Called(ctx, tx, req)

	if len(ret) == 0 {
		panic("no return value specified for InsertContract")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.TxProvider, request.CreateContactRequest) (int64, error)); ok {
		return rf(ctx, tx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.TxProvider, request.CreateContactRequest) int64); ok {
		r0 = rf(ctx, tx, req)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.TxProvider, request.CreateContactRequest) error); ok {
		r1 = rf(ctx, tx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertCustomer provides a mock function with given fields: ctx, tx, req
func (_m *CustomerRepository) InsertCustomer(ctx context.Context, tx repository.TxProvider, req request.CreateCustomerRequest) error {
	ret := _m.Called(ctx, tx, req)

	if len(ret) == 0 {
		panic("no return value specified for InsertCustomer")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.TxProvider, request.CreateCustomerRequest) error); ok {
		r0 = rf(ctx, tx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateLimitTransaction provides a mock function with given fields: ctx, tx, otrPrice, customerId
func (_m *CustomerRepository) UpdateLimitTransaction(ctx context.Context, tx repository.TxProvider, otrPrice int, customerId int64) error {
	ret := _m.Called(ctx, tx, otrPrice, customerId)

	if len(ret) == 0 {
		panic("no return value specified for UpdateLimitTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.TxProvider, int, int64) error); ok {
		r0 = rf(ctx, tx, otrPrice, customerId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewCustomerRepository creates a new instance of CustomerRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCustomerRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *CustomerRepository {
	mock := &CustomerRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}