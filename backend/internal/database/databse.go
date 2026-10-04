package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, databaseUrl string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return nil, err
	}

	return db, nil
}