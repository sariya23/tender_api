package tender

import (
	"context"
	"log/slog"
	"tender/internal/domain/models"
)

func (s *TenderService) CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	const op = "internal.service.tender.create.CreateTender"
	logger := s.logger.With("op", op)

	_, err := s.employeeRepo.GetByUsername(ctx, tender.CreatorUsername)
	if err != nil {
		logger.Error("cannot get user with username", slog.String("username", tender.CreatorUsername), slog.String("err", err.Error()))
		return models.Tender{}, err
	}

	_, err = s.orgRepo.GetById(ctx, tender.OrganizationId)
	if err != nil {
		logger.Error("cannot get organization with id", slog.Int("org is", tender.OrganizationId), slog.String("err", err.Error()))
		return models.Tender{}, err
	}

	createdTender, err := s.tenderRepo.Create(ctx, tender)
	if err != nil {
		logger.Error("cannot create tender", slog.String("err", err.Error()))
		return models.Tender{}, err
	}
	return createdTender, nil
}
