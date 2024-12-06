package postgres

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	connection *pgxpool.Pool
}

func MustNewConnection(ctx context.Context, dbURL string) *Storage {
	const op = "storage.postgres.MustNewConnection"
	ctx, cancel := context.WithTimeout(ctx, time.Second*4)
	defer cancel()
	conn, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("%s: cannot connect to db with URL: %s, with error: %v", op, dbURL, err)
	}
	err = conn.Ping(ctx)
	if err != nil {
		log.Fatalf("%s: db is unreachable: %v", op, err)
	}
	return &Storage{connection: conn}
}
