// Code generated by MockGen. DO NOT EDIT.
// Source: repInterface.go

// Package mock_repInterface is a generated GoMock package.
package mock_repInterface

import (
	reflect "reflect"
	models "services/internal/models"

	gomock "github.com/golang/mock/gomock"
)

// MockRepInterface is a mock of RepInterface interface
type MockRepInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepInterfaceMockRecorder
}

// MockRepInterfaceMockRecorder is the mock recorder for MockRepInterface
type MockRepInterfaceMockRecorder struct {
	mock *MockRepInterface
}

// NewMockRepInterface creates a new mock instance
func NewMockRepInterface(ctrl *gomock.Controller) *MockRepInterface {
	mock := &MockRepInterface{ctrl: ctrl}
	mock.recorder = &MockRepInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepInterface) EXPECT() *MockRepInterfaceMockRecorder {
	return m.recorder
}

// PersonCreate mocks base method
func (m *MockRepInterface) PersonCreate(person models.Person) (models.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PersonCreate", person)
	ret0, _ := ret[0].(models.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PersonCreate indicates an expected call of PersonCreate
func (mr *MockRepInterfaceMockRecorder) PersonCreate(person interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PersonCreate", reflect.TypeOf((*MockRepInterface)(nil).PersonCreate), person)
}

// GetPersonByID mocks base method
func (m *MockRepInterface) GetPersonByID(id int) (models.Person, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPersonByID", id)
	ret0, _ := ret[0].(models.Person)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPersonByID indicates an expected call of GetPersonByID
func (mr *MockRepInterfaceMockRecorder) GetPersonByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPersonByID", reflect.TypeOf((*MockRepInterface)(nil).GetPersonByID), id)
}

// GetAllPersonsInfo mocks base method
func (m *MockRepInterface) GetAllPersonsInfo() ([]models.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllPersonsInfo")
	ret0, _ := ret[0].([]models.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllPersonsInfo indicates an expected call of GetAllPersonsInfo
func (mr *MockRepInterfaceMockRecorder) GetAllPersonsInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllPersonsInfo", reflect.TypeOf((*MockRepInterface)(nil).GetAllPersonsInfo))
}

// UpdatePersonInfo mocks base method
func (m *MockRepInterface) UpdatePersonInfo(person models.Person) (models.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePersonInfo", person)
	ret0, _ := ret[0].(models.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePersonInfo indicates an expected call of UpdatePersonInfo
func (mr *MockRepInterfaceMockRecorder) UpdatePersonInfo(person interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePersonInfo", reflect.TypeOf((*MockRepInterface)(nil).UpdatePersonInfo), person)
}

// DeletePersonInfo mocks base method
func (m *MockRepInterface) DeletePersonInfo(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePersonInfo", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePersonInfo indicates an expected call of DeletePersonInfo
func (mr *MockRepInterfaceMockRecorder) DeletePersonInfo(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePersonInfo", reflect.TypeOf((*MockRepInterface)(nil).DeletePersonInfo), id)
}
