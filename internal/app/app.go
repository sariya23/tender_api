package app

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	dbapp "github.com/sariya23/tender/internal/app/db"
	serverapp "github.com/sariya23/tender/internal/app/server"
	tenderapp "github.com/sariya23/tender/internal/app/tender"
	"github.com/sariya23/tender/internal/route"
)

type App struct {
	Server *serverapp.ServerApp
}

func New(
	ctx context.Context,
	dbURL string,
	logger *slog.Logger,
	serverAddr string,
	serverPort string,
) *App {
	db := dbapp.New(ctx, dbURL)
	logger.Info("DB init success")
	tender := tenderapp.New(logger, db.Storage, db.Storage, db.Storage, db.Storage)
	logger.Info("tender service init success")

	router := gin.Default()
	apiRouterGroup := router.Group("/api")
	route.AddTenderRoutes(ctx, tender.TenderHandlers, apiRouterGroup)
	route.AddPingRoute(apiRouterGroup)

	serverApp := serverapp.New(serverAddr, serverPort, router)

	return &App{serverApp}
}
