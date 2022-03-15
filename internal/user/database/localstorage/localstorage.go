package database

import (
	"context"
	uerrors "github.com/dkischenko/chat/internal/errors"
	"github.com/dkischenko/chat/internal/user"
	"github.com/dkischenko/chat/internal/user/models"
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/dkischenko/chat/pkg/uuid"
	"sync"
)

type localstorage struct {
	users   []*models.User
	rwMutex *sync.RWMutex
	logger  *logger.Logger
}

func NewStorage(logger *logger.Logger) user.Repository {
	return &localstorage{
		users:   make([]*models.User, 1),
		rwMutex: new(sync.RWMutex),
		logger:  logger,
	}
}

func (ls *localstorage) Create(ctx context.Context, user *models.User) (id string, err error) {
	ls.rwMutex.Lock()
	defer ls.rwMutex.Unlock()
	user.ID = uuid.GetUUID()
	ls.users = append(ls.users, user)

	return user.ID, nil
}

func (ls *localstorage) FindOne(ctx context.Context, username string) (u *models.User, err error) {
	ls.rwMutex.RLock()
	defer ls.rwMutex.RUnlock()

	for _, u := range ls.users {
		if u.Username == username {
			return u, nil
		}
	}

	return nil, uerrors.ErrUserNotFound
}

func (ls *localstorage) FindAll(ctx context.Context) (u []*models.User, err error) {
	if len(ls.users) < 0 {
		ls.logger.Entry.Error(uerrors.ErrUserNotFound)
		return nil, uerrors.ErrUserNotFound
	}
	return ls.users, nil
}

func (ls *localstorage) UpdateKey(ctx context.Context, user *models.User, key string) (err error) {
	panic("Implement me")
}

func (ls *localstorage) FindByUUID(ctx context.Context, uuid string) (u *models.User, err error) {
	panic("Implement me")
}

func (ls *localstorage) UpdateOnline(ctx context.Context, user *models.User, isOnline bool) (err error) {
	panic("Implement me")
}

func (ls *localstorage) GetOnline(ctx context.Context) (count int, err error) {
	panic("Implement me")
}
