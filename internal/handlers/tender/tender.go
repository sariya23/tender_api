package tender

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"tender/internal/domain/models"
	"tender/internal/storage"

	"github.com/gin-gonic/gin"
)

type GetTendersResponse struct {
	Message     string          `json:"message"`
	Tenders     []models.Tender `json:"tenders"`
	ServiceType string          `json:"service_type"`
}

type CreateTenderRequest struct {
	Tender models.Tender `json:"tender"`
}

type CreateTenderResponse struct {
	Tender  models.Tender `json:"tender"`
	Message string        `json:"message"`
}

type TenderGetter interface {
	Tenders(ctx context.Context) ([]models.Tender, error)
	TendersByServiceType(ctx context.Context, serviceType string) ([]models.Tender, error)
}

type TenderCreater interface {
	Create(ctx context.Context, tender models.Tender) error
}

type UserProvider interface {
	GetUserByUsername(ctx context.Context, username string) error
}

type OrganizationProvider interface {
	GetOrganizationById(ctx context.Context, ogranizationId int) error
}

type UserResponsibler interface {
	CheckUserResponsible(ctx context.Context, username string, organixationId int) error
}

func GetTenders(ctx context.Context, logger *slog.Logger, tenderGetter TenderGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.tender.GetAllTenders"
		logger := logger.With("op", op)
		logger.Info("request to /api/tender/")

		serviceType := c.DefaultQuery("srv_type", "all")
		var err error
		var tenders []models.Tender

		if serviceType == "all" {
			tenders, err = tenderGetter.Tenders(ctx)
		} else {
			tenders, err = tenderGetter.TendersByServiceType(ctx, serviceType)
		}

		if err != nil {
			if errors.Is(err, storage.ErrNoTenderPresence) {
				logger.Info("no tenders found")
				c.JSON(http.StatusOK, GetTendersResponse{Message: "no tenders found", Tenders: []models.Tender{}, ServiceType: serviceType})
				return
			}
			logger.Error("cannot get tenders", slog.String("err", err.Error()))
			c.JSON(http.StatusBadRequest, GetTendersResponse{Message: "internal error", Tenders: []models.Tender{}})
			return
		}
		logger.Info("success get tenders")
		c.JSON(http.StatusOK, GetTendersResponse{Message: "ok", Tenders: tenders, ServiceType: serviceType})
	}
}

func CreateTender(
	ctx context.Context,
	logger *slog.Logger,
	tenderCreater TenderCreater,
	userProvider UserProvider,
	organizationProvider OrganizationProvider,
	userREponsilber UserResponsibler,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.tender.CreateTender"
		logger := logger.With("op", op)
		logger.Info("request to /api/tender/new")

		logger.Info("unmarshall request body")
		b := c.Request.Body
		defer func() {
			if err := b.Close(); err != nil {
				logger.Error("cannot close body", slog.String("err", err.Error()))
				return
			}
		}()

		var req CreateTenderRequest
		err := c.ShouldBindBodyWithJSON(&req)
		if err != nil {
			logger.Error("cannot unmarshall body", slog.String("err", err.Error()))
			c.JSON(
				http.StatusInternalServerError,
				CreateTenderResponse{Message: "internal error"})
			return
		}
		logger.Info("success unmarhsall", slog.String("tender", fmt.Sprintf("%+v", req.Tender)))

		logger.Info("checking existens of user")
		username := req.Tender.CreatorUsername
		err = userProvider.GetUserByUsername(ctx, username)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				logger.Warn("user not found", slog.String("username", username))
				c.JSON(http.StatusBadRequest, CreateTenderResponse{Message: "user not found"})
				return
			}
			logger.Error("cannot find user", slog.String("err", err.Error()))
			c.JSON(http.StatusInternalServerError, CreateTenderResponse{Message: "internal error"})
		}
		logger.Info("user exist")

		logger.Info("checking existens of existens")
		orgId := req.Tender.OrganizationId
		err = organizationProvider.GetOrganizationById(ctx, orgId)
		if err != nil {
			if errors.Is(err, storage.ErrOrganizationNotFound) {
				logger.Warn("organization not found", slog.Int("org_id", orgId))
				c.JSON(http.StatusBadRequest, CreateTenderResponse{Message: "org not found"})
				return
			}
			logger.Error("cannot find org", slog.String("err", err.Error()))
			c.JSON(http.StatusInternalServerError, CreateTenderResponse{Message: "internal error"})
		}
		logger.Info("org exist")

		logger.Info("checking responsible of user in organization")
		err = userREponsilber.CheckUserResponsible(ctx, username, orgId)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotReponsibleForOrg) {
				logger.Warn("user not responsible for organization", slog.String("username", username), slog.Int("org_id", orgId))
				c.JSON(http.StatusBadRequest, CreateTenderResponse{Message: "user not responsible for organization"})
				return
			}
			logger.Error("cannot check user responsibility", slog.String("err", err.Error()))
			c.JSON(http.StatusInternalServerError, CreateTenderResponse{Message: "internal error"})
		}

		logger.Info("user responsible", slog.String("username", username), slog.Int("org_id", orgId))

		logger.Info("processing create tender")
		err = tenderCreater.Create(ctx, req.Tender)
		if err != nil {
			logger.Error("cannot create tender", slog.String("err", err.Error()))
			c.JSON(
				http.StatusInternalServerError,
				GetTendersResponse{Message: "internal error"},
			)
		}
		logger.Info("tender created success", slog.String("tender name", req.Tender.TenderName))
		c.JSON(http.StatusOK, CreateTenderResponse{Tender: req.Tender, Message: "ok"})
	}
}
