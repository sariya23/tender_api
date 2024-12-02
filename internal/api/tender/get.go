package tenderapi

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sariya23/tender/internal/api"
	outerror "github.com/sariya23/tender/internal/out_error"
)

func (s *TenderService) GetTenders(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "internal.api.tender.service.GetTenders"
		logger := s.logger.With("op", op)
		logger.Info(fmt.Sprintf("request to %v", c.Request.URL))

		serviceType := c.DefaultQuery("srv_type", "all")
		tenders, err := s.tenderService.GetTenders(ctx, serviceType)
		if err != nil {
			if errors.Is(err, outerror.ErrTendersWithThisServiceTypeNotFound) {
				c.JSON(
					http.StatusBadRequest,
					api.GetTendersResponse{
						Message: fmt.Sprintf(
							"no tenders found with service type: %s",
							serviceType,
						),
					},
				)
				return
			} else {
				logger.Error("unexpected error", slog.String("err", err.Error()))
				c.JSON(http.StatusInternalServerError, api.GetTendersResponse{Message: "internal error"})
				return
			}
		}

		logger.Info("send success response")
		c.JSON(http.StatusOK, api.GetTendersResponse{Message: "ok", Tenders: tenders})
	}
}
