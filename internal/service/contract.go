package service

import (
	"context"

	"github.com/sariya23/tender/internal/domain/models"
)

// TenderServiceProvider представляет набор методов, которые нужны
// для взаимодействия с тендерами.
type TenderServiceProvider interface {
	CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error)
	GetTenders(ctx context.Context, serviceType string) ([]models.Tender, error)
	GetEmployeeTendersByUsername(ctx context.Context, username string) ([]models.Tender, error)
	Edit(ctx context.Context, tenderId int, updateTender models.TenderToUpdate) (models.TenderToUpdate, error)
}
