package repository

import (
	"context"
	"tender/internal/domain/models"
)

type TenderRepository interface {
	Create(ctx context.Context, tender models.Tender) (models.Tender, error)
	GetAll(ctx context.Context) ([]models.Tender, error)
	GetByServiceType(ctx context.Context, serviceType string) ([]models.Tender, error)
	GetUserTenders(ctx context.Context, username string) ([]models.Tender, error)
	Edit(ctx context.Context, updateTender models.Tender) (models.Tender, error)
	Rollback(ctx context.Context, tenderId int, toVersionRollback int) (models.Tender, error)
}

type EmployeeRepository interface {
	GetByUsername(ctx context.Context, username string) (models.Employee, error)
	GetById(ctx context.Context, id int) (models.Employee, error)
}

type OrganizationRepository interface {
	GetById(ctx context.Context, orgId int) (models.Organization, error)
}

type EmployeeResponsibler interface {
	CheckResponsibility(ctx context.Context, emplId int, orgId int) error
}
