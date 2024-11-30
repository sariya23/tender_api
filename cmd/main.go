package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sariya23/tender/internal/config"
	"github.com/sariya23/tender/internal/service/tender"
)

func main() {
	// ctx := context.Background()

	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	log.Info("starting app at", slog.String("addr", cfg.ServerAddress))

	tenderSerice := tender.New(log)
	// tenderAPI := tender.New(log, )

	r := gin.Default()
	api := r.Group("/api")
	tender := api.Group("/tenders")
	{
		tender.GET("/")
	}
}
