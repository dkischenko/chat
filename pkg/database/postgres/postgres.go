package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

const protocol = "postgres://"

func NewClient(ctx context.Context, host, port, username, password, database string) (dbpool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("%s%s:%s@%s:%s/%s", protocol, username, password, host, port, database)
	dbpool, err = pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}

	if err = dbpool.Ping(ctx); err != nil {
		return nil, err
	}

	return dbpool, nil
}
