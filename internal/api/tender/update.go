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

		_, err = s.tenderService.Edit(ctx, convertedTenderId, updatedReq.UpdateTenderData)
	}
}
