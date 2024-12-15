package tenderapi

import (
	"context"
	"log/slog"

	"github.com/sariya23/tender/internal/domain/models"
)

type TenderServiceProvider interface {
	CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error)
	GetTenders(ctx context.Context, serviceType string) ([]models.Tender, error)
	GetEmployeeTendersByUsername(ctx context.Context, username string) ([]models.Tender, error)
	EditTender(ctx context.Context, tenderId int, updateTender models.TenderToUpdate, username string) (models.Tender, error)
	RollbackTender(ctx context.Context, tenderId int, version int, username string) (models.Tender, error)
}

type TenderService struct {
	logger        *slog.Logger
	tenderService TenderServiceProvider
}

func New(logger *slog.Logger, tenderService TenderServiceProvider) *TenderService {
	return &TenderService{
		logger:        logger,
		tenderService: tenderService,
	}
}
