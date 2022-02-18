package database

import (
	"context"
	"github.com/dkischenko/chat/internal/user"
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/dkischenko/chat/pkg/uuid"
	"sync"
)

type localstorage struct {
	users   []*user.User
	rwMutex *sync.RWMutex
	logger  *logger.Logger
}

func NewStorage(logger *logger.Logger) user.Repository {
	return &localstorage{
		users:   make([]*user.User, 1),
		rwMutex: new(sync.RWMutex),
		logger:  logger,
	}
}

func (ls *localstorage) Create(ctx context.Context, user *user.User) (id string, err error) {
	ls.rwMutex.Lock()
	defer ls.rwMutex.Unlock()
	user.ID = uuid.GetUUID()
	ls.users = append(ls.users, user)

	return user.ID, nil
}

func (ls *localstorage) FindOne(ctx context.Context, username string) (u *user.User, err error) {
	ls.rwMutex.RLock()
	defer ls.rwMutex.RUnlock()

	for _, u := range ls.users {
		if u.Username == username {
			return u, nil
		}
	}

	return nil, user.ErrUserNotFound
}

func (ls *localstorage) FindAll(ctx context.Context) (u []*user.User, err error) {
	if len(ls.users) < 0 {
		ls.logger.Entry.Error(user.ErrUserNotFound)
		return nil, user.ErrUserNotFound
	}
	return ls.users, nil
}

func (ls *localstorage) UpdateKey(ctx context.Context, user *user.User, m map[string]string) (err error) {
	panic("Implement me")
}

func (ls *localstorage) FindByUUID(ctx context.Context, uuid string) (u *user.User, err error) {
	panic("Implement me")
}

func (ls *localstorage) UpdateOnline(ctx context.Context, user *user.User, m map[string]bool) (err error) {
	panic("Implement me")
}

func (ls *localstorage) GetOnline(ctx context.Context) (count int, err error) {
	panic("Implement me")
}
