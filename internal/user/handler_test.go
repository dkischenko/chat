package user_test

import (
	"context"
	"fmt"
	"github.com/dkischenko/chat/internal/config"
	"github.com/dkischenko/chat/internal/user"
	mock_user "github.com/dkischenko/chat/internal/user/mocks"
	"github.com/dkischenko/chat/pkg/auth"
	"github.com/dkischenko/chat/pkg/hasher"
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/dkischenko/chat/pkg/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandler_RegisterOk(t *testing.T) {
	t.Run("[Ok] Register handlers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			ctx  = context.Background()
			uDTO = user.UserDTO{
				Username: "bill",
				Password: "password",
			}
			getUUID = uuid.GetUUID()
			cfg     = &config.Config{}
			payload = `
				{
					"userName": "bill",
					"password": "password"
				}`
		)

		req := httptest.NewRequest(http.MethodPost, "/v1/user", strings.NewReader(payload))
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()
		mockService := mock_user.NewMockIService(ctrl)
		mockService.EXPECT().Create(ctx, uDTO).Return(getUUID, nil).AnyTimes()
		h := user.NewHandler(l, mockService, cfg)
		router := http.NewServeMux()
		h.Register(router)
		h.CreateUser(w, req)
		assert.Equal(t, w.Code, http.StatusOK)
	})
}

func TestHandler_CreateUserOk(t *testing.T) {
	t.Run("[Ok] Create user handler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			ctx  = context.Background()
			uDTO = user.UserDTO{
				Username: "bill",
				Password: "password",
			}
			getUUID = uuid.GetUUID()
			cfg     = &config.Config{}
			payload = `
				{
					"userName": "bill",
					"password": "password"
				}`
		)

		req := httptest.NewRequest(http.MethodPost, "/v1/user", strings.NewReader(payload))
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()
		mockService := mock_user.NewMockIService(ctrl)
		mockService.EXPECT().Create(ctx, uDTO).Return(getUUID, nil).AnyTimes()
		h := user.NewHandler(l, mockService, cfg)
		h.CreateUser(w, req)
		assert.Equal(t, w.Code, http.StatusOK)
	})
}

func TestHandler_CreateUserWrongHttpMethod(t *testing.T) {
	t.Run("[Wrong http method] Create user handler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			ctx  = context.Background()
			uDTO = user.UserDTO{
				Username: "bill",
				Password: "password",
			}
			getUUID = uuid.GetUUID()
			cfg     = &config.Config{}
			payload = `
				{
					"userName": "bill",
					"password": "password"
				}`
		)

		req := httptest.NewRequest(http.MethodGet, "/v1/user", strings.NewReader(payload))
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()
		mockService := mock_user.NewMockIService(ctrl)
		mockService.EXPECT().Create(ctx, uDTO).Return(getUUID, nil).AnyTimes()
		h := user.NewHandler(l, mockService, cfg)
		h.CreateUser(w, req)
		assert.Equal(t, w.Code, http.StatusInternalServerError)
	})
}

func TestHandler_CreateUserWrongJson(t *testing.T) {
	t.Run("[Wrong Json] Create user handler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			ctx  = context.Background()
			uDTO = user.UserDTO{
				Username: "bill",
				Password: "password",
			}
			getUUID = uuid.GetUUID()
			cfg     = &config.Config{}
			payload = `
				{
					"userName": "bill",
					"password": "password
				}`
		)

		req := httptest.NewRequest(http.MethodPost, "/v1/user", strings.NewReader(payload))
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()
		mockService := mock_user.NewMockIService(ctrl)
		mockService.EXPECT().Create(ctx, uDTO).Return(getUUID, nil).AnyTimes()
		h := user.NewHandler(l, mockService, cfg)
		h.CreateUser(w, req)
		assert.Equal(t, w.Code, http.StatusInternalServerError)
	})
}

func TestHandler_CreateUserWrongUserData(t *testing.T) {
	t.Run("[Wrong Json] Create user handler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			ctx  = context.Background()
			uDTO = user.UserDTO{
				Username: "bill",
				Password: "password",
			}
			getUUID = uuid.GetUUID()
			cfg     = &config.Config{}
			payload = `
				{
					"userName": "bill1",
					"password": "password"
				}`
		)

		req := httptest.NewRequest(http.MethodPost, "/v1/user", strings.NewReader(payload))
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()
		mockService := mock_user.NewMockIService(ctrl)
		mockService.EXPECT().Create(ctx, uDTO).Return(getUUID, nil).AnyTimes()
		h := user.NewHandler(l, mockService, cfg)
		h.CreateUser(w, req)
		assert.Equal(t, w.Code, http.StatusBadRequest)
	})
}

func TestHandler_LoginUserOk(t *testing.T) {
	t.Run("[Ok] Login user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		hash, _ := hasher.HashPassword("password")
		var (
			ctx  = context.Background()
			uDTO = &user.UserDTO{
				Username: "bill",
				Password: "password",
			}
			getUUID = uuid.GetUUID()
			u       = &user.User{
				ID:           getUUID,
				Username:     "bill",
				PasswordHash: hash,
				Key:          "",
				IsOnline:     false,
			}
			cfg     = &config.Config{}
			payload = `
				{
					"userName": "bill",
					"password": "password"
				}`
		)
		req := httptest.NewRequest(http.MethodPost, "/v1/user/login", strings.NewReader(payload))
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()
		mockService := mock_user.NewMockIService(ctrl)
		mockService.EXPECT().Login(ctx, uDTO).Return(u, nil).AnyTimes()
		mockService.EXPECT().CreateToken(ctx, u).Return(hash, nil).AnyTimes()
		h := user.NewHandler(l, mockService, cfg)
		h.LoginUser(w, req)
		assert.Equal(t, w.Code, http.StatusOK)
	})
}

func TestHandler_LoginUserWrongHttpMethod(t *testing.T) {
	t.Run("[Wrong http method] Login user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		hash, _ := hasher.HashPassword("password")
		var (
			ctx  = context.Background()
			uDTO = &user.UserDTO{
				Username: "bill",
				Password: "password",
			}
			getUUID = uuid.GetUUID()
			u       = &user.User{
				ID:           getUUID,
				Username:     "bill",
				PasswordHash: hash,
				Key:          "",
				IsOnline:     false,
			}
			cfg     = &config.Config{}
			payload = `
				{
					"userName": "bill",
					"password": "password"
				}`
		)
		req := httptest.NewRequest(http.MethodGet, "/v1/user/login", strings.NewReader(payload))
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()
		mockService := mock_user.NewMockIService(ctrl)
		mockService.EXPECT().Login(ctx, uDTO).Return(u, nil).AnyTimes()
		mockService.EXPECT().CreateToken(ctx, u).Return(hash, nil).AnyTimes()
		h := user.NewHandler(l, mockService, cfg)
		h.LoginUser(w, req)
		assert.Equal(t, w.Code, http.StatusInternalServerError)
	})
}

func TestHandler_LoginWrongJson(t *testing.T) {
	t.Run("[Wrong Json] Login user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		hash, _ := hasher.HashPassword("password")
		var (
			ctx  = context.Background()
			uDTO = &user.UserDTO{
				Username: "bill",
				Password: "password",
			}
			getUUID = uuid.GetUUID()
			u       = &user.User{
				ID:           getUUID,
				Username:     "bill",
				PasswordHash: hash,
				Key:          "",
				IsOnline:     false,
			}
			cfg     = &config.Config{}
			payload = `
				{
					"userName": "bill",
					"password": "password
				}`
		)
		req := httptest.NewRequest(http.MethodPost, "/v1/user/login", strings.NewReader(payload))
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()
		mockService := mock_user.NewMockIService(ctrl)
		mockService.EXPECT().Login(ctx, uDTO).Return(u, nil).AnyTimes()
		mockService.EXPECT().CreateToken(ctx, u).Return(hash, nil).AnyTimes()
		h := user.NewHandler(l, mockService, cfg)
		h.LoginUser(w, req)
		assert.Equal(t, w.Code, http.StatusBadRequest)
	})
}

func TestHandler_LoginUserWrongUserData(t *testing.T) {
	t.Run("[Wrong user data] Login user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		hash, _ := hasher.HashPassword("password")
		var (
			ctx  = context.Background()
			uDTO = &user.UserDTO{
				Username: "bill",
				Password: "password",
			}
			getUUID = uuid.GetUUID()
			u       = &user.User{
				ID:           getUUID,
				Username:     "bill",
				PasswordHash: hash,
				Key:          "",
				IsOnline:     false,
			}
			cfg     = &config.Config{}
			payload = `
				{
					"userName": "bill1",
					"password": "password1"
				}`
		)
		req := httptest.NewRequest(http.MethodPost, "/v1/user/login", strings.NewReader(payload))
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()
		mockService := mock_user.NewMockIService(ctrl)
		mockService.EXPECT().Login(ctx, uDTO).Return(u, nil).AnyTimes()
		mockService.EXPECT().CreateToken(ctx, u).Return(hash, nil).AnyTimes()
		h := user.NewHandler(l, mockService, cfg)
		h.LoginUser(w, req)
		assert.Equal(t, w.Code, http.StatusBadRequest)
	})
}

func TestHandler_ActiveUserOk(t *testing.T) {
	t.Run("[Ok] Active user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			ctx = context.Background()
			cfg = &config.Config{}
		)

		req := httptest.NewRequest(http.MethodGet, "/v1/user/active", nil)
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()
		mockService := mock_user.NewMockIService(ctrl)
		mockService.EXPECT().GetOnlineUsers(ctx).Return(10, nil).AnyTimes()
		h := user.NewHandler(l, mockService, cfg)
		h.ActiveUser(w, req)
		assert.Equal(t, w.Code, http.StatusOK)
	})
}

func TestHandler_ActiveUserWrongHttpMethod(t *testing.T) {
	t.Run("[Ok] Active user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			ctx = context.Background()
			cfg = &config.Config{}
		)

		req := httptest.NewRequest(http.MethodPost, "/v1/user/active", nil)
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()
		mockService := mock_user.NewMockIService(ctrl)
		mockService.EXPECT().GetOnlineUsers(ctx).Return(10, nil).AnyTimes()
		h := user.NewHandler(l, mockService, cfg)
		h.ActiveUser(w, req)
		assert.Equal(t, w.Code, http.StatusInternalServerError)
	})
}

func TestHandler_ActiveUserErrorGettingUsers(t *testing.T) {
	t.Run("[Ok] Active user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			ctx = context.Background()
			cfg = &config.Config{}
		)

		req := httptest.NewRequest(http.MethodGet, "/v1/user/active", nil)
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()
		mockService := mock_user.NewMockIService(ctrl)
		mockService.EXPECT().GetOnlineUsers(ctx).Return(0, fmt.Errorf("")).AnyTimes()
		h := user.NewHandler(l, mockService, cfg)
		h.ActiveUser(w, req)
		assert.Equal(t, w.Code, http.StatusInternalServerError)
	})
}

func TestHandler_ChatStartOk(t *testing.T) {
	t.Run("[OK] Chat start", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		hash, _ := hasher.HashPassword("password")
		uid := uuid.GetUUID()
		var (
			ctx = context.Background()
			u   = &user.User{
				ID:           uid,
				Username:     "bill",
				PasswordHash: hash,
				Key:          "",
				IsOnline:     false,
			}
			cfg = &config.Config{}
		)

		m, _ := auth.NewManager(3600 * time.Second)
		jwt, _ := m.CreateJWT(uid)
		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("/v1/chat/ws.rtm.start/?token=%s", jwt), strings.NewReader(jwt))
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()

		mockService := mock_user.NewMockIService(ctrl)
		mockService.EXPECT().
			ChatStart(ctx, jwt).
			Return(u, http.StatusOK, nil)
		mockService.EXPECT().
			StartWS(w, req, u).
			Return(nil)
		handler := user.NewHandler(l, mockService, cfg)
		handler.ChatStart(w, req)
		assert.Equal(t, w.Code, http.StatusOK)
	})
}

func TestHandler_ChatStartWrongHttpMethod(t *testing.T) {
	t.Run("[OK] Wrong http method", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		uid := uuid.GetUUID()
		var (
			cfg = &config.Config{}
		)

		m, _ := auth.NewManager(3600 * time.Second)
		jwt, _ := m.CreateJWT(uid)
		req := httptest.NewRequest(http.MethodPost,
			fmt.Sprintf("/v1/chat/ws.rtm.start/?token=%s", jwt), strings.NewReader(jwt))
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()

		mockService := mock_user.NewMockIService(ctrl)
		handler := user.NewHandler(l, mockService, cfg)
		handler.ChatStart(w, req)
		assert.Equal(t, w.Code, http.StatusInternalServerError)
	})
}

func TestHandler_ChatStartEmptyToken(t *testing.T) {
	t.Run("[Empty Token] Chat start", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		uid := uuid.GetUUID()
		cfg := &config.Config{}

		m, _ := auth.NewManager(3600 * time.Second)
		jwt, _ := m.CreateJWT(uid)
		req := httptest.NewRequest(http.MethodGet,
			"/v1/chat/ws.rtm.start/?token=", strings.NewReader(jwt))
		w := httptest.NewRecorder()
		l, _ := logger.GetLogger()

		mockService := mock_user.NewMockIService(ctrl)
		handler := user.NewHandler(l, mockService, cfg)
		handler.ChatStart(w, req)
		assert.Equal(t, w.Code, http.StatusBadRequest)
	})
}
