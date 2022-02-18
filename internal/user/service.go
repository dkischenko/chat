package user

import (
	"context"
	"github.com/dkischenko/chat/pkg/auth"
	"github.com/dkischenko/chat/pkg/hasher"
	"github.com/dkischenko/chat/pkg/logger"
	"net/http"
	"time"
)

type service struct {
	logger       *logger.Logger
	storage      Repository
	tokenManager *auth.Manager
}

func NewService(logger *logger.Logger, storage Repository, tokenTTL time.Duration) *service {
	tm, err := auth.NewManager(tokenTTL)
	if err != nil {
		logger.Entry.Errorf("error with token manager: %s", err)
	}
	return &service{
		logger:       logger,
		storage:      storage,
		tokenManager: tm,
	}
}

func (s *service) Create(ctx context.Context, user UserDTO) (id string, err error) {
	hashPassword, err := hasher.HashPassword(user.Password)
	if err != nil {
		s.logger.Entry.Errorf("troubles with hashing password: %s", user.Password)
		return "", err
	}
	usr := &User{
		Username:     user.Username,
		PasswordHash: hashPassword,
	}

	id, err = s.storage.Create(ctx, usr)

	if err != nil {
		return id, err
	}

	return
}

func (s *service) Login(ctx context.Context, dto *UserDTO) (u *User, err error) {
	u, err = s.storage.FindOne(ctx, dto.Username)
	if err != nil {
		s.logger.Entry.Errorf("failed find user with error: %s", err)
		return
	}

	if !hasher.CheckPasswordHash(u.PasswordHash, dto.Password) {
		s.logger.Entry.Errorf("user used wrong password: %s", err)
		return
	}

	return
}

func (s *service) FindByUUID(ctx context.Context, uid string) (u *User, err error) {
	u, err = s.storage.FindByUUID(ctx, uid)
	if err != nil {
		s.logger.Entry.Errorf("failed find user with error: %s", err)
		return nil, err
	}
	return
}

func (s *service) RevokeToken(ctx context.Context, u *User) (ok bool) {
	err := s.storage.UpdateKey(ctx, u, map[string]string{"key": ""})
	if err != nil {
		s.logger.Entry.Errorf("issue due error: %s", err)
		return false
	}
	return true
}

func (s *service) CreateToken(ctx context.Context, u *User) (hash string, err error) {
	if err != nil {
		s.logger.Entry.Errorf("issue due error: %s", err)
		return "", err
	}

	hash, err = s.tokenManager.CreateJWT(u.ID)
	if err != nil {
		s.logger.Entry.Errorf("problems with creating jwt token: %s", err)
		return "", err
	}

	if err := s.storage.UpdateKey(ctx, u, map[string]string{"key": hash}); err != nil {
		s.logger.Entry.Errorf("error with user update: %s", err)
	}

	return
}

func (s *service) ParseToken(tokenString string) (uuid string, err error) {
	uuid, err = s.tokenManager.ParseJWT(tokenString)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func (s *service) StartWS(w http.ResponseWriter, r *http.Request, u *User) error {
	err := s.InitSocketConnection(w, r, u)
	if err != nil {
		s.logger.Entry.Errorf("error with websocket initialization: %s", err)
		return err
	}

	return nil
}

func (s *service) GetOnlineUsers(ctx context.Context) (count int, err error) {
	count, err = s.storage.GetOnline(ctx)
	if err != nil {
		s.logger.Entry.Errorf("Error with getting online users: %s", err)
		return 0, err
	}
	return
}
