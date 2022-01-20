package postgres

import (
	"context"
	"github.com/dkischenko/chat/internal/user"
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/jackc/pgx/v4"
)

type postgres struct {
	user.Repository
	logger     *logger.Logger
	connection *pgx.Conn
}

func NewStorage(conn *pgx.Conn, logger *logger.Logger) user.Repository {
	return &postgres{
		logger:     logger,
		connection: conn,
	}
}

func (db *postgres) Create(ctx context.Context, user *user.User) (id string, err error) {
	return "", nil
}

func (db *postgres) FindOne(ctx context.Context, username string) (u *user.User, err error) {
	panic("findone not implemented")
}

func (db *postgres) FindAll(ctx context.Context) (u []*user.User, err error) {
	panic("findall not implemented")
}
