package tender

import (
	"log/slog"

	"github.com/sariya23/tender/internal/repository"
)

type TenderService struct {
	logger               *slog.Logger
	tenderRepo           repository.TenderRepository
	employeeRepo         repository.EmployeeRepository
	orgRepo              repository.OrganizationRepository
	employeeResponsibler repository.EmployeeResponsibler
}

func New(
	logger *slog.Logger,
	tenderRepo repository.TenderRepository,
	employeeRepo repository.EmployeeRepository,
	orgRepo repository.OrganizationRepository,
	employeeOrgResponsibler repository.EmployeeResponsibler,
) *TenderService {
	return &TenderService{
		logger:               logger,
		tenderRepo:           tenderRepo,
		employeeRepo:         employeeRepo,
		orgRepo:              orgRepo,
		employeeResponsibler: employeeOrgResponsibler,
	}
}
