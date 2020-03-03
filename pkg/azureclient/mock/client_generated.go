// Code generated by MockGen. DO NOT EDIT.
// Source: ./client.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	compute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
	gomock "github.com/golang/mock/gomock"
	azureclient "github.com/openshift/hive/pkg/azureclient"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// ListResourceSKUs mocks base method
func (m *MockClient) ListResourceSKUs(ctx context.Context) (azureclient.ResourceSKUsPage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListResourceSKUs", ctx)
	ret0, _ := ret[0].(azureclient.ResourceSKUsPage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListResourceSKUs indicates an expected call of ListResourceSKUs
func (mr *MockClientMockRecorder) ListResourceSKUs(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListResourceSKUs", reflect.TypeOf((*MockClient)(nil).ListResourceSKUs), ctx)
}

// MockResourceSKUsPage is a mock of ResourceSKUsPage interface
type MockResourceSKUsPage struct {
	ctrl     *gomock.Controller
	recorder *MockResourceSKUsPageMockRecorder
}

// MockResourceSKUsPageMockRecorder is the mock recorder for MockResourceSKUsPage
type MockResourceSKUsPageMockRecorder struct {
	mock *MockResourceSKUsPage
}

// NewMockResourceSKUsPage creates a new mock instance
func NewMockResourceSKUsPage(ctrl *gomock.Controller) *MockResourceSKUsPage {
	mock := &MockResourceSKUsPage{ctrl: ctrl}
	mock.recorder = &MockResourceSKUsPageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockResourceSKUsPage) EXPECT() *MockResourceSKUsPageMockRecorder {
	return m.recorder
}

// NextWithContext mocks base method
func (m *MockResourceSKUsPage) NextWithContext(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NextWithContext", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// NextWithContext indicates an expected call of NextWithContext
func (mr *MockResourceSKUsPageMockRecorder) NextWithContext(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextWithContext", reflect.TypeOf((*MockResourceSKUsPage)(nil).NextWithContext), ctx)
}

// NotDone mocks base method
func (m *MockResourceSKUsPage) NotDone() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotDone")
	ret0, _ := ret[0].(bool)
	return ret0
}

// NotDone indicates an expected call of NotDone
func (mr *MockResourceSKUsPageMockRecorder) NotDone() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotDone", reflect.TypeOf((*MockResourceSKUsPage)(nil).NotDone))
}

// Values mocks base method
func (m *MockResourceSKUsPage) Values() []compute.ResourceSku {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Values")
	ret0, _ := ret[0].([]compute.ResourceSku)
	return ret0
}

// Values indicates an expected call of Values
func (mr *MockResourceSKUsPageMockRecorder) Values() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Values", reflect.TypeOf((*MockResourceSKUsPage)(nil).Values))
}
