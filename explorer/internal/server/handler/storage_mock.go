// Code generated by mockery v2.15.0. DO NOT EDIT.

package handler

import (
	domain "github.com/MiyukiMori11/weatherstat/explorer/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockStorage is an autogenerated mock type for the Storage type
type MockStorage struct {
	mock.Mock
}

// AddCity provides a mock function with given fields: _a0
func (_m *MockStorage) AddCity(_a0 domain.City) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.City) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteCity provides a mock function with given fields: cityName, countryName
func (_m *MockStorage) DeleteCity(cityName string, countryName string) error {
	ret := _m.Called(cityName, countryName)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(cityName, countryName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCountryCode provides a mock function with given fields: countryName
func (_m *MockStorage) GetCountryCode(countryName string) (string, error) {
	ret := _m.Called(countryName)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(countryName)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(countryName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubscribedCities provides a mock function with given fields:
func (_m *MockStorage) GetSubscribedCities() (domain.Cities, error) {
	ret := _m.Called()

	var r0 domain.Cities
	if rf, ok := ret.Get(0).(func() domain.Cities); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(domain.Cities)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTemperatureStat provides a mock function with given fields: cityName, countryName
func (_m *MockStorage) GetTemperatureStat(cityName string, countryName string) (float64, float64, error) {
	ret := _m.Called(cityName, countryName)

	var r0 float64
	if rf, ok := ret.Get(0).(func(string, string) float64); ok {
		r0 = rf(cityName, countryName)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 float64
	if rf, ok := ret.Get(1).(func(string, string) float64); ok {
		r1 = rf(cityName, countryName)
	} else {
		r1 = ret.Get(1).(float64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string) error); ok {
		r2 = rf(cityName, countryName)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewMockStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockStorage creates a new instance of MockStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockStorage(t mockConstructorTestingTNewMockStorage) *MockStorage {
	mock := &MockStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
