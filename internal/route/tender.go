package route

import (
	"context"

	"github.com/gin-gonic/gin"
	tenderapi "github.com/sariya23/tender/internal/api/tender"
)

func AddTenderRoutes(ctx context.Context, tn *tenderapi.TenderService, r *gin.RouterGroup) {
	tender := r.Group("/tenders")
	{
		tender.GET("/", tn.GetTenders(ctx))
		tender.GET("/my", tn.GetEmployeeTendersByUsername(ctx))
		tender.POST("/new", tn.CreateTender(ctx))
		tender.PATCH("/:tenderId/edit", tn.EditTender(ctx))
		tender.PUT("/:tenderId/rollback/:version", tn.RollbackTender(ctx))
	}
}
