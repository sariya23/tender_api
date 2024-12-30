package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sariya23/tender/internal/domain/models"
	outerror "github.com/sariya23/tender/internal/out_error"
)

func (storage *Storage) GetOrganizationById(ctx context.Context, orgId int) (models.Organization, error) {
	const operationPlace = "repository.postgres.organization.GetOrganizationById"
	queryGetOrg := "select organization_id, name, description, organization_type_id from organization where organization_id = $1"
	queryGetOrgType := "select type from nsi_organization_type where nsi_organization_type_id = $1"

	var organization models.Organization
	var typeId int

	row := storage.connection.QueryRow(ctx, queryGetOrg, orgId)
	err := row.Scan(&organization.ID, &organization.Name, &organization.Description, &typeId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Organization{}, fmt.Errorf("%s: %w", operationPlace, outerror.ErrOrganizationNotFound)
		} else {
			return models.Organization{}, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}

	typeRow := storage.connection.QueryRow(ctx, queryGetOrgType, typeId)
	err = typeRow.Scan(&organization.Type)
	if err != nil {
		return models.Organization{}, fmt.Errorf("%s: %w", operationPlace, err)
	}

	return organization, nil
}

func (storage *Storage) CheckResponsibility(ctx context.Context, emplId int, orgId int) error {
	const operationPlace = "repository.postgres.organization.CheckResponsibility"
	query := "select organization_responsible_id from organization_responsible where organization_id = $1 and employee_id = $2"

	var id int
	row := storage.connection.QueryRow(ctx, query, orgId, emplId)
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", operationPlace, outerror.ErrEmployeeNotResponsibleForOrganization)
		} else {
			return fmt.Errorf("%s: %w", operationPlace, err)
		}
	}

	return nil
}

func (storage *Storage) CreateOrganization(ctx context.Context, organization models.Organization) error {
	const operationPlace = "repository.postgres.organization.CreateOrganization"
	queryGetOrgType := "select type_id from nsi_organization_type where type = $1"
	insertOrg := "insert into otganization(name, description, organization_type_id) values (@name, @desc, @orgTypeId)"
	var typeId int

	row := storage.connection.QueryRow(ctx, queryGetOrgType, organization.Type)
	err := row.Scan(&typeId)
	if err != nil {
		return fmt.Errorf("%s: %w", operationPlace, err)
	}

	_, err = storage.connection.Exec(
		ctx,
		insertOrg,
		pgx.NamedArgs{
			"name":      organization.Name,
			"desc":      organization.Description,
			"orgTypeId": typeId,
		},
	)
	if err != nil {
		return fmt.Errorf("%s: %w", operationPlace, err)
	}

	return nil
}
