package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultDSN = "postgres://exchange:exchange_local@localhost:5432/exchange_dev?sslmode=disable"

// Connect creates a pgx connection pool. If dsn is empty, uses the local Docker default.
func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	if dsn == "" {
		dsn = defaultDSN
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}
	return pool, nil
}
