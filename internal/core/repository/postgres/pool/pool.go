package core_postgres_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool interface {
	QueryRow(
		ctx context.Context,
		sql string,
		args ...any,
	) pgx.Row
	Query(ctx context.Context,
		sql string,
		args ...any,
	) (pgx.Rows, error)
	Exec(
		ctx context.Context,
		sql string,
		arguments ...any,
	) (pgconn.CommandTag, error)
	Close()
	OpTimeout() time.Duration
}

type ConnectionPool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func NewConnectionPool(ctx context.Context, config Config) (*ConnectionPool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.Database,
	)
	pgxconfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse pgx config: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping pgx pool: %w", err)
	}
	return &ConnectionPool{pool, config.Timeout}, nil
}

func (c *ConnectionPool) OpTimeout() time.Duration {
	return c.opTimeout
}
