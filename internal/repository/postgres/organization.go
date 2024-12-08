package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sariya23/tender/internal/domain/models"
	outerror "github.com/sariya23/tender/internal/out_error"
)

func (s Storage) GetOrganizationById(ctx context.Context, orgId int) (models.Organization, error) {
	const op = "repository.postgres.employee.GetOrganizationById"
	queryGetOrg := "select organization_id, name, description, organization_type_id from organization where organization_id = $1"
	queryGetOrgType := "select type from nsi_organization_type where nsi_organization_type_id = $1"

	var organization models.Organization
	var typeId int

	row := s.connection.QueryRow(ctx, queryGetOrg, orgId)
	err := row.Scan(&organization.ID, &organization.Name, &organization.Description, &typeId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Organization{}, fmt.Errorf("%s: %w", op, outerror.ErrOrganizationNotFound)
		} else {
			return models.Organization{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	typeRow := s.connection.QueryRow(ctx, queryGetOrgType, typeId)
	err = typeRow.Scan(&organization.Type)
	if err != nil {
		return models.Organization{}, fmt.Errorf("%s: %w", op, err)
	}

	return organization, nil
}

func (s Storage) CheckResponsibility(ctx context.Context, emplId int, orgId int) error {
	const op = "repository.postgres.employee.CheckResponsibility"
	query := "select organization_responsible_id from organization_responsible where organization_id = $1 and employee_id = $2"

	var id int
	row := s.connection.QueryRow(ctx, query, orgId, emplId)
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, outerror.ErrEmployeeNotResponsibleForOrganization)
		} else {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
