package app

import (
	"log/slog"
	serverapp "tender/internal/app/server"
)

type App struct {
	HttpServer *serverapp.ServerApp
	Conn       string
}

func New(logger *slog.Logger, addr string) *App {
	srv := serverapp.New(logger, addr)
	return &App{
		HttpServer: srv,
		Conn:       "qwe",
	}
}
