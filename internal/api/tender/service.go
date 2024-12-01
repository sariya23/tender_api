package tender

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
				c.JSON(
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
		logger.Info("success read body")
		createReq, err := unmarshal.CreateRequest([]byte(body))
		if err != nil {
			if errors.Is(err, unmarshal.ErrSyntax) {
				logger.Warn("req syntax error", slog.String("err", err.Error()))
				c.JSON(http.StatusBadRequest, api.CreateTenderResponse{Message: fmt.Sprintf("json err syntax: %s", err.Error())})
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
			c.JSON(http.StatusBadRequest, api.CreateTenderResponse{Message: fmt.Sprintf("validation faild: %s", err.Error())})
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
				logger.Warn("user not responsible for irganization", slog.String("err", err.Error()))
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
			} else {
				logger.Error("unexpected error", slog.String("err", err.Error()))
				c.JSON(http.StatusInternalServerError, api.CreateTenderResponse{Message: "internal error"})
			}
		}
		logger.Info("tender created success")
		c.JSON(http.StatusOK, api.CreateTenderResponse{Message: "ok", Tender: tender})
	}
}
