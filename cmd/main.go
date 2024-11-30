package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sariya23/tender/internal/api/tender"
	"github.com/sariya23/tender/internal/config"
	"github.com/sariya23/tender/internal/repository/postgres"
	tendersrv "github.com/sariya23/tender/internal/service/tender"
)

func main() {
	ctx := context.Background()

	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	log.Info("starting app at", slog.String("addr", cfg.ServerAddress))

	db := postgres.MustNewConnection("qwe")
	tenderService := tendersrv.New(log, db, db, db, db)
	tenderAPI := tender.New(log, tenderService)

	r := gin.Default()
	api := r.Group("/api")
	tender := api.Group("/tenders")
	{
		tender.GET("/", tenderAPI.GetTenders(ctx))
	}

	srv := &http.Server{Addr: cfg.ServerAddress, Handler: r}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err.Error())
	}
}
