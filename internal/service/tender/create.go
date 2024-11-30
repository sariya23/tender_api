package tender

import (
	"context"
	"log/slog"

	"github.com/sariya23/tender/internal/domain/models"
)

// CreateTender создает тендер с данными, переданными в tender.
func (s *TenderService) CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	const op = "internal.service.tender.create.CreateTender"
	logger := s.logger.With("op", op)

	empl, err := s.employeeRepo.GetEmployeeByUsername(ctx, tender.CreatorUsername)
	if err != nil {
		logger.Error("cannot get user with username", slog.String("username", tender.CreatorUsername), slog.String("err", err.Error()))
		return models.Tender{}, err
	}
	logger.Info("success check employee by username")

	_, err = s.orgRepo.GetOrganizationById(ctx, tender.OrganizationId)
	if err != nil {
		logger.Error("cannot get organization with id", slog.Int("org is", tender.OrganizationId), slog.String("err", err.Error()))
		return models.Tender{}, err
	}
	logger.Info("success check organization by id")
	err = s.employeeResponsibler.CheckResponsibility(ctx, empl.ID, tender.OrganizationId)
	if err != nil {
		logger.Error(
			"cannot check that employee responble for organization",
			slog.Int("empl id", empl.ID),
			slog.Int("org id", tender.OrganizationId),
			slog.String("err", err.Error()),
		)
		return models.Tender{}, err
	}
	logger.Info("success check employee responsible")
	createdTender, err := s.tenderRepo.CreateTender(ctx, tender)
	if err != nil {
		logger.Error("cannot create tender", slog.String("err", err.Error()))
		return models.Tender{}, err
	}
	logger.Info("success create tender")
	return createdTender, nil
}
