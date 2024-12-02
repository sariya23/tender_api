package repository

import (
	"context"

	"github.com/sariya23/tender/internal/domain/models"
)

type TenderRepository interface {
	CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error)
	GetAllTenders(ctx context.Context) ([]models.Tender, error)
	GetTendersByServiceType(ctx context.Context, serviceType string) ([]models.Tender, error)
	GetEmployeeTendersByUsername(ctx context.Context, username string) ([]models.Tender, error)
	EditTender(ctx context.Context, tenderId int, updateTender models.TenderToUpdate) (models.TenderToUpdate, error)
	RollbackTender(ctx context.Context, tenderId int, toVersionRollback int) (models.Tender, error)
	GetTenderById(ctx context.Context, tenderId int) (models.Tender, error)
}

type EmployeeRepository interface {
	GetEmployeeByUsername(ctx context.Context, username string) (models.Employee, error)
	GetEmployeeById(ctx context.Context, id int) (models.Employee, error)
}

type OrganizationRepository interface {
	GetOrganizationById(ctx context.Context, orgId int) (models.Organization, error)
}

type EmployeeResponsibler interface {
	CheckResponsibility(ctx context.Context, emplId int, orgId int) error
}
