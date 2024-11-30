package postgres

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	connection *pgxpool.Pool
}

func MustNewConnection(dbURL string) *Storage {
	panic("impl me")
}
