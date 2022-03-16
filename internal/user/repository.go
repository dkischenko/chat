package user

import (
	"context"
	"github.com/dkischenko/chat/internal/user/models"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go
type Repository interface {
	Create(ctx context.Context, user *models.User) (id string, err error)
	FindOne(ctx context.Context, username string) (u *models.User, err error)
	FindByUUID(ctx context.Context, uuid string) (u *models.User, err error)
	UpdateKey(ctx context.Context, user *models.User, key string) (err error)
	UpdateOnline(ctx context.Context, user *models.User, isOnline bool) (err error)
	GetOnline(ctx context.Context) (count int, err error)
}
