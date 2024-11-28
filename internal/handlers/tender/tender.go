package tender

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"tender/internal/domain/models"
	"tender/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type GetTendersResponse struct {
	Message     string          `json:"message"`
	Tenders     []models.Tender `json:"tenders"`
	ServiceType string          `json:"service_type"`
}

type CreateTenderRequest struct {
	Tender models.Tender `json:"tender" validate:"required"`
}

type CreateTenderResponse struct {
	Tender  models.Tender `json:"tender"`
	Message string        `json:"message"`
}

type GetUserTendersResponse struct {
	Message string          `json:"message"`
	Tenders []models.Tender `json:"tenders"`
}

type EditTenderRequest struct {
	TenderName      string `json:"name"`
	Description     string `json:"description"`
	ServiceType     string `json:"serviceType"`
	Status          string `json:"status"`
	OrganizationId  int    `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}

type EditTenderResponse struct {
	UpdatedTender models.Tender `json:"updated_tender"`
	Message       string        `json:"message"`
}

type TenderGetter interface {
	Tenders(ctx context.Context) ([]models.Tender, error)
	TendersByServiceType(ctx context.Context, serviceType string) ([]models.Tender, error)
}

type TenderCreater interface {
	CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error)
}

type UserProvider interface {
	UserByUsername(ctx context.Context, username string) error
}

type OrganizationProvider interface {
	OrganizationById(ctx context.Context, ogranizationId int) error
}

type UserResponsibler interface {
	CheckUserResponsible(ctx context.Context, username string, organizationId int) error
}

type UserTenderGetter interface {
	UserTenders(ctx context.Context, username string) ([]models.Tender, error)
}

type TenderEditor interface {
	EditTender(ctx context.Context, tenderId int, newTender models.Tender) (models.Tender, error)
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
	userResponsibler UserResponsibler,
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
			logMsg, respMsg, code := handleErrorWhileUnmarshallCreateTenderRequest(err)
			logger.Error(logMsg, slog.String("err", err.Error()))
			c.JSON(code, CreateTenderResponse{Message: respMsg})
			return
		}
		logger.Info("success unmarhsall", slog.String("tender", fmt.Sprintf("%+v", req.Tender)))

		err = validateCreateTenderRequest(req)
		if err != nil {
			logger.Error("invalid request struct", slog.String("err", err.Error()))
			c.JSON(http.StatusBadRequest, CreateTenderResponse{Message: "invalid fields"})
			return
		}

		logger.Info("checking existens of user")
		username := req.Tender.CreatorUsername
		err = userProvider.UserByUsername(ctx, username)
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

		logger.Info("checking existens of organization")
		orgId := req.Tender.OrganizationId
		err = organizationProvider.OrganizationById(ctx, orgId)
		if err != nil {
			if errors.Is(err, storage.ErrOrganizationNotFound) {
				logger.Warn("organization not found", slog.Int("org_id", orgId))
				c.JSON(http.StatusBadRequest, CreateTenderResponse{Message: "organization not found"})
				return
			}
			logger.Error("cannot find org", slog.String("err", err.Error()))
			c.JSON(http.StatusInternalServerError, CreateTenderResponse{Message: "internal error"})
		}
		logger.Info("org exist")
		// TODO: посмотреть получится ли отправлять один запрос. Не проверять наличие юзера и организации, а получать
		// в случае чего ошибку из проверки ответсвенности.
		logger.Info("checking responsible of user in organization")
		err = userResponsibler.CheckUserResponsible(ctx, username, orgId)
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
		tender, err := tenderCreater.CreateTender(ctx, req.Tender)
		if err != nil {
			logger.Error("cannot create tender", slog.String("err", err.Error()))
			c.JSON(
				http.StatusInternalServerError,
				GetTendersResponse{Message: "internal error"},
			)
		}
		logger.Info("tender created success", slog.String("tender name", req.Tender.TenderName))
		c.JSON(http.StatusOK, CreateTenderResponse{Tender: tender, Message: "ok"})
	}
}

func GetUserTenders(ctx context.Context, logger *slog.Logger, userTenderGetter UserTenderGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.tender.GetUserTenders"
		logger := logger.With("op", op)
		username := c.Query("username")
		logger.Info(fmt.Sprintf("request to /api/tender/my?username=%v", username))
		if username == "" {
			logger.Info("username not specified - redirect to get all tenders")
			c.Redirect(http.StatusMovedPermanently, "/api/tender/")
			return
		}

		tenders, err := userTenderGetter.UserTenders(ctx, username)
		if err != nil {
			if errors.Is(err, storage.ErrNoTenderForThisUser) {
				logger.Warn("no tenders for this user", slog.String("username", username))
				c.JSON(http.StatusOK, GetUserTendersResponse{Message: "no tenders for this user"})
				return
			}
			logger.Error("unexpected error", slog.String("err", err.Error()))
			c.JSON(http.StatusInternalServerError, GetTendersResponse{Message: "internal error"})
			return
		}
		logger.Info("success get tenders", slog.String("username", username), "tenders", tenders)
		c.JSON(http.StatusOK, GetUserTendersResponse{Message: "ok", Tenders: tenders})
	}
}

func EditTender(ctx context.Context, logger *slog.Logger, tenderEditor TenderEditor) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.tender.EditTender"
		logger := logger.With("op", op)
		logger.Info(fmt.Sprintf("request to /api/tender/%v/edit/", c.Param("tender_id")))
		tenderId, err := strconv.Atoi(c.Param("tender_id"))
		if err != nil {
			logger.Error("cannot parse tender id to int", slog.String("err", err.Error()))
			c.JSON(http.StatusBadRequest, EditTenderResponse{Message: "cannot parse tender id"})
			return
		}

		b := c.Request.Body
		defer func() {
			if err := b.Close(); err != nil {
				logger.Error("cannot close body", slog.String("err", err.Error()))
				return
			}
		}()

		var req EditTenderRequest
		err = c.ShouldBindBodyWithJSON(&req)
		logger.Info("unmarshal body")
		if err != nil {
			logger.Error("cannot unmarshall body", slog.String("err", err.Error()))
			logMsg, respMsg, code := handleErrorWhileEditTender(err)
			logger.Error(logMsg, slog.String("err", err.Error()))
			c.JSON(code, CreateTenderResponse{Message: respMsg})
			return
		}
		logger.Info("success unmarhsall", slog.String("new tender data", fmt.Sprintf("%+v", req)))

		editedTender, err := tenderEditor.EditTender(ctx, tenderId, models.Tender{
			TenderName:      req.TenderName,
			Description:     req.Description,
			ServiceType:     req.ServiceType,
			Status:          req.Status,
			OrganizationId:  req.OrganizationId,
			CreatorUsername: req.CreatorUsername,
		})
		if err != nil {
			logMsg, respMsg, code := handleErrorWhileUnmarshallCreateTenderRequest(err)
			logger.Error(logMsg, slog.String("err", err.Error()))
			c.JSON(code, EditTenderResponse{Message: respMsg})
			return
		}
		logger.Info("edit success")
		c.JSON(http.StatusOK, EditTenderResponse{Message: "ok", UpdatedTender: editedTender})
	}
}

func handleErrorWhileUnmarshallCreateTenderRequest(err error) (logMessage string, reponseMessage string, code int) {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	switch {
	case errors.As(err, &syntaxError):
		return fmt.Sprintf("JSON syntax error at byte %d", syntaxError.Offset), "syntax error", http.StatusBadRequest
	case errors.As(err, &unmarshalTypeError):
		return fmt.Sprintf("JSON type error: field '%s' expects %s but got %s",
			unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value), "JSON type error", http.StatusBadRequest
	case errors.Is(err, io.EOF):
		return "empty json body", "empty json body", http.StatusBadRequest
	default:
		return "smth went wrong", "internal error", http.StatusInternalServerError
	}
}

func handleErrorWhileEditTender(err error) (logMessage string, reponseMessage string, code int) {
	if errors.Is(err, storage.ErrUserNotFound) {
		return "new user not exists", "new user not found", http.StatusBadRequest
	} else if errors.Is(err, storage.ErrOrganizationNotFound) {
		return "new organization not found", "new organization not found", http.StatusBadRequest
	} else if errors.Is(err, storage.ErrUserNotReponsibleForOrg) {
		return "new user not respobsible for organization", "new user not respobsible for organization", http.StatusBadRequest
	}
	return "unexpected error", "internal error", http.StatusInternalServerError

}

func validateCreateTenderRequest(req CreateTenderRequest) error {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return err
	}

	return err
}
