package routes

import (
	"context"
	"log/slog"
	tn "tender/internal/handlers/tender"

	"github.com/gin-gonic/gin"
)

type TenderProveder interface {
	tn.TenderGetter
}

func TenderRoutes(ctx context.Context, logger *slog.Logger, conn TenderProveder, r *gin.RouterGroup) {
	tender := r.Group("/tender")
	{
		tender.GET("/")
		tender.POST("/new", tn.GetTenders(ctx, logger, conn))
	}
}
