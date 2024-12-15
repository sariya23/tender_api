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

func (tenderSrv *TenderService) EditTender(ctx context.Context) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		const opeartionPlace = "internal.api.tenderapi.EditTedner"
		logger := tenderSrv.logger.With("op", opeartionPlace)
		logger.Info(fmt.Sprintf("request to %v", ginContext.Request.URL.Path))

		tenderId := ginContext.Param("tenderId")
		convertedTenderId, err := strconv.Atoi(tenderId)
		if err != nil {
			logger.Error(
				"cannot convert tender id to int",
				slog.String("tender id", tenderId),
				slog.String("err", err.Error()),
			)
			ginContext.JSON(http.StatusBadRequest, schema.EditTenderResponse{Message: "wrong path", UpdatedTender: models.Tender{}})
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
			ginContext.JSON(http.StatusInternalServerError, schema.EditTenderResponse{Message: "internal error", UpdatedTender: models.Tender{}})
			return
		}
		logger.Info("success read body")

		updatedReq, err := unmarshal.EditRequest(bodyData)
		if err != nil {
			if errors.Is(err, unmarshal.ErrSyntax) {
				logger.Warn("req syntax error", slog.String("err", err.Error()))
				ginContext.JSON(
					http.StatusBadRequest,
					schema.EditTenderResponse{
						Message:       fmt.Sprintf("json syntax err: %s", err.Error()),
						UpdatedTender: models.Tender{},
					},
				)
				return
			} else if errors.Is(err, unmarshal.ErrType) {
				logger.Warn("req type error", slog.String("err", err.Error()))
				ginContext.JSON(
					http.StatusBadRequest,
					schema.EditTenderResponse{
						Message:       fmt.Sprintf("json type err: %s", err.Error()),
						UpdatedTender: models.Tender{},
					},
				)
				return
			} else {
				logger.Error("unexpected error", slog.String("err", err.Error()))
				ginContext.JSON(
					http.StatusInternalServerError,
					schema.EditTenderResponse{
						Message:       "internal error",
						UpdatedTender: models.Tender{},
					},
				)
				return
			}
		}
		logger.Info("success unmarshal request")

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(&updatedReq)
		if err != nil {
			logger.Error("validation error", slog.String("err", err.Error()))
			ginContext.JSON(
				http.StatusBadRequest,
				schema.EditTenderResponse{
					Message:       fmt.Sprintf("validation failed: %s", err.Error()),
					UpdatedTender: models.Tender{},
				},
			)
			return
		}
		logger.Info("validate success")

		tender, err := tenderSrv.tenderService.EditTender(ctx, convertedTenderId, updatedReq.UpdateTenderData, updatedReq.Username)

		if err != nil {
			if errors.Is(err, outerror.ErrUnknownTenderStatus) {
				logger.Warn(fmt.Sprintf("tender status \"%s\" is unknown", *updatedReq.UpdateTenderData.Status))
				ginContext.JSON(
					http.StatusBadRequest,
					schema.EditTenderResponse{
						Message:       fmt.Sprintf("tender status \"%s\" is unknown", *updatedReq.UpdateTenderData.Status),
						UpdatedTender: models.Tender{},
					},
				)
				return
			} else if errors.Is(err, outerror.ErrCannotSetThisTenderStatus) {
				logger.Warn("cannot set this tender status")
				ginContext.JSON(
					http.StatusBadRequest,
					schema.EditTenderResponse{
						Message:       "cannot set this tender status. Cannot set tender status from PUBLISHED to CREATED and from CLOSED to CREATED",
						UpdatedTender: models.Tender{},
					},
				)
				return
			} else if errors.Is(err, outerror.ErrCannotSetThisTenderStatus) {
				logger.Warn("cannot update tender status")
				ginContext.JSON(
					http.StatusBadRequest,
					schema.EditTenderResponse{
						Message:       fmt.Sprintf("cannot update tender status: %v", err.Error()),
						UpdatedTender: models.Tender{},
					},
				)
				return
			} else if errors.Is(err, outerror.ErrEmployeeNotResponsibleForTender) {
				logger.Warn("employee not respobsible for tender")
				ginContext.JSON(
					http.StatusBadRequest,
					schema.EditTenderResponse{
						Message: fmt.Sprintf(
							"employee with username \"%s\" not respobsible for tender with id \"%d\"",
							updatedReq.Username,
							convertedTenderId,
						),
						UpdatedTender: models.Tender{},
					},
				)
			} else if errors.Is(err, outerror.ErrTenderNotFound) {
				logger.Warn(fmt.Sprintf("tender with id=\"%d\" not found", convertedTenderId))
				ginContext.JSON(
					http.StatusBadRequest,
					schema.EditTenderResponse{
						Message:       fmt.Sprintf("tender with id=\"%d\" not found", convertedTenderId),
						UpdatedTender: models.Tender{},
					},
				)
				return
			} else if errors.Is(err, outerror.ErrEmployeeNotFound) {
				logger.Warn(
					fmt.Sprintf(
						"updated employee with username=\"%s\" not found",
						*updatedReq.UpdateTenderData.CreatorUsername,
					),
				)
				ginContext.JSON(
					http.StatusBadRequest,
					schema.EditTenderResponse{
						Message: fmt.Sprintf(
							"updated employee with username=\"%s\" not found",
							*updatedReq.UpdateTenderData.CreatorUsername,
						),
						UpdatedTender: models.Tender{},
					},
				)
				return
			} else if errors.Is(err, outerror.ErrOrganizationNotFound) {
				logger.Warn(
					fmt.Sprintf(
						"updated organization with id=\"%d\" not found",
						*updatedReq.UpdateTenderData.OrganizationId,
					),
				)
				ginContext.JSON(
					http.StatusBadRequest,
					schema.EditTenderResponse{
						Message: fmt.Sprintf(
							"updated organization with id=\"%d\" not found",
							*updatedReq.UpdateTenderData.OrganizationId,
						),
						UpdatedTender: models.Tender{},
					},
				)
				return
			} else if errors.Is(err, outerror.ErrEmployeeNotResponsibleForOrganization) {
				logger.Warn(
					fmt.Sprintf(
						"new employee with username=\"%s\" not responsible for new organization with id=\"%d\"",
						*updatedReq.UpdateTenderData.CreatorUsername,
						*updatedReq.UpdateTenderData.OrganizationId,
					),
				)
				ginContext.JSON(
					http.StatusBadRequest,
					schema.EditTenderResponse{
						Message: fmt.Sprintf(
							"employee with username=\"%s\" not responsible for organization with id=\"%d\"",
							*updatedReq.UpdateTenderData.CreatorUsername,
							*updatedReq.UpdateTenderData.OrganizationId,
						),
						UpdatedTender: models.Tender{},
					},
				)
				return
			} else {
				logger.Error("unexpected error", slog.String("err", err.Error()))
				ginContext.JSON(http.StatusInternalServerError, schema.EditTenderResponse{Message: "internal error", UpdatedTender: models.Tender{}})
				return
			}
		}
		logger.Info("tender updated success")
		ginContext.JSON(http.StatusOK, schema.EditTenderResponse{Message: "ok", UpdatedTender: tender})
	}
}
