package tender

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"tender/internal/domain/models"
	"tender/internal/storage"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string          `json:"message"`
	Tenders []models.Tedner `json:"tenders"`
}

type TenderGetter interface {
	Tenders(ctx context.Context) ([]models.Tedner, error)
	TendersByServiceType(ctx context.Context, serviceType string) ([]models.Tedner, error)
}

func GetTenders(ctx context.Context, logger *slog.Logger, tenderGetter TenderGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.tender.GetAllTenders"
		logger := logger.With("op", op)
		logger.Info("request to /api/tender/")

		serviceType := c.DefaultQuery("srv_type", "all")
		var err error
		var tenders []models.Tedner

		if serviceType == "all" {
			tenders, err = tenderGetter.Tenders(ctx)
		} else {
			tenders, err = tenderGetter.TendersByServiceType(ctx, serviceType)
		}

		if err != nil {
			if errors.Is(err, storage.ErrNoTenderPresence) {
				logger.Info("no tenders found")
				c.JSON(http.StatusOK, Response{Message: "no tenders found", Tenders: []models.Tedner{}})
			}
			logger.Error("cannot get tenders", slog.String("err", err.Error()))
			c.JSON(http.StatusBadRequest, Response{Message: "internal error", Tenders: []models.Tedner{}})
			return
		}
		logger.Info("success get tenders")
		c.JSON(http.StatusOK, Response{Message: "ok", Tenders: tenders})
	}
}
