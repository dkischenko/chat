// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_user is a generated GoMock package.
package mock_user

import (
	context "context"
	"github.com/dkischenko/chat/internal/user/models"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIService is a mock of IService interface.
type MockIService struct {
	ctrl     *gomock.Controller
	recorder *MockIServiceMockRecorder
}

// MockIServiceMockRecorder is the mock recorder for MockIService.
type MockIServiceMockRecorder struct {
	mock *MockIService
}

// NewMockIService creates a new mock instance.
func NewMockIService(ctrl *gomock.Controller) *MockIService {
	mock := &MockIService{ctrl: ctrl}
	mock.recorder = &MockIServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIService) EXPECT() *MockIServiceMockRecorder {
	return m.recorder
}

// ChatStart mocks base method.
func (m *MockIService) ChatStart(ctx context.Context, token string) (*models.User, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChatStart", ctx, token)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChatStart indicates an expected call of ChatStart.
func (mr *MockIServiceMockRecorder) ChatStart(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChatStart", reflect.TypeOf((*MockIService)(nil).ChatStart), ctx, token)
}

// Create mocks base method.
func (m *MockIService) Create(ctx context.Context, user models.UserDTO) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockIServiceMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIService)(nil).Create), ctx, user)
}

// CreateToken mocks base method.
func (m *MockIService) CreateToken(ctx context.Context, u *models.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateToken", ctx, u)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateToken indicates an expected call of CreateToken.
func (mr *MockIServiceMockRecorder) CreateToken(ctx, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateToken", reflect.TypeOf((*MockIService)(nil).CreateToken), ctx, u)
}

// FindByUUID mocks base method.
func (m *MockIService) FindByUUID(ctx context.Context, uid string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUUID", ctx, uid)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUUID indicates an expected call of FindByUUID.
func (mr *MockIServiceMockRecorder) FindByUUID(ctx, uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUUID", reflect.TypeOf((*MockIService)(nil).FindByUUID), ctx, uid)
}

// GetOnlineUsers mocks base method.
func (m *MockIService) GetOnlineUsers(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOnlineUsers", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOnlineUsers indicates an expected call of GetOnlineUsers.
func (mr *MockIServiceMockRecorder) GetOnlineUsers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOnlineUsers", reflect.TypeOf((*MockIService)(nil).GetOnlineUsers), ctx)
}

// InitSocketConnection mocks base method.
func (m *MockIService) InitSocketConnection(w http.ResponseWriter, r *http.Request, u *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitSocketConnection", w, r, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// InitSocketConnection indicates an expected call of InitSocketConnection.
func (mr *MockIServiceMockRecorder) InitSocketConnection(w, r, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitSocketConnection", reflect.TypeOf((*MockIService)(nil).InitSocketConnection), w, r, u)
}

// Login mocks base method.
func (m *MockIService) Login(ctx context.Context, dto *models.UserDTO) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, dto)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockIServiceMockRecorder) Login(ctx, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockIService)(nil).Login), ctx, dto)
}

// RevokeToken mocks base method.
func (m *MockIService) RevokeToken(ctx context.Context, u *models.User) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RevokeToken", ctx, u)
	ret0, _ := ret[0].(bool)
	return ret0
}

// RevokeToken indicates an expected call of RevokeToken.
func (mr *MockIServiceMockRecorder) RevokeToken(ctx, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeToken", reflect.TypeOf((*MockIService)(nil).RevokeToken), ctx, u)
}

// StartWS mocks base method.
func (m *MockIService) StartWS(w http.ResponseWriter, r *http.Request, u *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartWS", w, r, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartWS indicates an expected call of StartWS.
func (mr *MockIServiceMockRecorder) StartWS(w, r, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartWS", reflect.TypeOf((*MockIService)(nil).StartWS), w, r, u)
}
