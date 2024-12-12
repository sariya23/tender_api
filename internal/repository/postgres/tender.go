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

	lastTenderId, err := s.getLastInsertedTenderId(ctx)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", op, err)
	}
	args := pgx.NamedArgs{
		"tender_id":    lastTenderId + 1,
		"name":         tender.TenderName,
		"desc":         tender.Description,
		"service_type": tender.ServiceType,
		"status":       tender.Status,
		"org_id":       tender.OrganizationId,
		"username":     tender.CreatorUsername,
		"version":      1,
	}
	createQuery := `insert into tender values (@tender_id, @name, @desc, @service_type, @status, @org_id, @username, @version) 
	returning name, description, service_type, organization_id, creator_username, status`
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
		args,
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
	where service_type=$1
	order by version desc
	limit 1`
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
	where creator_username = $1
	order by versrion desc
	limit 1`
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
	panic("impl me")
}
func (s *Storage) RollbackTender(ctx context.Context, tenderId int, toVersionRollback int) (models.Tender, error) {
	panic("impl me")
}
func (s *Storage) GetTenderById(ctx context.Context, tenderId int) (models.Tender, error) {
	const op = "repository.postgres.tender.GetTenderById"
	query := `select name, description, service_type, status, organization_id, creator_username 
	from tender
	where tender_id = $1
	order by version desc
	limit 1`

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

func (s *Storage) getLastInsertedTenderId(ctx context.Context) (int, error) {
	const op = "repository.postgres.tender.getLastInsertedTenderId"
	query := `select tender_id from tender order by version desc limit 1`
	row := s.connection.QueryRow(ctx, query)
	var tenderId int
	err := row.Scan(&tenderId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return tenderId, nil
}
