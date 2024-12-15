package tenderapi

import (
	"log/slog"

	"github.com/sariya23/tender/internal/service"
)

type TenderService struct {
	logger        *slog.Logger
	tenderService service.TenderServiceProvider
}

func New(logger *slog.Logger, tenderService service.TenderServiceProvider) *TenderService {
	return &TenderService{
		logger:        logger,
		tenderService: tenderService,
	}
}
