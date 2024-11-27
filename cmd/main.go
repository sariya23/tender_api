package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"tender/internal/app"
	"tender/internal/config"
)

func main() {
	ctx := context.Background()

	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	log.Info("starting app at", slog.String("addr", cfg.ServerAddress))

	app := app.New(ctx, log, cfg.ServerAddress, cfg.PostgresConn)
	go app.HttpServer.MustRun()
	log.Info("http server is running", slog.String("addr", cfg.ServerAddress))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")
	app.HttpServer.GracefullStop(ctx)
}
