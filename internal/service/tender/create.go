package tender

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sariya23/tender/internal/domain/models"
	outerror "github.com/sariya23/tender/internal/out_error"
)

// CreateTender создает тендер с данными, переданными в tender.
func (s *TenderService) CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	const op = "internal.service.tender.create.CreateTender"
	logger := s.logger.With("op", op)

	empl, err := s.employeeRepo.GetEmployeeByUsername(ctx, tender.CreatorUsername)
	if err != nil {
		if errors.Is(err, outerror.ErrEmployeeNotFound) {
			logger.Warn("employee not found", slog.String("username", tender.CreatorUsername))
			return models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrEmployeeNotFound)
		}
		logger.Error("cannot get employee with username", slog.String("username", tender.CreatorUsername), slog.String("err", err.Error()))
		return models.Tender{}, fmt.Errorf("cannot get employee: %w", err)
	}
	logger.Info("success check employee by username")

	_, err = s.orgRepo.GetOrganizationById(ctx, tender.OrganizationId)
	if err != nil {
		if errors.Is(err, outerror.ErrOrganizationNotFound) {
			logger.Warn("organization not found", slog.Int("org id", tender.OrganizationId))
			return models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrOrganizationNotFound)
		}
		logger.Error("cannot get organization with id", slog.Int("org id", tender.OrganizationId), slog.String("err", err.Error()))
		return models.Tender{}, fmt.Errorf("cannot get organization with id: %w", err)
	}
	logger.Info("success check organization by id")
	err = s.employeeResponsibler.CheckResponsibility(ctx, empl.ID, tender.OrganizationId)
	if err != nil {
		if errors.Is(err, outerror.ErrEmployeeNotResponsibleForOrganization) {
			logger.Warn("employee not responsible for organization", slog.Int("empl id", empl.ID), slog.Int("org id", tender.OrganizationId))
			return models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrEmployeeNotResponsibleForOrganization)
		}
		logger.Error(
			"cannot check that employee responsible for organization",
			slog.Int("empl id", empl.ID),
			slog.Int("org id", tender.OrganizationId),
			slog.String("err", err.Error()),
		)
		return models.Tender{}, fmt.Errorf("cannot check that employee responsible for organization: %w", err)
	}
	logger.Info("success check employee responsible")
	createdTender, err := s.tenderRepo.CreateTender(ctx, tender)
	if err != nil {
		logger.Error("cannot create tender", slog.String("err", err.Error()))
		return models.Tender{}, fmt.Errorf("cannot create tender: %w", err)
	}
	logger.Info("success create tender")
	return createdTender, nil
}
