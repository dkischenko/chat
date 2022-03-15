package postgres

import (
	"context"
	"fmt"
	"github.com/dkischenko/chat/internal/user/models"
	"strings"

	uerrors "github.com/dkischenko/chat/internal/errors"
	"github.com/dkischenko/chat/internal/user"
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/jackc/pgx/v4/pgxpool"
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

func (db *postgres) Create(ctx context.Context, user *models.User) (id string, err error) {
	q := `
		INSERT INTO users(username, password_hash) 
		    VALUES
		           ($1, $2)
		    RETURNING id
	`
	err = db.pool.QueryRow(ctx, formatQuery(q), user.Username, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		db.logger.Entry.Error(err)
		return "", fmt.Errorf("Error occurs: %s. %s", err, uerrors.ErrCreateUser)
	}
	return user.ID, nil
}

func (db *postgres) FindOne(ctx context.Context, username string) (u *models.User, err error) {
	u = &models.User{}
	q := `
		SELECT id, username, password_hash, key, is_online 
		FROM users WHERE username = $1
	`
	row := db.pool.QueryRow(ctx, q, username)
	err = row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Key, &u.IsOnline)
	if err != nil {
		db.logger.Entry.Error(err)
		return u, err
	}

	return u, nil
}

func (db *postgres) FindByUUID(ctx context.Context, uuid string) (u *models.User, err error) {
	u = &models.User{}
	q := `
		SELECT id, username, password_hash, key, is_online
		FROM users WHERE id = $1
	`
	row := db.pool.QueryRow(ctx, q, uuid)
	err = row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Key, &u.IsOnline)
	if err != nil {
		db.logger.Entry.Error(err)
		return
	}

	return
}

func (db *postgres) UpdateKey(ctx context.Context, user *models.User, key string) (err error) {
	q := `
		UPDATE users
		SET key = $1
		WHERE id = $2
	`
	_, err = db.pool.Exec(ctx, q, key, user.ID)
	if err != nil {
		db.logger.Entry.Error(err)
		return fmt.Errorf("Error occurs: %w", uerrors.ErrUserUpdateKey)
	}
	return
}

func (db *postgres) UpdateOnline(ctx context.Context, user *models.User, isOnline bool) (err error) {
	q := `
		UPDATE users
		SET is_online = $1
		WHERE id = $2
	`
	_, err = db.pool.Exec(ctx, q, isOnline, user.ID)
	if err != nil {
		db.logger.Entry.Error(err)
		return nil
	}
	return
}

func (db *postgres) GetOnline(ctx context.Context) (count int, err error) {
	q := `
		SELECT count(id) 
		FROM users WHERE is_online = true
	`
	row := db.pool.QueryRow(ctx, q)
	err = row.Scan(&count)
	if err != nil {
		db.logger.Entry.Error(err)
		return 0, err
	}

	return count, nil
}
