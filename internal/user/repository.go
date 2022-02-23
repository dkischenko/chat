package user

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, user *User) (id string, err error)
	FindOne(ctx context.Context, username string) (u *User, err error)
	FindByUUID(ctx context.Context, uuid string) (u *User, err error)
	UpdateKey(ctx context.Context, user *User, key string) (err error)
	UpdateOnline(ctx context.Context, user *User, isOnline bool) (err error)
	GetOnline(ctx context.Context) (count int, err error)
}
