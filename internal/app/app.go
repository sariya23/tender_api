package app

import (
	"context"
	"log/slog"
	serverapp "tender/internal/app/server"
	"tender/internal/storage/postgres"
)

type App struct {
	HttpServer *serverapp.ServerApp
	Conn       string
}

func New(ctx context.Context, logger *slog.Logger, addr string, dbURL string) *App {
	conn := postgres.MustNewConnection(ctx, dbURL)
	srv := serverapp.New(ctx, logger, addr, conn)
	return &App{
		HttpServer: srv,
		Conn:       "qwe",
	}
}
