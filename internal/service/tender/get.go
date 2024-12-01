package tender

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sariya23/tender/internal/domain/models"
	"github.com/sariya23/tender/internal/repository"
	"github.com/sariya23/tender/internal/service"
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
		if errors.Is(err, repository.ErrNoTendersWithThisServiceType) {
			logger.Warn("no tenders found", slog.String("err", err.Error()))
			return []models.Tender{}, fmt.Errorf("no tenders found: %w", service.ErrNoTendersFound)
		}
		logger.Error("cannot get tenders", slog.String("err", err.Error()))
		return []models.Tender{}, fmt.Errorf("cannot get tenders: %w", err)
	}
	logger.Info("success get tenders")
	return tenders, nil
}

// GetUserTenders возвращает список тендоров, которые связаны с переданным юзером.
func (s *TenderService) GetUserTenders(ctx context.Context, username string) ([]models.Tender, error) {
	const op = "internal.service.tender.getall.GetUserTenders"
	logger := s.logger.With("op", op)

	_, err := s.employeeRepo.GetEmployeeByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrEmployeeNotFound) {
			logger.Warn("employee not found", slog.String("username", username))
			return []models.Tender{}, service.ErrEmployeeNotFound
		}
		logger.Error("cannot get user", slog.String("username", username), slog.String("err", err.Error()))
		return []models.Tender{}, err
	}
	logger.Info("success check employee by username")
	tenders, err := s.tenderRepo.GetUserTenders(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrNoUserTenders) {
			logger.Warn("no tenders for user", slog.String("username", username), slog.String("err", err.Error()))
			return []models.Tender{}, service.ErrUserTendersNotFound
		}
		logger.Error("cannot get tenders", slog.String("err", err.Error()))
		return []models.Tender{}, err
	}
	logger.Info("success get employee tenders")
	return tenders, nil
}
