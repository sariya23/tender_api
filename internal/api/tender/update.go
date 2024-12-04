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
	"github.com/sariya23/tender/internal/api"
	"github.com/sariya23/tender/internal/lib/unmarshal"
	outerror "github.com/sariya23/tender/internal/out_error"
)

func (s *TenderService) EditTedner(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "internal.api.tenderapi.EditTedner"
		logger := s.logger.With("op", op)
		logger.Info(fmt.Sprintf("request to %v", c.Request.URL.Path))

		tenderId := c.Param("tenderId")
		convertedTenderId, err := strconv.Atoi(tenderId)
		if err != nil {
			logger.Error(
				"cannot convert tender id to int",
				slog.String("tender id", tenderId),
				slog.String("err", err.Error()),
			)
			c.JSON(http.StatusBadRequest, api.EditTenderResponse{Message: "wrong path"})
			return
		}

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
			c.JSON(http.StatusInternalServerError, api.EditTenderResponse{Message: "internal error"})
			return
		}
		logger.Info("success read body")

		updatedReq, err := unmarshal.EditRequest(body)
		if err != nil {
			if errors.Is(err, unmarshal.ErrSyntax) {
				logger.Warn("req syntax error", slog.String("err", err.Error()))
				c.JSON(http.StatusBadRequest, api.EditTenderResponse{Message: fmt.Sprintf("json syntax err: %s", err.Error())})
				return
			} else if errors.Is(err, unmarshal.ErrType) {
				logger.Warn("req type error", slog.String("err", err.Error()))
				c.JSON(http.StatusBadRequest, api.EditTenderResponse{Message: fmt.Sprintf("json type err: %s", err.Error())})
				return
			} else {
				logger.Error("unexpected error", slog.String("err", err.Error()))
				c.JSON(http.StatusInternalServerError, api.EditTenderResponse{Message: "internal error"})
				return
			}
		}
		logger.Info("success unmarshal request")

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(&updatedReq)
		if err != nil {
			logger.Error("validation error", slog.String("err", err.Error()))
			c.JSON(http.StatusBadRequest, api.CreateTenderResponse{Message: fmt.Sprintf("validation faild: %s", err.Error())})
			return
		}
		logger.Info("validate success")

		tender, err := s.tenderService.EditTender(ctx, convertedTenderId, updatedReq.UpdateTenderData)

		if err != nil {
			if errors.Is(err, outerror.ErrTenderNotFound) {
				logger.Warn(fmt.Sprintf("tender with id=\"%d\" not found", convertedTenderId))
				c.JSON(
					http.StatusBadRequest,
					api.EditTenderResponse{
						Message: fmt.Sprintf("tender with id=\"%d\" not found", convertedTenderId),
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
				c.JSON(
					http.StatusBadRequest,
					api.EditTenderResponse{
						Message: fmt.Sprintf(
							"updated employee with username=\"%s\" not found",
							*updatedReq.UpdateTenderData.CreatorUsername,
						),
					},
				)
				return
			} else if errors.Is(err, outerror.ErrOrganizationNotFound) {
				logger.Warn(
					fmt.Sprintf(
						"updated organization with id=\"%d\" not found",
						updatedReq.UpdateTenderData.OrganizationId,
					),
				)
				c.JSON(
					http.StatusBadRequest,
					api.EditTenderResponse{
						Message: fmt.Sprintf(
							"updated organization with id=\"%d\" not found",
							updatedReq.UpdateTenderData.OrganizationId,
						),
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
				c.JSON(
					http.StatusBadRequest,
					api.EditTenderResponse{
						Message: fmt.Sprintf(
							"employee with username=\"%s\" not responsible for organization with id=\"%d\"",
							*updatedReq.UpdateTenderData.CreatorUsername,
							*updatedReq.UpdateTenderData.OrganizationId,
						),
					},
				)
				return
			} else {
				logger.Error("unexpected error", slog.String("err", err.Error()))
				c.JSON(http.StatusInternalServerError, api.EditTenderResponse{Message: "internal error"})
				return
			}
		}
		logger.Info("tender updated success")
		c.JSON(http.StatusOK, api.EditTenderResponse{Message: "ok", UpdatedTender: tender})
	}
}
