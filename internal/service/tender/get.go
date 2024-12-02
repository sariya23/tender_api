package tender

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sariya23/tender/internal/domain/models"
	outerror "github.com/sariya23/tender/internal/out_error"
)

// GetTenders возвращает список тендеров, который удовлетворяют переданному serviceType.
func (s *TenderService) GetTenders(ctx context.Context, serviceType string) ([]models.Tender, error) {
	const op = "internal.service.tender.getall.GetTenders"
	logger := s.logger.With("op", op)

	var err error
	var tenders []models.Tender

	if serviceType == "all" {
		logger.Info("get all tenders")
		tenders, err = s.tenderRepo.GetAllTenders(ctx)
	} else {
		logger.Info("get tenders with service type", slog.String("service type", serviceType))
		tenders, err = s.tenderRepo.GetTendersByServiceType(ctx, serviceType)
	}

	if err != nil {
		if errors.Is(err, outerror.ErrTendersWithThisServiceTypeNotFound) {
			logger.Warn("no tenders found", slog.String("err", err.Error()))
			return []models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrTendersWithThisServiceTypeNotFound)
		}
		logger.Error("cannot get tenders", slog.String("err", err.Error()))
		return []models.Tender{}, fmt.Errorf("cannot get tenders: %w", err)
	}
	logger.Info("success get tenders")
	return tenders, nil
}

// GetEmployeeTendersByUsername возвращает список тендоров, которые связаны с переданным юзером.
func (s *TenderService) GetEmployeeTendersByUsername(ctx context.Context, username string) ([]models.Tender, error) {
	const op = "internal.service.tender.getall.GetEmployeeTendersByUsername"
	logger := s.logger.With("op", op)

	_, err := s.employeeRepo.GetEmployeeByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, outerror.ErrEmployeeNotFound) {
			logger.Warn("employee not found", slog.String("username", username))
			return []models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrEmployeeNotFound)
		}
		logger.Error("cannot get employee", slog.String("username", username), slog.String("err", err.Error()))
		return []models.Tender{}, fmt.Errorf("cannot get employee: %w", err)
	}
	logger.Info("success check employee by username")
	tenders, err := s.tenderRepo.GetEmployeeTendersByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, outerror.ErrEmployeeTendersNotFound) {
			logger.Warn("no tenders for employee", slog.String("username", username), slog.String("err", err.Error()))
			return []models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrEmployeeTendersNotFound)
		}
		logger.Error("cannot get tenders", slog.String("err", err.Error()))
		return []models.Tender{}, fmt.Errorf("cannot get tenders: %w", err)
	}
	logger.Info("success get employee tenders")
	return tenders, nil
}
