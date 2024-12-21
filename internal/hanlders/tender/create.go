package tenderapi

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sariya23/tender/internal/domain/models"
	schema "github.com/sariya23/tender/internal/hanlders"
	"github.com/sariya23/tender/internal/lib/unmarshal"
	outerror "github.com/sariya23/tender/internal/out_error"
)

func (tenderSrv *TenderService) CreateTender(ctx context.Context) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		const operationPlace = "internal.api.tenderapi.CreateTender"
		logger := tenderSrv.logger.With("op", operationPlace)
		logger.Info(fmt.Sprintf("request to %v", ginContext.Request.URL))

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
			ginContext.JSON(http.StatusInternalServerError, schema.CreateTenderResponse{Message: "internal error", Tender: models.Tender{}})
			return
		}
		logger.Info("success read body")
		createReq, err := unmarshal.CreateRequest([]byte(bodyData))
		if err != nil {
			if errors.Is(err, unmarshal.ErrSyntax) {
				logger.Warn("req syntax error", slog.String("err", err.Error()))
				ginContext.JSON(
					http.StatusBadRequest,
					schema.CreateTenderResponse{
						Message: fmt.Sprintf("json syntax err: %s", err.Error()),
						Tender:  models.Tender{},
					},
				)
				return
			} else if errors.Is(err, unmarshal.ErrType) {
				logger.Warn("req type error", slog.String("err", err.Error()))
				ginContext.JSON(
					http.StatusBadRequest,
					schema.CreateTenderResponse{
						Message: fmt.Sprintf("json type err: %s", err.Error()),
						Tender:  models.Tender{},
					},
				)
				return
			} else {
				logger.Error("unexpected error", slog.String("err", err.Error()))
				ginContext.JSON(http.StatusInternalServerError, schema.CreateTenderResponse{Message: "internal error", Tender: models.Tender{}})
				return
			}
		}
		logger.Info("success unmarshal request")

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(&createReq)
		if err != nil {
			logger.Error("validation error", slog.String("err", err.Error()))
			ginContext.JSON(
				http.StatusBadRequest,
				schema.CreateTenderResponse{
					Message: fmt.Sprintf("validation failed: %s", err.Error()),
					Tender:  models.Tender{},
				},
			)
			return
		}
		logger.Info("validate success")
		tender, err := tenderSrv.tenderService.CreateTender(ctx, createReq.Tender)
		if err != nil {
			if errors.Is(err, outerror.ErrEmployeeNotFound) {
				logger.Warn("employee not found", slog.String("err", err.Error()))
				ginContext.JSON(
					http.StatusUnprocessableEntity,
					schema.CreateTenderResponse{
						Message: fmt.Sprintf(
							"employee with username=<%s> not found",
							createReq.Tender.CreatorUsername,
						),
						Tender: models.Tender{},
					},
				)
				return
			} else if errors.Is(err, outerror.ErrOrganizationNotFound) {
				logger.Warn("organization not found", slog.String("err", err.Error()))
				ginContext.JSON(
					http.StatusBadRequest,
					schema.CreateTenderResponse{
						Message: fmt.Sprintf(
							"organization with id=<%d> not found",
							createReq.Tender.OrganizationId,
						),
						Tender: models.Tender{},
					},
				)
				return
			} else if errors.Is(err, outerror.ErrEmployeeNotResponsibleForOrganization) {
				logger.Warn("employee not responsible for organization", slog.String("err", err.Error()))
				ginContext.JSON(
					http.StatusBadRequest,
					schema.CreateTenderResponse{
						Message: fmt.Sprintf(
							"employee <%s> not responsible for organization with id=<%d>",
							createReq.Tender.CreatorUsername,
							createReq.Tender.OrganizationId,
						),
						Tender: models.Tender{},
					},
				)
				return
			} else if errors.Is(err, outerror.ErrNewTenderCannotCreatedWithStatusNotCreated) {
				logger.Warn("cannot create tender with status", slog.String("status", createReq.Tender.Status))
				ginContext.JSON(
					http.StatusBadRequest,
					schema.CreateTenderResponse{
						Message: fmt.Sprintf("cannot create tender with status <%s>", createReq.Tender.Status),
						Tender:  models.Tender{},
					},
				)
				return
			} else {
				logger.Error("unexpected error", slog.String("err", err.Error()))
				ginContext.JSON(
					http.StatusInternalServerError,
					schema.CreateTenderResponse{
						Message: "internal error",
						Tender:  models.Tender{},
					},
				)
				return
			}
		}
		logger.Info("tender created success")
		ginContext.JSON(http.StatusOK, schema.CreateTenderResponse{Message: "ok", Tender: tender})
	}
}
