package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	tenderapi "github.com/sariya23/tender/internal/api/tender"
	"github.com/sariya23/tender/internal/config"
	"github.com/sariya23/tender/internal/repository/postgres"
	tendersrv "github.com/sariya23/tender/internal/service/tender"
)

func main() {
	ctx := context.Background()

	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	log.Info("starting app at", slog.String("addr", cfg.ServerAddress))
	db := postgres.MustNewConnection(ctx, cfg.PostgresConn)
	log.Info("db connect success")
	tenderService := tendersrv.New(log, db, db, db, db)
	tenderAPI := tenderapi.New(log, tenderService)

	r := gin.Default()
	api := r.Group("/api")
	ping := api.Group("/ping")
	{
		ping.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, struct{ Message string }{Message: "ok"}) })
	}
	tender := api.Group("/tenders")
	{
		tender.GET("/", tenderAPI.GetTenders(ctx))
		tender.GET("/my", tenderAPI.GetEmployeeTendersByUsername(ctx))
		tender.POST("/new", tenderAPI.CreateTender(ctx))
		tender.PATCH("/:tenderId/edit", tenderAPI.EditTedner(ctx))
		tender.PUT("/:tenderId/rollback/:version", tenderAPI.RollbackTender(ctx))
	}

	srv := &http.Server{Addr: cfg.ServerAddress, Handler: r}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err.Error())
	}
}
