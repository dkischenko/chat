package user

import (
	"context"
	"fmt"
	"github.com/dkischenko/chat/internal/user/models"
	"net/http"
	"sync"
	"time"

	uerrors "github.com/dkischenko/chat/internal/errors"
	"github.com/dkischenko/chat/pkg/auth"
	"github.com/dkischenko/chat/pkg/hasher"
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/gorilla/websocket"
)

type Service struct {
	logger       *logger.Logger
	storage      Repository
	tokenManager *auth.Manager
	Upgrader     websocket.Upgrader
	rwMutex      *sync.RWMutex
	clients      map[*websocket.Conn]bool
}

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type IService interface {
	Create(ctx context.Context, user models.UserDTO) (id string, err error)
	Login(ctx context.Context, dto *models.UserDTO) (u *models.User, err error)
	FindByUUID(ctx context.Context, uid string) (u *models.User, err error)
	RevokeToken(ctx context.Context, u *models.User) (ok bool)
	CreateToken(ctx context.Context, u *models.User) (hash string, err error)
	GetOnlineUsers(ctx context.Context) (count int, err error)
	StartWS(w http.ResponseWriter, r *http.Request, u *models.User) error
	ChatStart(ctx context.Context, token string) (u *models.User, code int, err error)
	InitSocketConnection(w http.ResponseWriter, r *http.Request, u *models.User) error
}

func NewService(logger *logger.Logger, storage Repository, tokenTTL time.Duration) IService {
	tm, err := auth.NewManager(tokenTTL)
	if err != nil {
		logger.Entry.Errorf("error with token manager: %s", err)
	}
	return &Service{
		logger:       logger,
		storage:      storage,
		tokenManager: tm,
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		rwMutex: new(sync.RWMutex),
		clients: make(map[*websocket.Conn]bool),
	}
}

func (s Service) Create(ctx context.Context, user models.UserDTO) (id string, err error) {
	if len(user.Username) == 0 {
		s.logger.Entry.Errorf("error occurs: %s", uerrors.ErrEmptyUsername)
		return "", fmt.Errorf("error occurs: %w", uerrors.ErrEmptyUsername)
	}
	hashPassword, err := hasher.HashPassword(user.Password)
	if err != nil {
		s.logger.Entry.Errorf("troubles with hashing password: %s", user.Password)
		return "", err
	}
	usr := &models.User{
		Username:     user.Username,
		PasswordHash: hashPassword,
	}

	id, err = s.storage.Create(ctx, usr)

	if err != nil {
		return id, err
	}

	return
}

func (s Service) Login(ctx context.Context, dto *models.UserDTO) (u *models.User, err error) {
	u, err = s.storage.FindOne(ctx, dto.Username)
	if err != nil {
		s.logger.Entry.Errorf("failed find user with error: %s", err)
		return nil, fmt.Errorf("error occurs: %w", uerrors.ErrFindOneUser)
	}

	if !hasher.CheckPasswordHash(u.PasswordHash, dto.Password) {
		s.logger.Entry.Errorf("user used wrong password: %s", err)
		return nil, fmt.Errorf("error occurs: %w", uerrors.ErrCheckUserPasswordHash)
	}

	return
}

func (s Service) FindByUUID(ctx context.Context, uid string) (u *models.User, err error) {
	u, err = s.storage.FindByUUID(ctx, uid)
	if err != nil {
		s.logger.Entry.Errorf("failed find user with error: %s", err)
		return nil, fmt.Errorf("error occurs: %w with uuid %s", uerrors.ErrFindUserByUIID, uid)
	}
	return
}

func (s Service) RevokeToken(ctx context.Context, u *models.User) (ok bool) {
	err := s.storage.UpdateKey(ctx, u, "")
	if err != nil {
		s.logger.Entry.Errorf("error occurs: %s. %w", err, uerrors.ErrRevokeToken)
		return false
	}
	return true
}

func (s Service) CreateToken(ctx context.Context, u *models.User) (hash string, err error) {
	hash, err = s.tokenManager.CreateJWT(u.ID)
	if err != nil {
		s.logger.Entry.Errorf("problems with creating jwt token: %s", err)
		return "", fmt.Errorf("error occurs: %w", uerrors.ErrCreateJWTToken)
	}

	if err := s.storage.UpdateKey(ctx, u, hash); err != nil {
		s.logger.Entry.Errorf("error with user update: %s", err)
		return "", fmt.Errorf("error occurs: %w", uerrors.ErrUserUpdateKey)
	}

	return
}

func (s Service) parseToken(tokenString string) (uuid string, err error) {
	uuid, err = s.tokenManager.ParseJWT(tokenString)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func (s *Service) StartWS(w http.ResponseWriter, r *http.Request, u *models.User) error {
	err := s.InitSocketConnection(w, r, u)
	if err != nil {
		s.logger.Entry.Errorf("error with websocket initialization: %s", err)
		return err
	}

	return nil
}

func (s Service) GetOnlineUsers(ctx context.Context) (count int, err error) {
	count, err = s.storage.GetOnline(ctx)
	if err != nil {
		s.logger.Entry.Errorf("error occurs: %s. %w", err, uerrors.ErrGetOnlineUsers)
		return 0, err
	}
	return
}

func (s Service) ChatStart(ctx context.Context, token string) (u *models.User, code int, err error) {
	uuid, err := s.parseToken(token)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	u, err = s.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	if len(u.Key) == 0 {
		return u, http.StatusBadRequest, fmt.Errorf("error occurs: %w", uerrors.ErrEmptyUserKey)
	}
	ok := s.RevokeToken(ctx, u)
	if !ok {
		return u, http.StatusInternalServerError, fmt.Errorf("error occurs: %w", uerrors.ErrRevokeToken)
	}

	return u, http.StatusOK, nil
}
