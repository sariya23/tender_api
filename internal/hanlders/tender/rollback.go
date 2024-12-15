package tenderapi

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	schema "github.com/sariya23/tender/internal/hanlders"
	outerror "github.com/sariya23/tender/internal/out_error"
)

func (tenderSrv *TenderService) RollbackTender(ctx context.Context) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		const operationPlace = "internal.api.tenderapi.RollbackTender"
		logger := tenderSrv.logger.With("op", operationPlace)
		logger.Info(fmt.Sprintf("request to %v", ginContext.Request.URL.Path))

		tenderId := ginContext.Param("tenderId")
		convertedTenderId, err := strconv.Atoi(tenderId)
		if err != nil {
			logger.Error(
				"cannot convert tender id to int",
				slog.String("tender id", tenderId),
				slog.String("err", err.Error()),
			)
			ginContext.JSON(
				http.StatusBadRequest,
				schema.RollbackTenderResponse{
					Message: fmt.Sprintf("wrong path: %s", ginContext.Request.URL.Path),
				},
			)
			return
		}

		version := ginContext.Param("version")
		convertedVersion, err := strconv.Atoi(version)
		if err != nil {
			logger.Error(
				"cannot convert version to int",
				slog.String("version", version),
				slog.String("err", err.Error()),
			)
			ginContext.JSON(
				http.StatusBadRequest,
				schema.RollbackTenderResponse{
					Message: fmt.Sprintf("wrong path: %s", ginContext.Request.URL.Path),
				},
			)
			return
		}

		tender, err := tenderSrv.tenderService.RollbackTender(ctx, convertedTenderId, convertedVersion)
		if err != nil {
			if errors.Is(err, outerror.ErrTenderNotFound) {
				logger.Warn(fmt.Sprintf("tender with id=\"%d\" not found", convertedTenderId))
				ginContext.JSON(
					http.StatusBadRequest,
					schema.RollbackTenderResponse{
						Message: fmt.Sprintf("tender with id=\"%d\" not found", convertedTenderId),
					},
				)
				return
			} else if errors.Is(err, outerror.ErrTenderVersionNotFound) {
				logger.Warn(fmt.Sprintf("tender with id=\"%d\" doesnt have version \"%d\"", convertedTenderId, convertedVersion))
				ginContext.JSON(
					http.StatusBadRequest,
					schema.RollbackTenderResponse{
						Message: fmt.Sprintf("tender with id=\"%d\" doesnt have version \"%d\"", convertedTenderId, convertedVersion),
					},
				)
				return
			} else {
				logger.Error("unexpected error", slog.String("err", err.Error()))
				ginContext.JSON(http.StatusInternalServerError, schema.RollbackTenderResponse{Message: "internal error"})
				return
			}
		}

		logger.Info("rollback success")
		ginContext.JSON(http.StatusOK, schema.RollbackTenderResponse{Message: "ok", RollbackTender: tender})
	}
}
