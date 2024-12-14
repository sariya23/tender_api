package tenderapi

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sariya23/tender/internal/api"
	outerror "github.com/sariya23/tender/internal/out_error"
)

func (s *TenderService) RollbackTender(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		const operationPlace = "internal.api.tenderapi.RollbackTender"
		logger := s.logger.With("op", operationPlace)
		logger.Info(fmt.Sprintf("request to %v", c.Request.URL.Path))

		tenderId := c.Param("tenderId")
		convertedTenderId, err := strconv.Atoi(tenderId)
		if err != nil {
			logger.Error(
				"cannot convert tender id to int",
				slog.String("tender id", tenderId),
				slog.String("err", err.Error()),
			)
			c.JSON(http.StatusBadRequest, api.RollbackTenderResponse{Message: fmt.Sprintf("wrong path: %s", c.Request.URL.Path)})
			return
		}

		version := c.Param("version")
		convertedVersion, err := strconv.Atoi(version)
		if err != nil {
			logger.Error(
				"cannot convert version to int",
				slog.String("version", version),
				slog.String("err", err.Error()),
			)
			c.JSON(http.StatusBadRequest, api.RollbackTenderResponse{Message: fmt.Sprintf("wrong path: %s", c.Request.URL.Path)})
			return
		}

		tender, err := s.tenderService.RollbackTender(ctx, convertedTenderId, convertedVersion)
		if err != nil {
			if errors.Is(err, outerror.ErrTenderNotFound) {
				logger.Warn(fmt.Sprintf("tender with id=\"%d\" not found", convertedTenderId))
				c.JSON(http.StatusBadRequest, api.RollbackTenderResponse{Message: fmt.Sprintf("tender with id=\"%d\" not found", convertedTenderId)})
				return
			} else if errors.Is(err, outerror.ErrTenderVersionNotFound) {
				logger.Warn(fmt.Sprintf("tender with id=\"%d\" doesnt have version \"%d\"", convertedTenderId, convertedVersion))
				c.JSON(
					http.StatusBadRequest,
					api.RollbackTenderResponse{
						Message: fmt.Sprintf("tender with id=\"%d\" doesnt have version \"%d\"", convertedTenderId, convertedVersion),
					},
				)
				return
			} else {
				logger.Error("unexpected error", slog.String("err", err.Error()))
				c.JSON(http.StatusInternalServerError, api.RollbackTenderResponse{Message: "internal error"})
				return
			}
		}

		logger.Info("rollback success")
		c.JSON(http.StatusOK, api.RollbackTenderResponse{Message: "ok", RollbackTender: tender})
	}
}
