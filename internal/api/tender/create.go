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
	"github.com/sariya23/tender/internal/api"
	"github.com/sariya23/tender/internal/lib/unmarshal"
	outerror "github.com/sariya23/tender/internal/out_error"
)

func (s *TenderService) CreateTender(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "internal.api.tenderapi.CreateTender"
		logger := s.logger.With("op", op)
		logger.Info(fmt.Sprintf("request to %v", c.Request.URL))

		b := c.Request.Body
		defer func() {
			err := b.Close()
			if err != nil {
				logger.Error("cannot close body", slog.String("err", err.Error()))
			}
		}()

		body, err := io.ReadAll(b)
		if err != nil {
			logger.Error("cannot read body", slog.String("err", err.Error()))
			c.JSON(http.StatusInternalServerError, api.CreateTenderResponse{Message: "internal error"})
			return
		}
		logger.Info("success read body")
		createReq, err := unmarshal.CreateRequest([]byte(body))
		if err != nil {
			if errors.Is(err, unmarshal.ErrSyntax) {
				logger.Warn("req syntax error", slog.String("err", err.Error()))
				c.JSON(http.StatusBadRequest, api.CreateTenderResponse{Message: fmt.Sprintf("json syntax err: %s", err.Error())})
				return
			} else if errors.Is(err, unmarshal.ErrType) {
				logger.Warn("req type error", slog.String("err", err.Error()))
				c.JSON(http.StatusBadRequest, api.CreateTenderResponse{Message: fmt.Sprintf("json type err: %s", err.Error())})
				return
			} else {
				logger.Error("unexpected error", slog.String("err", err.Error()))
				c.JSON(http.StatusInternalServerError, api.CreateTenderResponse{Message: "internal error"})
				return
			}
		}
		logger.Info("success unmarshal request")

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(&createReq)
		if err != nil {
			logger.Error("validation error", slog.String("err", err.Error()))
			c.JSON(http.StatusBadRequest, api.CreateTenderResponse{Message: fmt.Sprintf("validation failed: %s", err.Error())})
			return
		}
		logger.Info("validate success")
		tender, err := s.tenderService.CreateTender(ctx, createReq.Tender)
		if err != nil {
			if errors.Is(err, outerror.ErrEmployeeNotFound) {
				logger.Warn("employee not found", slog.String("err", err.Error()))
				c.JSON(
					http.StatusBadRequest,
					api.CreateTenderResponse{
						Message: fmt.Sprintf(
							"employee with username=\"%s\" not found",
							createReq.Tender.CreatorUsername,
						),
					},
				)
				return
			} else if errors.Is(err, outerror.ErrOrganizationNotFound) {
				logger.Warn("organization not found", slog.String("err", err.Error()))
				c.JSON(
					http.StatusBadRequest,
					api.CreateTenderResponse{
						Message: fmt.Sprintf(
							"organization with id=%d not found",
							createReq.Tender.OrganizationId,
						),
					},
				)
				return
			} else if errors.Is(err, outerror.ErrEmployeeNotResponsibleForOrganization) {
				logger.Warn("employee not responsible for organization", slog.String("err", err.Error()))
				c.JSON(
					http.StatusBadRequest,
					api.CreateTenderResponse{
						Message: fmt.Sprintf(
							"employee \"%s\" not responsible for organization with id=%d",
							createReq.Tender.CreatorUsername,
							createReq.Tender.OrganizationId,
						),
					},
				)
				return
			} else if errors.Is(err, outerror.ErrUnknownTenderStatus) {
				logger.Warn("unknown tender status", slog.String("status", createReq.Tender.Status))
				c.JSON(http.StatusBadRequest, api.CreateTenderResponse{Message: fmt.Sprintf("unknown tender status \"%s\"", createReq.Tender.Status)})
				return
			} else {
				logger.Error("unexpected error", slog.String("err", err.Error()))
				c.JSON(http.StatusInternalServerError, api.CreateTenderResponse{Message: "internal error"})
				return
			}
		}
		logger.Info("tender created success")
		c.JSON(http.StatusOK, api.CreateTenderResponse{Message: "ok", Tender: tender})
	}
}
