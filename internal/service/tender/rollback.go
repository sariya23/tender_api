package tender

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sariya23/tender/internal/domain/models"
	outerror "github.com/sariya23/tender/internal/out_error"
)

func (tenderSrv *TenderService) RollbackTender(ctx context.Context, tenderId int, version int, username string) (models.Tender, error) {
	const operationPlace = "internal.service.tender.rollback.RollbackTender"
	logger := tenderSrv.logger.With("op", operationPlace)

	tender, err := tenderSrv.tenderRepo.GetTenderById(ctx, tenderId)
	if err != nil {
		if errors.Is(err, outerror.ErrTenderNotFound) {
			logger.Warn(fmt.Sprintf("tender with id=\"%d\" not found", tenderId))
			return models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
		}
		logger.Error(fmt.Sprintf("cannot get tender with id=\"%d\"", tenderId))
		return models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
	}

	if tender.CreatorUsername != username {
		logger.Warn(fmt.Sprintf("employee with username=<%s> not responsible for tender with id=<%d>", tender.CreatorUsername, tenderId))
		return models.Tender{}, fmt.Errorf("%s: %w", operationPlace, outerror.ErrEmployeeNotResponsibleForTender)
	}

	err = tenderSrv.tenderRepo.FindTenderVersion(ctx, tenderId, version)
	if err != nil {
		if errors.Is(err, outerror.ErrTenderVersionNotFound) {
			logger.Warn(fmt.Sprintf("tender version=\"%d\" not found", version))
			return models.Tender{}, err
		}
		logger.Error("cannot get tender version")
		return models.Tender{}, err
	}

	err = tenderSrv.tenderRepo.RollbackTender(ctx, tenderId, version)
	if err != nil {
		logger.Error("cannot rollback tender", slog.String("err", err.Error()))
		return models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
	}

	tender, err = tenderSrv.tenderRepo.GetTenderById(ctx, tenderId)
	if err != nil {
		logger.Error("unexpected error", slog.String("err", err.Error()))
		return models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
	}
	return tender, nil
}
