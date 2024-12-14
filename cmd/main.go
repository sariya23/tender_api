package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sariya23/tender/internal/app"
	"github.com/sariya23/tender/internal/config"
)

func main() {
	ctx := context.Background()

	cfg := config.MustLoad()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("starting app at", slog.String("addr", cfg.ServerAddress))
	app := app.New(ctx, cfg.PostgresConn, logger, cfg.ServerAddress)
	logger.Info("app init success")
	go app.Server.MustRun()

	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-quitSignal
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := app.Server.Server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		logger.Info("timeout of 5 seconds.")
	}
	logger.Info("Server exiting")
}
