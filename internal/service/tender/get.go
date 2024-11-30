package tender

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sariya23/tender/internal/domain/models"
)

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
		logger.Error("cannot get tenders", slog.String("err", err.Error()))
		return []models.Tender{}, fmt.Errorf("cannot get tenders: %w", err)
	}
	logger.Info("success get tenders")
	return tenders, nil
}

func (s *TenderService) GetUserTenders(ctx context.Context, username string) ([]models.Tender, error) {
	const op = "internal.service.tender.getall.GetUserTenders"
	logger := s.logger.With("op", op)

	_, err := s.employeeRepo.GetEmployeeByUsername(ctx, username)
	if err != nil {
		logger.Error("cannot get user with username", slog.String("username", username), slog.String("err", err.Error()))
		return []models.Tender{}, err
	}
	logger.Info("success check employee by username")
	tenders, err := s.tenderRepo.GetUserTenders(ctx, username)
	if err != nil {
		logger.Error("cannot get tenders", slog.String("err", err.Error()))
		return []models.Tender{}, err
	}
	logger.Info("success get employee tenders")
	return tenders, nil
}
