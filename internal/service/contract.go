package service

import (
	"context"
	"tender/internal/domain/models"
)

type TenderServiceProvider interface {
	CreateTender(ctx context.Context, tender models.Tender)
	GetTenders(ctx context.Context, serviceType string) ([]models.Tender, error)
	GetUserTenders(ctx context.Context, username string) ([]models.Tender, error)
	Edit(ctx context.Context, tenderId int, updateTender models.TenderToUpdate) (models.TenderToUpdate, error)
}
