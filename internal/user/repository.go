package user

import (
	"context"
	"github.com/dkischenko/chat/internal/user/models"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go
type Repository interface {
	Create(ctx context.Context, user *models.User) (id string, err error)
	StoreMessage(ctx context.Context, message *models.Message) (id int, err error)
	GetUnreadMessages(ctx context.Context, u *models.User, unreadMC int) (messages []models.Message, err error)
	GetUnreadMessagesCount(ctx context.Context, u *models.User) (count int, err error)
	FindOneUser(ctx context.Context, username string) (u *models.User, err error)
	FindOneMessage(ctx context.Context, mid int) (m *models.Message, err error)
	FindByUUID(ctx context.Context, uuid string) (u *models.User, err error)
	UpdateKey(ctx context.Context, user *models.User, key string) (err error)
	UpdateOnline(ctx context.Context, user *models.User, isOnline bool) (err error)
	GetOnline(ctx context.Context) (count int, err error)
}
