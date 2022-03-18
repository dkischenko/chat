package postgres

import (
	"context"
	"fmt"
	uerrors "github.com/dkischenko/chat/internal/errors"
	"github.com/dkischenko/chat/internal/user"
	"github.com/dkischenko/chat/internal/user/models"
	"github.com/dkischenko/chat/pkg/logger"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
	"time"
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

func (db *postgres) StoreMessage(ctx context.Context, message *models.Message) (id int, err error) {
	q := `
		INSERT INTO messages(text, u_from, created_at)
		VALUES
			($1, $2, $3)
		RETURNING id
	`
	err = db.pool.QueryRow(ctx, q, message.Text, message.UFrom, time.Now().Unix()).
		Scan(&message.ID)

	if err != nil {
		db.logger.Entry.Error(err)
		return 0, fmt.Errorf("Error occurs: %w. %w", err, uerrors.ErrCreateMessage)
	}

	return message.ID, nil
}

func (db *postgres) GetUnreadMessagesCount(ctx context.Context, u *models.User) (count int, err error) {
	q := `
		SELECT count(id)
		FROM messages
		WHERE created_at > $1 AND u_from <> $2
	`
	row := db.pool.QueryRow(ctx, q, u.LastOnline, u.ID)
	err = row.Scan(&count)
	if err != nil {
		db.logger.Entry.Error(err)
		return 0, err
	}

	return count, nil
}

func (db *postgres) GetUnreadMessages(ctx context.Context, u *models.User, unreadMC int) (messages []models.Message, err error) {
	q := `
		SELECT id, text, u_from, created_at
		FROM messages
		WHERE created_at > $1 AND u_from <> $2
	`
	mRows, err := db.pool.Query(ctx, q, u.LastOnline, u.ID)
	if err != nil {
		db.logger.Entry.Error(err)
		return nil, err
	}
	var mSlice = make([]models.Message, 0, unreadMC)
	for mRows.Next() {
		var m models.Message
		err := mRows.Scan(&m.ID, &m.Text, &m.UFrom, &m.CreatedAt)
		if err != nil {
			db.logger.Entry.Error(err)
			return nil, err
		}
		mSlice = append(mSlice, m)
	}

	if err := mRows.Err(); err != nil {
		db.logger.Entry.Error(err)
		return nil, err
	}

	return mSlice, nil
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
		return "", fmt.Errorf("Error occurs: %w. %w", err, uerrors.ErrCreateUser)
	}
	return user.ID, nil
}

func (db *postgres) FindOneUser(ctx context.Context, username string) (u *models.User, err error) {
	u = &models.User{}
	q := `
		SELECT id, username, password_hash, key, is_online, last_online
		FROM users WHERE username = $1
	`
	row := db.pool.QueryRow(ctx, q, username)
	err = row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Key, &u.IsOnline, &u.LastOnline)
	if err != nil {
		db.logger.Entry.Error(err)
		return u, err
	}

	return u, nil
}
func (db *postgres) FindOneMessage(ctx context.Context, mid int) (m *models.Message, err error) {
	m = &models.Message{}
	q := `
		SELECT id, text, u_from, created_at
		FROM messages WHERE id = $1
	`
	row := db.pool.QueryRow(ctx, q, mid)
	err = row.Scan(&m.ID, &m.Text, &m.UFrom, &m.CreatedAt)
	if err != nil {
		db.logger.Entry.Error(err)
		return nil, err
	}

	return m, nil
}

func (db *postgres) FindByUUID(ctx context.Context, uuid string) (u *models.User, err error) {
	u = &models.User{}
	q := `
		SELECT id, username, password_hash, key, is_online, last_online
		FROM users WHERE id = $1
	`
	row := db.pool.QueryRow(ctx, q, uuid)
	err = row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Key, &u.IsOnline, &u.LastOnline)
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
	qOnline := `
		UPDATE users
		SET is_online = $1
		WHERE id = $2
	`

	qOffline := `
		UPDATE users
		SET is_online = $1, last_online = $2
		WHERE id = $3
	`
	if isOnline {
		_, err = db.pool.Exec(ctx, qOnline, isOnline, user.ID)
	} else {
		_, err = db.pool.Exec(ctx, qOffline, isOnline, time.Now().Unix(), user.ID)
	}

	if err != nil {
		db.logger.Entry.Error(err)
		return err
	}
	return nil
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
