package route

import (
	"context"

	"github.com/gin-gonic/gin"
)

type TenderServicer interface {
	GetTenders(ctx context.Context) gin.HandlerFunc
	GetEmployeeTendersByUsername(ctx context.Context) gin.HandlerFunc
	CreateTender(ctx context.Context) gin.HandlerFunc
	EditTender(ctx context.Context) gin.HandlerFunc
	RollbackTender(ctx context.Context) gin.HandlerFunc
}

func AddTenderRoutes(ctx context.Context, tn TenderServicer, r *gin.RouterGroup) {
	tender := r.Group("/tenders")
	{
		tender.GET("/", tn.GetTenders(ctx))
		tender.GET("/my", tn.GetEmployeeTendersByUsername(ctx))
		tender.POST("/new", tn.CreateTender(ctx))
		tender.PATCH("/:tenderId/edit", tn.EditTender(ctx))
		tender.PUT("/:tenderId/rollback/:version", tn.RollbackTender(ctx))
	}
}
