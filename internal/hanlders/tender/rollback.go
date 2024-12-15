package tenderapi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sariya23/tender/internal/domain/models"
	schema "github.com/sariya23/tender/internal/hanlders"
	"github.com/sariya23/tender/internal/lib/unmarshal"
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
				http.StatusNotFound,
				schema.RollbackTenderResponse{
					Message: "tenderId must be positive integer number",
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
				http.StatusNotFound,
				schema.RollbackTenderResponse{
					Message: "version must be positive integer number",
				},
			)
			return
		}

		body := ginContext.Request.Body
		defer func() {
			err := body.Close()
			if err != nil {
				logger.Error("cannot close body", slog.String("err", err.Error()))
			}
		}()

		bodyData, err := io.ReadAll(body)
		if err != nil {
			logger.Error("cannot read body", slog.String("err", err.Error()))
			ginContext.JSON(http.StatusInternalServerError, schema.RollbackTenderResponse{Message: "internal error", RollbackTender: models.Tender{}})
			return
		}
		logger.Info("success read body")
		rollbackReq, err := unmarshal.RollbackRequest([]byte(bodyData))
		if err != nil {
			if errors.Is(err, unmarshal.ErrSyntax) {
				logger.Warn("req syntax error", slog.String("err", err.Error()))
				ginContext.JSON(
					http.StatusBadRequest,
					schema.RollbackTenderResponse{
						Message:        fmt.Sprintf("json syntax err: %s", err.Error()),
						RollbackTender: models.Tender{},
					},
				)
				return
			} else if errors.Is(err, unmarshal.ErrType) {
				logger.Warn("req type error", slog.String("err", err.Error()))
				ginContext.JSON(
					http.StatusBadRequest,
					schema.RollbackTenderResponse{
						Message:        fmt.Sprintf("json type err: %s", err.Error()),
						RollbackTender: models.Tender{},
					},
				)
				return
			} else {
				logger.Error("unexpected error", slog.String("err", err.Error()))
				ginContext.JSON(http.StatusInternalServerError, schema.RollbackTenderResponse{Message: "internal error", RollbackTender: models.Tender{}})
				return
			}
		}
		logger.Info("success unmarshal request")

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(&rollbackReq)
		if err != nil {
			logger.Error("validation error", slog.String("err", err.Error()))
			ginContext.JSON(
				http.StatusBadRequest,
				schema.RollbackTenderResponse{
					Message:        fmt.Sprintf("validation failed: %s", err.Error()),
					RollbackTender: models.Tender{},
				},
			)
			return
		}
		logger.Info("validate success")

		tender, err := tenderSrv.tenderService.RollbackTender(ctx, convertedTenderId, convertedVersion, rollbackReq.Username)
		if err != nil {
			if errors.Is(err, outerror.ErrTenderNotFound) {
				logger.Warn(fmt.Sprintf("tender with id=<%d> not found", convertedTenderId))
				ginContext.JSON(
					http.StatusNotFound,
					schema.RollbackTenderResponse{
						Message: fmt.Sprintf("tender with id=<%d> not found", convertedTenderId),
					},
				)
				return
			} else if errors.Is(err, outerror.ErrTenderVersionNotFound) {
				logger.Warn(fmt.Sprintf("tender with id=<%d> doesnt have version=<%d>", convertedTenderId, convertedVersion))
				ginContext.JSON(
					http.StatusNotFound,
					schema.RollbackTenderResponse{
						Message: fmt.Sprintf("tender with id=<%d> doesnt have version=<%d>", convertedTenderId, convertedVersion),
					},
				)
				return
			} else if errors.Is(err, outerror.ErrEmployeeNotResponsibleForTender) {
				logger.Warn(fmt.Sprintf("employee with username=<%s> not creator of tender with id=<%d>", rollbackReq.Username, convertedTenderId))
				ginContext.JSON(
					http.StatusForbidden,
					schema.RollbackTenderResponse{
						Message: fmt.Sprintf("employee with username=<%s> not creator of tender with id=<%d>", rollbackReq.Username, convertedTenderId),
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
