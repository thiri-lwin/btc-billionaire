// Code generated by MockGen. DO NOT EDIT.
// Source: btc_billionaire/postgres (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	postgres "btc_billionaire/postgres"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// GetBTCHistory mocks base method.
func (m *MockStore) GetBTCHistory(arg0, arg1 time.Time) ([]postgres.BTCInfoResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBTCHistory", arg0, arg1)
	ret0, _ := ret[0].([]postgres.BTCInfoResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBTCHistory indicates an expected call of GetBTCHistory.
func (mr *MockStoreMockRecorder) GetBTCHistory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBTCHistory", reflect.TypeOf((*MockStore)(nil).GetBTCHistory), arg0, arg1)
}

// InsertBTCInfo mocks base method.
func (m *MockStore) InsertBTCInfo(arg0 postgres.BTCInfo) (postgres.BTCInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertBTCInfo", arg0)
	ret0, _ := ret[0].(postgres.BTCInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertBTCInfo indicates an expected call of InsertBTCInfo.
func (mr *MockStoreMockRecorder) InsertBTCInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertBTCInfo", reflect.TypeOf((*MockStore)(nil).InsertBTCInfo), arg0)
}