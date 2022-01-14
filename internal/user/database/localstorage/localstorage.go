package database

import (
	"context"
	"github.com/dkischenko/chat/internal/user"
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/dkischenko/chat/pkg/uuid"
	"sync"
)

type localstorage struct {
	user.Repository
	users  []*user.User
	mutex  *sync.Mutex
	logger *logger.Logger
}

func NewStorage(logger *logger.Logger) user.Repository {
	return &localstorage{
		users:  make([]*user.User, 1),
		mutex:  new(sync.Mutex),
		logger: logger,
	}
}

func (ls *localstorage) Create(ctx context.Context, user *user.User) (id string, err error) {
	ls.mutex.Lock()
	user.ID = uuid.GetUUID()
	ls.users = append(ls.users, user)
	defer ls.mutex.Unlock()

	return user.ID, nil
}

func (ls *localstorage) FindOne(ctx context.Context, username string) (u *user.User, err error) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	for _, u := range ls.users {
		if u.Username == username {
			return u, nil
		}
	}

	return nil, user.ErrUserNotFound
}

func (ls *localstorage) FindAll(ctx context.Context) (u []*user.User, err error) {
	if len(ls.users) < 0 {
		ls.logger.Fatal(user.ErrUserNotFound)
		return nil, user.ErrUserNotFound
	}
	return ls.users, nil
}
