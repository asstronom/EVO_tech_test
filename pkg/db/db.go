package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type TransactionDB struct {
	pool *pgxpool.Pool
}

func Open(ctx context.Context, url string) (*TransactionDB, error) {
	pool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &TransactionDB{pool: pool}, nil
}

func (db *TransactionDB) Close() {
	db.pool.Close()
}
