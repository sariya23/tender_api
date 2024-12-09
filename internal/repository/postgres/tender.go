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
	values ($1, $2, $3, $4, $5, $6) returning name, description, service_type, organization_id, creator_username, status`
	createdTender = models.Tender{}

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
		tender.Status,
		tender.OrganizationId,
		tender.CreatorUsername,
	)
	err = row.Scan(
		&createdTender.TenderName,
		&createdTender.Description,
		&createdTender.ServiceType,
		&createdTender.OrganizationId,
		&createdTender.CreatorUsername,
		&createdTender.Status,
	)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w. Place = createQuery", op, err)
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
func (s *Storage) GetEmployeeTenders(ctx context.Context, empl models.Employee) ([]models.Tender, error) {
	const op = "repository.postgres.tender.GetEmployeeTenders"
	query := `select name, description, service_type, status, organization_id, creator_username 
	from tender
	where creator_username = $1`
	tenders := []models.Tender{}

	rows, err := s.connection.Query(ctx, query, empl.Username)
	if err != nil {
		return []models.Tender{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		tender := models.Tender{}
		err := rows.Scan(
			&tender.TenderName,
			&tender.Description,
			&tender.ServiceType,
			&tender.Status,
			&tender.OrganizationId,
			&tender.CreatorUsername,
		)
		if err != nil {
			return []models.Tender{}, fmt.Errorf("%s: %w", op, err)
		}

		tenders = append(tenders, tender)

		if err := rows.Err(); err != nil {
			return []models.Tender{}, fmt.Errorf("%s: %w", op, err)
		}
	}
	if len(tenders) == 0 {
		return []models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrEmployeeTendersNotFound)
	}
	return tenders, nil
}
func (s *Storage) EditTender(ctx context.Context, tenderId int, updateTender models.TenderToUpdate) (models.Tender, error) {
	const op = "repository.postgres.tender.EditTender"
	// query := "update %s set %s where %s = @%s"
	updateParams := make([]string, 0)
	updateFormat := "%s = @%s"
	fieldsToUpdate := pgx.NamedArgs{}

	if name := updateTender.TenderName; name != nil {
		updateParams = append(updateParams, fmt.Sprintf(updateFormat, "name", "name"))
		fieldsToUpdate["name"] = *name
	}
	if desc := updateTender.Description; desc != nil {
		updateParams = append(updateParams, fmt.Sprintf(updateFormat, "description", "description"))
		fieldsToUpdate["description"] = *desc
	}
	if srvType := updateTender.ServiceType; srvType != nil {
		updateParams = append(updateParams, fmt.Sprintf(updateFormat, "service_type", "service_type"))
		fieldsToUpdate["service_type"] = *srvType
	}
	if status := updateTender.Status; status != nil {
		updateParams = append(updateParams, fmt.Sprintf(updateFormat, "status", "status"))
		fieldsToUpdate["status"] = *status
	}
	if orgId := updateTender.OrganizationId; orgId != nil {
		updateParams = append(updateParams, fmt.Sprintf(updateFormat, "organization_id", "organization_id"))
		fieldsToUpdate["organization_id"] = *orgId
	}
	if username := updateTender.CreatorUsername; username != nil {
		updateParams = append(updateParams, fmt.Sprintf(updateFormat, "creator_username", "creator_username"))
		fieldsToUpdate["creator_username"] = *username
	}

	if len(fieldsToUpdate) == 0 {
		return models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrNothingToUpdate)
	}

	return models.Tender{}, nil
}
func (s *Storage) RollbackTender(ctx context.Context, tenderId int, toVersionRollback int) (models.Tender, error) {
	panic("impl me")
}
func (s *Storage) GetTenderById(ctx context.Context, tenderId int) (models.Tender, error) {
	const op = "repository.postgres.tender.GetTenderById"
	query := `select name, description, service_type, status, organization_id, creator_username 
	from tender
	where tender_id = $1`

	var tender models.Tender

	row := s.connection.QueryRow(ctx, query, tenderId)
	err := row.Scan(
		&tender.TenderName,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationId,
		&tender.CreatorUsername,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrTenderNotFound)
		} else {
			return models.Tender{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	return tender, nil
}

func (s *Storage) FindTenderVersion(ctx context.Context, tenderId int, version int) error {
	panic("impl me")
}

func (s *Storage) GetTenderStatus(ctx context.Context, tenderStatus string) (string, error) {
	panic("impl me")
}
