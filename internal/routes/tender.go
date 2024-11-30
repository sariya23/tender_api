package routes

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func TenderRoutes(ctx context.Context, logger *slog.Logger, conn TenderProvider, r *gin.RouterGroup) {
	tender := r.Group("/tender")
	{
		tender.GET("/", tn.GetTenders(ctx, logger, conn))
		tender.POST("/new", tn.CreateTender(ctx, logger, conn, conn, conn, conn))
		tender.GET("/my", tn.GetUserTenders(ctx, logger, conn))
		tender.PATCH("/:tender_id/edit", tn.EditTender(ctx, logger, conn))
	}
}
