package tender

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sariya23/tender/internal/domain/models"
	"github.com/sariya23/tender/internal/repository"
	"github.com/sariya23/tender/internal/service"
)

type TenderService struct {
	logger        *slog.Logger
	tenderService service.TenderServiceProvider
}

type GetTendersResponse struct {
	Tenders []models.Tender `json:"tenders,omitempty"`
	Message string          `json:"message"`
}

func New(logger *slog.Logger, tenderService service.TenderServiceProvider) *TenderService {
	return &TenderService{
		logger:        logger,
		tenderService: tenderService,
	}
}

func (s *TenderService) GetTenders(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "internal.api.tender.service.GetTenders"
		logger := s.logger.With("op", op)
		logger.Info(fmt.Sprintf("request to %v", c.Request.URL))

		serviceType := c.DefaultQuery("srv_type", "all")
		tenders, err := s.tenderService.GetTenders(ctx, serviceType)

		if err != nil {
			if errors.Is(err, repository.ErrNoTendersWithThisServiceType) {
				c.JSON(http.StatusBadRequest, GetTendersResponse{Message: fmt.Sprintf("no tenders found with service type: %s", serviceType)})
				return
			} else {
				logger.Error("unexpected error", slog.String("err", err.Error()))
				c.JSON(http.StatusInternalServerError, GetTendersResponse{Message: "internal error"})
				return
			}
		}

		logger.Info("send success response")
		c.JSON(http.StatusOK, GetTendersResponse{Message: "ok", Tenders: tenders})
	}
}
