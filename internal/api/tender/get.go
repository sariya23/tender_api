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
		const op = "internal.api.tenderapi.GetTenders"
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

func (s *TenderService) GetEmployeeTendersByUsername(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "internal.api.tenderapi.GetEmployeeTendersByUsername"
		logger := s.logger.With("op", op)
		logger.Info(fmt.Sprintf("request to %s", c.Request.URL))

		username := c.Query("username")
		if username == "" {
			logger.Info("username not specified. Rediredr to /api/tenders")
			c.Redirect(http.StatusMovedPermanently, "/api/tenders")
			return
		}

		tenders, err := s.tenderService.GetEmployeeTendersByUsername(ctx, username)
		if err != nil {
			if errors.Is(err, outerror.ErrEmployeeNotFound) {
				logger.Warn(fmt.Sprintf("employee with username=\"%s\" not found", username))
				c.JSON(
					http.StatusBadRequest,
					api.GetEmployeeTendersResponse{
						Message: fmt.Sprintf("employee with username \"%s\" not found", username),
					},
				)
				return
			} else if errors.Is(err, outerror.ErrEmployeeTendersNotFound) {
				logger.Warn(fmt.Sprintf("not found tenders for employee with username \"%s\"", username))
				c.JSON(
					http.StatusBadRequest,
					api.GetEmployeeTendersResponse{
						Message: fmt.Sprintf(
							"not found tenders for employee with username \"%s\"",
							username,
						),
					},
				)
				return
			} else {
				logger.Error("unexpected error", slog.String("err", err.Error()))
				c.JSON(http.StatusInternalServerError, api.GetEmployeeTendersResponse{Message: "internal error"})
				return
			}
		}
		logger.Info("success get employee tenders")
		c.JSON(http.StatusOK, api.GetEmployeeTendersResponse{Tenders: tenders, Message: "ok"})
	}
}
