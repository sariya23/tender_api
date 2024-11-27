package main

import (
	"log/slog"
	"os"
	"tender/internal/config"
)

func main() {
	cfg := config.MustLoad()
	// TODO: init config
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	log.Info("starting app", slog.String("addr", cfg.ServerAddress))
	// TODO: init app
	// TODO: run app
}
