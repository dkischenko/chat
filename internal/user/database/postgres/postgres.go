package postgres

import (
	"context"
	"github.com/dkischenko/chat/internal/user"
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
)

type postgres struct {
	logger *logger.Logger
	pool   *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool, logger *logger.Logger) user.Repository {
	return &postgres{
		logger: logger,
		pool:   pool,
	}
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", "")
}

func (db *postgres) Create(ctx context.Context, user *user.User) (id string, err error) {
	q := `
		INSERT INTO users(username, password_hash) 
		    VALUES
		           ($1, $2)
		    RETURNING id
	`
	err = db.pool.QueryRow(ctx, formatQuery(q), user.Username, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		db.logger.Entry.Error(err)
		return "", err
	}
	return user.ID, nil
}

func (db *postgres) FindOne(ctx context.Context, username string) (u *user.User, err error) {
	u = &user.User{}
	q := `
		SELECT id, username, password_hash FROM users WHERE username = $1
	`
	row := db.pool.QueryRow(ctx, q, username)
	err = row.Scan(&u.ID, &u.Username, &u.PasswordHash)
	if err != nil {
		db.logger.Entry.Error(err)
		return u, err
	}

	return u, nil
}
