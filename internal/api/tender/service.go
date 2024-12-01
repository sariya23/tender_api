package tender

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sariya23/tender/internal/api"
	outerror "github.com/sariya23/tender/internal/out_error"
	"github.com/sariya23/tender/internal/service"
)

type TenderService struct {
	logger        *slog.Logger
	tenderService service.TenderServiceProvider
}

func New(logger *slog.Logger, tenderService service.TenderServiceProvider) *TenderService {
	return &TenderService{
		logger:        logger,
		tenderService: tenderService,
	}
}

func (s *TenderService) GetTenders(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "internal.api.tender.service.GetTenders"
		logger := s.logger.With("op", op)
		logger.Info(fmt.Sprintf("request to %v", c.Request.URL))

		serviceType := c.DefaultQuery("srv_type", "all")
		tenders, err := s.tenderService.GetTenders(ctx, serviceType)
		if err != nil {
			if errors.Is(err, outerror.ErrTendersWithThisServiceTypeNotFound) {
				c.JSON(http.StatusBadRequest, api.GetTendersResponse{Message: fmt.Sprintf("no tenders found with service type: %s", serviceType)})
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

func (s *TenderService) CreateTender(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "internal.api.tender.service.CreateTender"
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

		var createRequest api.CreateTenderRequest
		err = json.Unmarshal(body, &createRequest)
		if err != nil {
			var syntaxErr *json.SyntaxError
			var typeErr *json.UnmarshalTypeError
			if errors.As(err, &syntaxErr) {
				logger.Warn("cannot unmarhal types in go type", slog.String("err", typeErr.Error()))
				c.JSON(http.StatusBadRequest, api.CreateTenderResponse{Message: "wrong types"})
				return
			}
			logger.Error("cannot unmarshal body", slog.String("err", err.Error()))
			c.JSON(http.StatusInternalServerError, api.CreateTenderResponse{Message: "internal error"})
			return
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(&createRequest)
		if err != nil {
			logger.Error("cannot validate fields", slog.String("err", err.Error()))
			c.JSON(http.StatusInternalServerError, api.CreateTenderResponse{Message: "internal error"})
		}
	}
}
