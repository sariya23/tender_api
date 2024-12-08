package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sariya23/tender/internal/domain/models"
	outerror "github.com/sariya23/tender/internal/out_error"
)

func (s *Storage) CreateTender(ctx context.Context, tender models.Tender) (createdTender models.Tender, err error) {
	const op = "repository.postgres.tender.CreateTender"

	createQuery := `insert into tender(name, description, service_type, status, organization_id, creator_username)
	values ($1, $2, $3, $4, $5, $6) returning name, description, service_type, status, organization_id, creator_username`
	getEmployeeQuery := `select employee_id from employee where username = $1`
	getOrgQuery := `select organization_id from organization where name = $1`
	checkResponsobilityQuery := `select organization_responsible_id from organization_responsible where organization_id = $1 and employee_id = $2`
	getStatusIdQuery := `select nsi_tender_status_id from nsi_tender_status where status = $1`
	createdTender = models.Tender{}

	var employeeId int
	rowEmployeeId := s.connection.QueryRow(ctx, getEmployeeQuery, tender.CreatorUsername)
	err = rowEmployeeId.Scan(&employeeId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrEmployeeNotFound)
		} else {
			return models.Tender{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	var orgId int
	rowOrgId := s.connection.QueryRow(ctx, getOrgQuery, tender.OrganizationId)
	err = rowOrgId.Scan(&orgId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrOrganizationNotFound)
		} else {
			return models.Tender{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	var responsibleId int
	rowResponsibleId := s.connection.QueryRow(ctx, checkResponsobilityQuery, orgId, employeeId)
	err = rowResponsibleId.Scan(&responsibleId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrEmployeeNotResponsibleForOrganization)
		} else {
			return models.Tender{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	var statusId int
	rowStatusId := s.connection.QueryRow(ctx, getStatusIdQuery, tender.Status)
	err = rowStatusId.Scan(&statusId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrUnknownTenderStatus)
		} else {
			return models.Tender{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	tx, err := s.connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	row := tx.QueryRow(
		ctx,
		createQuery,
		tender.TenderName,
		tender.Description,
		tender.ServiceType,
		statusId,
		tender.OrganizationId,
		tender.CreatorUsername,
	)
	err = row.Scan(
		&createdTender.TenderName,
		&createdTender.Description,
		&createdTender.ServiceType,
		&createdTender.Status,
		&createdTender.OrganizationId,
		&createdTender.CreatorUsername,
	)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", op, err)
	}
	return createdTender, nil
}

func (s *Storage) GetAllTenders(ctx context.Context) ([]models.Tender, error) {
	const op = "repository.postgres.tender.GetAllTenders"

	query := `select name, description, service_type, status, organization_id, creator_username from tender`
	tenders := []models.Tender{}

	rows, err := s.connection.Query(ctx, query)
	if err != nil {
		return []models.Tender{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		tender := models.Tender{}
		err := rows.Scan(&tender.TenderName, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatorUsername)
		if err != nil {
			return []models.Tender{}, fmt.Errorf("%s: %w", op, err)
		}

		tenders = append(tenders, tender)

		if err := rows.Err(); err != nil {
			return []models.Tender{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	if len(tenders) == 0 {
		return []models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrTendersWithThisServiceTypeNotFound)
	}

	return tenders, nil
}

func (s *Storage) GetTendersByServiceType(ctx context.Context, serviceType string) ([]models.Tender, error) {
	const op = "repository.postgres.tender.GetAllTenders"

	query := `select name, description, service_type, status, organization_id, creator_username
	from tender
	where service_type=$1`
	tenders := []models.Tender{}

	rows, err := s.connection.Query(ctx, query, serviceType)

	if err != nil {
		return []models.Tender{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		tender := models.Tender{}
		err := rows.Scan(&tender.TenderName, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatorUsername)
		if err != nil {
			return []models.Tender{}, fmt.Errorf("%s: %w", op, err)
		}

		tenders = append(tenders, tender)

		if err := rows.Err(); err != nil {
			return []models.Tender{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	if len(tenders) == 0 {
		return []models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrTendersWithThisServiceTypeNotFound)
	}

	return tenders, nil
}
func (s *Storage) GetEmployeeTendersByUsername(ctx context.Context, username string) ([]models.Tender, error) {
	panic("impl me")
}
func (s *Storage) EditTender(ctx context.Context, tenderId int, updateTender models.TenderToUpdate) (models.Tender, error) {
	panic("impl me")
}
func (s *Storage) RollbackTender(ctx context.Context, tenderId int, toVersionRollback int) (models.Tender, error) {
	panic("impl me")
}
func (s *Storage) GetTenderById(ctx context.Context, tenderId int) (models.Tender, error) {
	panic("impl me")
}

func (s *Storage) FindTenderVersion(ctx context.Context, tenderId int, version int) error {
	panic("impl me")
}
