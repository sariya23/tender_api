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
	return func(ginContext *gin.Context) {
		const operationPlace = "internal.api.tenderapi.GetTenders"
		logger := s.logger.With("op", operationPlace)
		logger.Info(fmt.Sprintf("request to %v", ginContext.Request.URL))

		serviceType := ginContext.DefaultQuery("srv_type", "all")
		tenders, err := s.tenderService.GetTenders(ctx, serviceType)
		if err != nil {
			if errors.Is(err, outerror.ErrTendersWithThisServiceTypeNotFound) {
				ginContext.JSON(
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
				ginContext.JSON(http.StatusInternalServerError, api.GetTendersResponse{Message: "internal error"})
				return
			}
		}

		logger.Info("send success response")
		ginContext.JSON(http.StatusOK, api.GetTendersResponse{Message: "ok", Tenders: tenders})
	}
}

func (s *TenderService) GetEmployeeTendersByUsername(ctx context.Context) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		const op = "internal.api.tenderapi.GetEmployeeTendersByUsername"
		logger := s.logger.With("op", op)
		logger.Info(fmt.Sprintf("request to %s", ginContext.Request.URL))

		username := ginContext.Query("username")
		if username == "" {
			logger.Info("username not specified. Rediredr to /api/tenders")
			ginContext.Redirect(http.StatusMovedPermanently, "/api/tenders")
			return
		}
		logger.Info("try get employee tenders", slog.String("username", username))
		tenders, err := s.tenderService.GetEmployeeTendersByUsername(ctx, username)
		if err != nil {
			if errors.Is(err, outerror.ErrEmployeeNotFound) {
				logger.Warn(fmt.Sprintf("employee with username=\"%s\" not found", username))
				ginContext.JSON(
					http.StatusBadRequest,
					api.GetEmployeeTendersResponse{
						Message: fmt.Sprintf("employee with username \"%s\" not found", username),
					},
				)
				return
			} else if errors.Is(err, outerror.ErrEmployeeTendersNotFound) {
				logger.Warn(fmt.Sprintf("not found tenders for employee with username \"%s\"", username))
				ginContext.JSON(
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
				ginContext.JSON(http.StatusInternalServerError, api.GetEmployeeTendersResponse{Message: "internal error"})
				return
			}
		}
		logger.Info("success get employee tenders")
		ginContext.JSON(http.StatusOK, api.GetEmployeeTendersResponse{Tenders: tenders, Message: "ok"})
	}
}
