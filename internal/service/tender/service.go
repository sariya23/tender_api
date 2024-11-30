package tender

import (
	"log/slog"
	"tender/internal/repository"
)

type TenderService struct {
	logger                  *slog.Logger
	tenderRepo              repository.TenderRepository
	employeeRepo            repository.EmployeeRepository
	orgRepo                 repository.OrganizationRepository
	employeeOrgResponsibler repository.EmployeeResponsibler
}

func New(
	logger *slog.Logger,
	tenderRepo repository.TenderRepository,
	employeeRepo repository.EmployeeRepository,
	orgRepo repository.OrganizationRepository,
	employeeOrgResponsibler repository.EmployeeResponsibler,
) *TenderService {
	return &TenderService{
		logger:                  logger,
		tenderRepo:              tenderRepo,
		employeeRepo:            employeeRepo,
		orgRepo:                 orgRepo,
		employeeOrgResponsibler: employeeOrgResponsibler,
	}
}
