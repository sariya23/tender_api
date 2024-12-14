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

	lastTenderId, err := s.GetLastInsertedTenderId(ctx)
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

	query := `select name, description, service_type, status, organization_id, creator_username from tender
	where selected_version = $1`
	tenders := []models.Tender{}

	rows, err := s.connection.Query(ctx, query, true)
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
	where service_type=$1 and selected_version=$2`
	tenders := []models.Tender{}

	rows, err := s.connection.Query(ctx, query, serviceType, true)

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
	where creator_username = $1 and selected_version=$2`
	tenders := []models.Tender{}

	rows, err := s.connection.Query(ctx, query, empl.Username, true)
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
func (s *Storage) EditTender(
	ctx context.Context,
	oldTender models.Tender,
	tenderId int,
	updateTender models.TenderToUpdate,
) (models.Tender, error) {
	const op = "repository.postgres.tender.EditTender"

	insertQuery := `
	insert into tender values (@name, @desc, @srv_type, @status, @org_id, @username, @version, @selected_version)
	returning name, description, service_type, status, organization_id, creator_username`

	args := pgx.NamedArgs{"selected_version": true}

	if newName := updateTender.TenderName; newName == nil {
		args["name"] = oldTender.TenderName
	} else {
		args["name"] = newName
	}

	if newDesc := updateTender.Description; newDesc == nil {
		args["desc"] = oldTender.Description
	} else {
		args["desc"] = newDesc
	}

	if newSrvType := updateTender.ServiceType; newSrvType == nil {
		args["srv_type"] = oldTender.ServiceType
	} else {
		args["srv_type"] = newSrvType
	}

	if newStatus := updateTender.Status; newStatus == nil {
		args["status"] = oldTender.Status
	} else {
		args["status"] = newStatus
	}

	if newOrgId := updateTender.OrganizationId; newOrgId == nil {
		args["org_id"] = oldTender.OrganizationId
	} else {
		args["org_id"] = newOrgId
	}

	if newUsername := updateTender.CreatorUsername; newUsername == nil {
		args["username"] = oldTender.CreatorUsername
	} else {
		args["username"] = newUsername
	}

	err := s.deactivateTenderVersion(ctx, tenderId)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", op, err)
	}
	var tender models.Tender

	row := s.connection.QueryRow(ctx, insertQuery, args)
	err = row.Scan(
		&tender.TenderName,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationId,
		&tender.CreatorUsername,
	)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", op, err)
	}

	return tender, nil
}
func (s *Storage) RollbackTender(ctx context.Context, tenderId int, toVersionRollback int) (models.Tender, error) {
	panic("impl me")
}
func (s *Storage) GetTenderById(ctx context.Context, tenderId int) (models.Tender, error) {
	const op = "repository.postgres.tender.GetTenderById"
	query := `select name, description, service_type, status, organization_id, creator_username 
	from tender
	where tender_id = $1 and selected_version=$2`

	var tender models.Tender

	row := s.connection.QueryRow(ctx, query, tenderId, true)
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

func (s *Storage) GetLastInsertedTenderId(ctx context.Context) (int, error) {
	const op = "repository.postgres.tender.getLastInsertedTenderId"
	query := `select tender_id from tender order by tender_id desc limit 1`
	row := s.connection.QueryRow(ctx, query)
	var tenderId int
	err := row.Scan(&tenderId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return tenderId, nil
}

func (s *Storage) getLastTenderVersion(ctx context.Context, tenderId int) (int, error) {
	const op = "repository.postgres.tender.getLastTenderVersion"
	query := "select version from tender where tender_id = $1"
	row := s.connection.QueryRow(ctx, query, tenderId)
	var version int
	err := row.Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return version, nil
}

func (s *Storage) deactivateTenderVersion(ctx context.Context, tenderId int) error {
	const op = "repository.postgres.tender.deactivateTenderVersion"
	query := "update tender set selected_version = $1 where tender_id = $2"
	_, err := s.connection.Exec(ctx, query, false, tenderId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) activateTenderVersion(ctx context.Context, tenderId int, version int) error {
	const op = "repository.postgres.tender.activateTenderVersion"
	query := "update tender set selected_version = $1 where tender_id = $2 and version = $3"
	_, err := s.connection.Exec(ctx, query, true, tenderId, version)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
