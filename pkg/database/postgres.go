package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresConnection() (*pgxpool.Pool, error) {

	dbURL := "postgres://utkarshpravind:postgres@localhost:5432/cloudkitchen"

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to PostgreSQL")

	return pool, nil
}
