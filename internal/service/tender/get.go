package tender

import (
	"context"
	"fmt"
	"log/slog"
	"tender/internal/domain/models"
)

// GetTenders...
func (s *TenderService) GetTenders(ctx context.Context, serviceType string) ([]models.Tender, error) {
	const op = "internal.service.tender.getall.GetTenders"
	logger := s.logger.With("op", op)

	var err error
	var tenders []models.Tender

	if serviceType == "all" {
		logger.Info("get all tenders")
		tenders, err = s.tenderRepo.GetAll(ctx)
	} else {
		logger.Info("get tenders with service type", slog.String("service type", serviceType))
		tenders, err = s.tenderRepo.GetByServiceType(ctx, serviceType)
	}

	if err != nil {
		logger.Error("cannot get tenders", slog.String("err", err.Error()))
		return []models.Tender{}, fmt.Errorf("cannot get tenders: %w", err)
	}
	return tenders, nil
}
