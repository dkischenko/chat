package user_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/dkischenko/chat/internal/config"
	database "github.com/dkischenko/chat/internal/user/database/postgres"
	"github.com/dkischenko/chat/internal/user/models"
	"github.com/dkischenko/chat/pkg/database/postgres"
	"os"
	"testing"
	"time"

	uerrors "github.com/dkischenko/chat/internal/errors"
	"github.com/dkischenko/chat/internal/user"
	mock_user "github.com/dkischenko/chat/internal/user/mocks"
	"github.com/dkischenko/chat/pkg/auth"
	"github.com/dkischenko/chat/pkg/hasher"
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/dkischenko/chat/pkg/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	l, _ := logger.GetLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_user.NewMockRepository(ctrl)
	assert.NotNil(t, user.NewService(l, mockRepo, 3600), "NewService should be not nil")
}

func TestService_Create(t *testing.T) {
	testCases := []struct {
		name      string
		ctx       context.Context
		user      *models.UserDTO
		wantError bool
	}{
		{
			name: "OK case",
			ctx:  context.Background(),
			user: &models.UserDTO{
				Username: "Bill",
				Password: "password",
			},
			wantError: false,
		},
		{
			name: "Empty password (skip)",
			ctx:  context.Background(),
			user: &models.UserDTO{
				Username: "Bill",
				Password: "",
			},
			wantError: true,
		},
		{
			name: "Empty name",
			ctx:  context.Background(),
			user: &models.UserDTO{
				Username: "",
				Password: "password",
			},
			wantError: true,
		},
	}

	for _, tcase := range testCases {
		t.Run(tcase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			l, _ := logger.GetLogger()

			mockRepo := mock_user.NewMockRepository(ctrl)
			mockRepo.EXPECT().
				Create(tcase.ctx, gomock.Any()).Return(uuid.GetUUID(), nil).AnyTimes()

			service := user.NewService(l, mockRepo, 3600)
			if len(tcase.user.Username) == 0 {
				if tcase.wantError {
					t.Skip("Username can't be empty")
				}
				t.Error("Unexpected error")
			}
			hash, err := hasher.HashPassword(tcase.user.Password)
			if err != nil {
				if tcase.wantError {
					assert.Equal(t, errors.New("String must not be empty"), err)
					t.Skipf("Expected error: %s", err)
				}
				t.Errorf("Unexpected error: %s", err)
			}

			u := &models.User{
				Username:     tcase.user.Username,
				PasswordHash: hash,
			}

			id, err := mockRepo.Create(tcase.ctx, u)
			if err != nil {
				t.Fatalf("Cannot store user due error: %s", err)
			}
			assert.Equal(t, len(id), 36, "Got wrong UUID format")

			usr := *tcase.user
			id, err = service.Create(tcase.ctx, usr)
			if err != nil {
				t.Fatalf("Cannot store user via service due error: %s", err)
			}
			assert.NotNil(t, id, "User id can't be nil")
		})
	}
}

func TestService_CreateUsernameErr(t *testing.T) {
	t.Run("Create User(empty username)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		l, _ := logger.GetLogger()
		ctx := context.Background()
		uDTO := models.UserDTO{
			Username: "",
			Password: "password",
		}

		mockRepo := mock_user.NewMockRepository(ctrl)
		mockRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return("", fmt.Errorf("Error occurs: %w", uerrors.ErrEmptyUsername)).AnyTimes()
		s := user.NewService(l, mockRepo, 3600)
		_, err := s.Create(ctx, uDTO)
		if err != nil {
			assert.ErrorIs(t, err, uerrors.ErrEmptyUsername)
		} else {
			t.Fatalf("Unexpected error: %s", err)
		}
	})
}

func TestService_CreateError(t *testing.T) {
	t.Run("Create User(error with store)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		l, _ := logger.GetLogger()
		ctx := context.Background()
		uDTO := models.UserDTO{
			Username: "Bob",
			Password: "password",
		}

		mockRepo := mock_user.NewMockRepository(ctrl)
		mockRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return("", fmt.Errorf("Error occurs: %w", uerrors.ErrCreateUser))

		s := user.NewService(l, mockRepo, 3600)
		_, err := s.Create(ctx, uDTO)
		if err != nil {
			assert.ErrorIs(t, err, uerrors.ErrCreateUser)
		} else {
			t.Fatalf("Unexpected error: %s", err)
		}
	})
}

func TestService_Login(t *testing.T) {
	t.Run("User login(Ok)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx := context.Background()
		l, _ := logger.GetLogger()
		hash, _ := hasher.HashPassword("password")
		uId := uuid.GetUUID()
		mockRepo := mock_user.NewMockRepository(ctrl)
		uDTO := &models.UserDTO{
			Username: "Bob",
			Password: "password",
		}
		mockRepo.EXPECT().
			FindOneUser(ctx, uDTO.Username).Return(&models.User{
			ID:           uId,
			Username:     uDTO.Username,
			PasswordHash: hash,
			Key:          "",
			IsOnline:     false,
		}, nil).AnyTimes()
		service := user.NewService(l, mockRepo, 3600)
		u, err := mockRepo.FindOneUser(ctx, uDTO.Username)
		if err != nil {
			t.Fatalf("Can't find user with credentials due error: %s", err)
		}

		if !hasher.CheckPasswordHash(u.PasswordHash, uDTO.Password) {
			t.Fatalf("User with wrong password. Error: %s", err)
		}

		usr, err := service.Login(ctx, uDTO)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		assert.NotNil(t, usr)
	})
}

func TestService_LoginFindOneError(t *testing.T) {
	t.Run("User login find one error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctrl.Finish()

		ctx := context.Background()
		mockRepo := mock_user.NewMockRepository(ctrl)
		mockRepo.EXPECT().
			FindOneUser(ctx, "Bob").
			Return(nil, fmt.Errorf("Error occurs: %w", uerrors.ErrFindOneUser)).AnyTimes()

		l, _ := logger.GetLogger()
		s := user.NewService(l, mockRepo, 3600)
		uDTO := &models.UserDTO{
			Username: "Bob",
			Password: "password",
		}
		_, err := s.Login(ctx, uDTO)
		if err != nil {
			assert.ErrorIs(t, err, uerrors.ErrFindOneUser)
		} else {
			t.Fatalf("Unexpected error.")
		}
	})
}

func TestService_FindByUUID(t *testing.T) {
	testCases := []struct {
		name      string
		ctx       context.Context
		uid       string
		wantError bool
	}{
		{
			name:      "Ok",
			ctx:       context.Background(),
			uid:       uuid.GetUUID(),
			wantError: false,
		},
		{
			name:      "Fail",
			ctx:       context.Background(),
			uid:       "",
			wantError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			l, _ := logger.GetLogger()
			mockRepo := mock_user.NewMockRepository(ctrl)
			mockRepo.EXPECT().FindByUUID(test.ctx, gomock.Any()).Return(&models.User{
				ID:           test.uid,
				Username:     "Bob",
				PasswordHash: "$2a$10$gYb3GF0v.o8ycTS3ClIDkOMSjojOazFMVRiwAqh5IXLHhDo4/iQJO",
				Key:          "",
				IsOnline:     false,
			}, nil).AnyTimes()
			service := user.NewService(l, mockRepo, 3600)
			u, err := service.FindByUUID(test.ctx, test.uid)
			if err != nil {
				if test.wantError {
					t.Skipf("Expected error: %s", err)
				}
				t.Errorf("Unexpected error: %s", err)
			}

			assert.NotNil(t, u)
		})
	}
}

func TestService_FindByUUIDError(t *testing.T) {
	t.Run("Error find user by UUID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		l, _ := logger.GetLogger()
		ctx := context.Background()
		mockRepo := mock_user.NewMockRepository(ctrl)
		mockRepo.EXPECT().FindByUUID(ctx, gomock.Any()).
			Return(nil, fmt.Errorf("Error occurs: %w", uerrors.ErrFindUserByUIID))
		s := user.NewService(l, mockRepo, 3600)
		_, err := s.FindByUUID(ctx, uuid.GetUUID())
		if err != nil {
			assert.ErrorIs(t, err, uerrors.ErrFindUserByUIID)
		} else {
			t.Fatalf("Unexpected error.")
		}
	})
}

func TestService_RevokeToken(t *testing.T) {
	testCases := []struct {
		name      string
		ctx       context.Context
		u         *models.User
		wantError bool
	}{
		{
			name: "Ok",
			ctx:  context.Background(),
			u: &models.User{
				ID:           "ce34e740-3f08-4292-964f-6f4ad096f8bc",
				Username:     "Bob",
				PasswordHash: "$2a$10$4nCeWqjuHWH9WtbaQArWkODCbSZNe9kDlmrwGY61dLMdi/5r3/G8K",
				Key:          "Kasdjhagrtewtyewgdeqwyuuyg",
				IsOnline:     false,
			},
			wantError: false,
		},
		{
			name: "Ok",
			ctx:  context.Background(),
			u: &models.User{
				ID:           "ce34e740-3f08-4292-964f-6f4ad096f8bc",
				Username:     "Bob",
				PasswordHash: "$2a$10$4nCeWqjuHWH9WtbaQArWkODCbSZNe9kDlmrwGY61dLMdi/5r3/G8K",
				Key:          "",
				IsOnline:     false,
			},
			wantError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			l, _ := logger.GetLogger()
			mockRepo := mock_user.NewMockRepository(ctrl)
			mockRepo.EXPECT().UpdateKey(test.ctx, test.u, "").Return(nil).AnyTimes()
			service := user.NewService(l, mockRepo, 3600)
			ok := service.RevokeToken(test.ctx, test.u)
			assert.True(t, ok)
		})
	}
}

func TestService_RevokeTokenFalse(t *testing.T) {
	t.Run("Revoke token with error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx := context.Background()
		u := &models.User{
			ID:           "ce34e740-3f08-4292-964f-6f4ad096f8bc",
			Username:     "Bob",
			PasswordHash: "$2a$10$4nCeWqjuHWH9WtbaQArWkODCbSZNe9kDlmrwGY61dLMdi/5r3/G8K",
			Key:          "Kasdjhagrtewtyewgdeqwyuuyg",
			IsOnline:     false,
		}
		mockRepo := mock_user.NewMockRepository(ctrl)
		mockRepo.EXPECT().UpdateKey(ctx, u, "").Return(fmt.Errorf("Error occurs: %w", uerrors.ErrRevokeToken))
		l, _ := logger.GetLogger()
		s := user.NewService(l, mockRepo, 3600)
		ok := s.RevokeToken(ctx, u)
		if !ok {
			assert.False(t, ok)
		} else {
			t.Fatalf("Unexpected error.")
		}
	})
}

func TestService_GetOnlineUsersOk(t *testing.T) {
	t.Run("Get online users (OK)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_user.NewMockRepository(ctrl)
		ctx := context.Background()
		mockRepo.EXPECT().GetOnline(ctx).Return(10, nil)
		l, _ := logger.GetLogger()
		s := user.NewService(l, mockRepo, 3600)
		c, err := s.GetOnlineUsers(ctx)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		assert.Equal(t, c, 10)
	})
}

func TestService_GetOnlineUsersErr(t *testing.T) {
	t.Run("Get online users (Err)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_user.NewMockRepository(ctrl)
		ctx := context.Background()
		mockRepo.EXPECT().GetOnline(ctx).Return(0, fmt.Errorf("Error occurs: %w", uerrors.ErrGetOnlineUsers))
		l, _ := logger.GetLogger()
		s := user.NewService(l, mockRepo, 3600)
		_, err := s.GetOnlineUsers(ctx)
		if err != nil {
			assert.ErrorIs(t, err, uerrors.ErrGetOnlineUsers)
		} else {
			t.Fatalf("Unexpected error: %s", err)
		}
	})
}

func TestService_ChatStartOk(t *testing.T) {
	hash, _ := hasher.HashPassword("password")
	t.Run("[Ok]Chat start", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		httpStatusCode := 200
		usr := &models.User{
			ID:           uuid.GetUUID(),
			Username:     "Bob",
			PasswordHash: hash,
			Key:          "asdasdasdasda",
			IsOnline:     false,
		}

		l, _ := logger.GetLogger()
		mockRepo := mock_user.NewMockRepository(ctrl)
		mockRepo.EXPECT().FindByUUID(ctx, usr.ID).
			Return(usr, nil).AnyTimes()

		mockRepo.EXPECT().
			UpdateKey(ctx, usr, "").
			Return(nil).AnyTimes()

		authRepo, _ := auth.NewManager(3600 * time.Second)
		s := user.NewService(l, mockRepo, 3600*time.Second)
		token, _ := authRepo.CreateJWT(usr.ID)

		_, err := authRepo.ParseJWT(token)
		if err != nil {
			t.Fatalf("unexpected error")
		}
		ok := s.RevokeToken(ctx, usr)
		if err == nil && ok {
			u, code, _ := s.ChatStart(ctx, token)
			assert.NotNil(t, u)
			assert.Equal(t, httpStatusCode, code)
		} else {
			t.Fatalf("Unexpected error.")
		}

	})
}

func TestService_ChatStartParseTokenError(t *testing.T) {
	hash, _ := hasher.HashPassword("password")
	t.Run("[Parse token error] Chat start", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		httpStatusCode := 400
		usr := &models.User{
			ID:           uuid.GetUUID(),
			Username:     "Bob",
			PasswordHash: hash,
			Key:          "asdasdasdasda",
			IsOnline:     false,
		}

		l, _ := logger.GetLogger()
		mockRepo := mock_user.NewMockRepository(ctrl)

		authRepo, _ := auth.NewManager(3600 * time.Second)
		token, _ := authRepo.CreateJWT(usr.ID)
		token = token + "123"

		s := user.NewService(l, mockRepo, 3600*time.Second)
		_, code, err := s.ChatStart(ctx, token)
		if err != nil {
			assert.Equal(t, code, httpStatusCode)
		} else {
			t.Fatalf("Unexpected error.")
		}
	})
}

func TestService_ChatStartFindByUIIDError(t *testing.T) {
	hash, _ := hasher.HashPassword("password")
	t.Run("[FindByUUID error] Chat start", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		httpStatusCode := 400
		usr := &models.User{
			ID:           uuid.GetUUID(),
			Username:     "Bob",
			PasswordHash: hash,
			Key:          "asdasdasdasda",
			IsOnline:     false,
		}

		l, _ := logger.GetLogger()
		mockRepo := mock_user.NewMockRepository(ctrl)
		mockRepo.EXPECT().FindByUUID(ctx, usr.ID).
			Return(nil, fmt.Errorf("")).AnyTimes()

		authRepo, _ := auth.NewManager(3600 * time.Second)
		token, _ := authRepo.CreateJWT(usr.ID)

		s := user.NewService(l, mockRepo, 3600*time.Second)
		_, code, err := s.ChatStart(ctx, token)
		if err != nil {
			assert.Equal(t, code, httpStatusCode)
		} else {
			t.Fatalf("Unexpected error.")
		}
	})
}

func TestService_ChatStartUserKeyEmptyError(t *testing.T) {
	hash, _ := hasher.HashPassword("password")
	t.Run("[User key empty error] Chat start", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		httpStatusCode := 400
		usr := &models.User{
			ID:           uuid.GetUUID(),
			Username:     "Bob",
			PasswordHash: hash,
			Key:          "",
			IsOnline:     false,
		}

		l, _ := logger.GetLogger()
		mockRepo := mock_user.NewMockRepository(ctrl)
		mockRepo.EXPECT().FindByUUID(ctx, usr.ID).
			Return(usr, nil).AnyTimes()

		authRepo, _ := auth.NewManager(3600 * time.Second)
		token, _ := authRepo.CreateJWT(usr.ID)

		s := user.NewService(l, mockRepo, 3600*time.Second)
		_, code, err := s.ChatStart(ctx, token)
		if err != nil {
			assert.Equal(t, code, httpStatusCode)
			assert.ErrorIs(t, err, uerrors.ErrEmptyUserKey)
		} else {
			t.Fatalf("Unexpected error.")
		}
	})
}

func TestService_ChatStartRevokeTokenError(t *testing.T) {
	hash, _ := hasher.HashPassword("password")
	t.Run("[Revoke Token error]Chat start", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		httpStatusCode := 500
		usr := &models.User{
			ID:           uuid.GetUUID(),
			Username:     "Bob",
			PasswordHash: hash,
			Key:          "asdasdasdasda",
			IsOnline:     false,
		}

		l, _ := logger.GetLogger()
		mockRepo := mock_user.NewMockRepository(ctrl)
		mockRepo.EXPECT().FindByUUID(ctx, usr.ID).
			Return(usr, nil).AnyTimes()

		mockRepo.EXPECT().
			UpdateKey(ctx, usr, "").
			Return(fmt.Errorf("Error occurs: %w", uerrors.ErrRevokeToken)).AnyTimes()

		authRepo, _ := auth.NewManager(3600 * time.Second)
		token, _ := authRepo.CreateJWT(usr.ID)

		s := user.NewService(l, mockRepo, 3600*time.Second)
		ok := s.RevokeToken(ctx, usr)
		_, code, err := s.ChatStart(ctx, token)
		if err != nil && !ok {
			assert.Equal(t, httpStatusCode, code)
			assert.ErrorIs(t, err, uerrors.ErrRevokeToken)
		} else {
			t.Fatalf("Unexpected error.")
		}
	})
}

func TestService_CreateTokenOk(t *testing.T) {
	t.Run("[Ok] Create token", func(t *testing.T) {
		hash, _ := hasher.HashPassword("password")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		usr := &models.User{
			ID:           uuid.GetUUID(),
			Username:     "Bob",
			PasswordHash: hash,
			Key:          "asdasdasdasda",
			IsOnline:     false,
		}

		l, _ := logger.GetLogger()
		mockRepo := mock_user.NewMockRepository(ctrl)

		authRepo, _ := auth.NewManager(3600 * time.Second)
		token, _ := authRepo.CreateJWT(usr.ID)

		mockRepo.EXPECT().
			UpdateKey(ctx, usr, token).
			Return(nil).AnyTimes()

		s := user.NewService(l, mockRepo, 3600*time.Second)
		hash, err := s.CreateToken(ctx, usr)
		if err != nil {
			t.Fatalf("Unexpected error.")
		}

		assert.NotNil(t, hash)
	})
}

func TestService_CreateTokenUpdateKeyErr(t *testing.T) {
	t.Run("[Update token error] Create token", func(t *testing.T) {
		hash, _ := hasher.HashPassword("password")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		usr := &models.User{
			ID:           uuid.GetUUID(),
			Username:     "Bob",
			PasswordHash: hash,
			Key:          "asdasdasdasda",
			IsOnline:     false,
		}

		l, _ := logger.GetLogger()
		mockRepo := mock_user.NewMockRepository(ctrl)

		authRepo, _ := auth.NewManager(3600 * time.Second)
		token, _ := authRepo.CreateJWT(usr.ID)

		mockRepo.EXPECT().
			UpdateKey(ctx, usr, token).
			Return(fmt.Errorf("Error occurs: %w", uerrors.ErrUserUpdateKey)).AnyTimes()

		s := user.NewService(l, mockRepo, 3600*time.Second)
		_, err := s.CreateToken(ctx, usr)
		if err != nil {
			assert.ErrorIs(t, err, uerrors.ErrUserUpdateKey)
		} else {
			t.Fatal("Unexpected error.")
		}
	})
}

func BenchmarkService_Login(b *testing.B) {
	var (
		ctx  = context.Background()
		uDTO = &models.UserDTO{
			Username: "bill",
			Password: "password",
		}
	)

	var cfg *config.Config

	configPath := os.Getenv("CONFIG")
	cfg = config.GetConfig(configPath, &config.Config{})

	l, err := logger.GetLogger()
	if err != nil {
		panic(err)
	}

	accessTokenTTL, err := time.ParseDuration(cfg.Auth.AccessTokenTTL)
	if err != nil {
		panic(err)
	}

	client, err := postgres.NewClient(context.Background(), cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Username,
		cfg.Storage.Password, cfg.Storage.Database)
	if err != nil {
		panic(err)
	}
	storage := database.NewStorage(client, l)

	s := user.NewService(l, storage, accessTokenTTL)

	for i := 0; i < b.N; i++ {
		_, err = s.Login(ctx, uDTO)
		if err != nil {
			l.Entry.Fatal(err)
		}
	}
}

func BenchmarkService_FindByUUID(b *testing.B) {
	var (
		ctx = context.Background()
	)

	var cfg *config.Config

	configPath := os.Getenv("CONFIG")
	cfg = config.GetConfig(configPath, &config.Config{})

	l, err := logger.GetLogger()
	if err != nil {
		panic(err)
	}

	accessTokenTTL, err := time.ParseDuration(cfg.Auth.AccessTokenTTL)
	if err != nil {
		panic(err)
	}

	client, err := postgres.NewClient(context.Background(), cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Username,
		cfg.Storage.Password, cfg.Storage.Database)
	if err != nil {
		panic(err)
	}
	storage := database.NewStorage(client, l)

	s := user.NewService(l, storage, accessTokenTTL)

	for i := 0; i < b.N; i++ {
		_, err = s.FindByUUID(ctx, "e103da7d-5336-422c-abb5-a2d2d26c1786")
		if err != nil {
			l.Entry.Fatal(err)
		}
	}
}
