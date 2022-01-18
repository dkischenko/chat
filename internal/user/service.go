package user

import (
	"context"
	"github.com/dkischenko/chat/pkg/hasher"
	"github.com/dkischenko/chat/pkg/logger"
)

type service struct {
	logger  *logger.Logger
	storage Repository
}

func NewService(logger *logger.Logger, storage Repository) *service {
	return &service{
		logger:  logger,
		storage: storage,
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

func (s *service) Login(ctx context.Context, username string) (string, error) {
	u, err := s.storage.FindOne(ctx, username)
	if err != nil {
		s.logger.Entry.Errorf("failed find user with error: %s", err)
	}

	// @todo: store user to context

	hash, err := hasher.HashPassword(u.Username + u.ID)
	if err != nil {
		s.logger.Entry.Errorf("problems with hashing user data: %s", err)
	}

	return hash, nil
}
