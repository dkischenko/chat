package user

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, user *User) (id string, err error)
	FindOne(ctx context.Context, username string) (u *User, err error)
}
