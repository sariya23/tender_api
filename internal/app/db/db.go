package dbapp

import (
	"context"

	"github.com/sariya23/tender/internal/repository/postgres"
)

type DBApp struct {
	Storage *postgres.Storage
}

func New(ctx context.Context, dbURL string) *DBApp {
	db := postgres.MustNewConnection(ctx, dbURL)
	return &DBApp{Storage: db}
}
