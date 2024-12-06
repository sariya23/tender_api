package tender

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sariya23/tender/internal/domain/models"
	outerror "github.com/sariya23/tender/internal/out_error"
)

func (s *TenderService) RollbackTender(ctx context.Context, tenderId int, version int) (models.Tender, error) {
	const op = "internal.service.tender.rollback.RollbackTender"
	logger := s.logger.With("op", op)

	_, err := s.tenderRepo.GetTenderById(ctx, tenderId)
	if err != nil {
		if errors.Is(err, outerror.ErrTenderNotFound) {
			logger.Warn(fmt.Sprintf("tender with id=\"%d\" not found", tenderId))
			return models.Tender{}, fmt.Errorf("%s: %w", op, err)
		}
		logger.Error(fmt.Sprintf("cannot get tender with id=\"%d\"", tenderId))
		return models.Tender{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.tenderRepo.FindTenderVersion(ctx, tenderId, version)
	if err != nil {
		if errors.Is(err, outerror.ErrTenderVersionNotFound) {
			logger.Warn(fmt.Sprintf("tender version=\"%d\" not found", version))
			return models.Tender{}, err
		}
		logger.Error("cannot get tender version")
		return models.Tender{}, err
	}

	tender, err := s.tenderRepo.RollbackTender(ctx, tenderId, version)
	if err != nil {
		logger.Error("cannot rollback tender", slog.String("err", err.Error()))
		return models.Tender{}, fmt.Errorf("%s: %w", op, err)
	}

	return tender, nil
}
