package routes

import (
	"context"
	"log/slog"
	tn "tender/internal/handlers/tender"

	"github.com/gin-gonic/gin"
)

type TenderProvider interface {
	tn.TenderGetter
	tn.TenderCreater
	tn.UserProvider
	tn.UserResponsibler
	tn.UserTenderGetter
	tn.OrganizationProvider
}

func TenderRoutes(ctx context.Context, logger *slog.Logger, conn TenderProvider, r *gin.RouterGroup) {
	tender := r.Group("/tender")
	{
		tender.GET("/", tn.GetTenders(ctx, logger, conn))
		tender.POST("/new", tn.CreateTender(ctx, logger, conn, conn, conn, conn))
		tender.GET("/my", tn.GetUserTenders(ctx, logger, conn))
	}
}
