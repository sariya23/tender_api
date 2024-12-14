package tenderapp

import (
	"log/slog"

	tenderapi "github.com/sariya23/tender/internal/api/tender"
	"github.com/sariya23/tender/internal/repository"
	tendersrv "github.com/sariya23/tender/internal/service/tender"
)

type TenderApp struct {
	TenderHandlers *tenderapi.TenderService
}

func New(
	logger *slog.Logger,
	tenderRepo repository.TenderRepository,
	employeeRepo repository.EmployeeRepository,
	orgRepo repository.OrganizationRepository,
	responsibler repository.EmployeeResponsibler,
) *TenderApp {
	tenderService := tendersrv.New(logger, tenderRepo, employeeRepo, orgRepo, responsibler)
	tenderHandlers := tenderapi.New(logger, tenderService)
	return &TenderApp{tenderHandlers}
}
