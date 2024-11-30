package postgres

import (
	"context"

	"github.com/sariya23/tender/internal/domain/models"
)

func (s Storage) GetOrganizationById(ctx context.Context, orgId int) (models.Organization, error) {
	panic("impl me")
}

func (s Storage) CheckResponsibility(ctx context.Context, emplId int, orgId int) error {
	panic("impl me")
}
