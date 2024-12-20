package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sariya23/tender/internal/domain/models"
	outerror "github.com/sariya23/tender/internal/out_error"
)

func (storage *Storage) CreateTender(ctx context.Context, tender models.Tender) (createdTender models.Tender, err error) {
	const operationPlace = "repository.postgres.tender.CreateTender"

	lastTenderId, err := storage.GetLastInsertedTenderId(ctx)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
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
						returning name, description, service_type, organization_id, creator_username, status
	`
	createdTender = models.Tender{}

	tx, err := storage.connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
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
		return models.Tender{}, fmt.Errorf("%s: %w. Place = createQuery", operationPlace, err)
	}
	return createdTender, nil
}

func (storage *Storage) GetAllTenders(ctx context.Context) ([]models.Tender, error) {
	const operationPlace = "repository.postgres.tender.GetAllTenders"

	query := `select name, description, service_type, status, organization_id, creator_username from tender
				where is_active_version = $1 and status = $2
	`
	tenders := []models.Tender{}

	rows, err := storage.connection.Query(ctx, query, true, models.TenderPublishedStatus)
	if err != nil {
		return []models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer rows.Close()

	for rows.Next() {
		tender := models.Tender{}
		err := rows.Scan(&tender.TenderName, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatorUsername)
		if err != nil {
			return []models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
		}

		tenders = append(tenders, tender)

		if err := rows.Err(); err != nil {
			return []models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}

	if len(tenders) == 0 {
		return []models.Tender{}, fmt.Errorf("%s: %w", operationPlace, outerror.ErrTendersWithThisServiceTypeNotFound)
	}

	return tenders, nil
}

func (storage *Storage) GetTendersByServiceType(ctx context.Context, serviceType string) ([]models.Tender, error) {
	const operationPlace = "repository.postgres.tender.GetAllTenders"

	query := `select name, description, service_type, status, organization_id, creator_username
				from tender
				where service_type=$1 and is_active_version=$2 and status = $3`
	tenders := []models.Tender{}

	rows, err := storage.connection.Query(ctx, query, serviceType, true, models.TenderPublishedStatus)

	if err != nil {
		return []models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer rows.Close()

	for rows.Next() {
		tender := models.Tender{}
		err := rows.Scan(&tender.TenderName, &tender.Description, &tender.ServiceType, &tender.Status, &tender.OrganizationId, &tender.CreatorUsername)
		if err != nil {
			return []models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
		}

		tenders = append(tenders, tender)

		if err := rows.Err(); err != nil {
			return []models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}

	if len(tenders) == 0 {
		return []models.Tender{}, fmt.Errorf("%s: %w", operationPlace, outerror.ErrTendersWithThisServiceTypeNotFound)
	}

	return tenders, nil
}
func (storage *Storage) GetEmployeeTenders(ctx context.Context, empl models.Employee) (t []models.Tender, err error) {
	const operationPlace = "repository.postgres.tender.GetEmployeeTenders"
	query := `select name, description, service_type, status, organization_id, creator_username 
				from tender
				where creator_username = $1 and is_active_version=$2`
	tenders := []models.Tender{}

	rows, err := storage.connection.Query(ctx, query, empl.Username, true)
	if err != nil {
		return []models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
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
			return []models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
		}

		tenders = append(tenders, tender)

		if err := rows.Err(); err != nil {
			return []models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}
	if len(tenders) == 0 {
		return []models.Tender{}, fmt.Errorf("%s: %w", operationPlace, outerror.ErrEmployeeTendersNotFound)
	}
	return tenders, nil
}
func (storage *Storage) EditTender(
	ctx context.Context,
	oldTender models.Tender,
	tenderId int,
	updateTender models.TenderToUpdate,
) (models.Tender, error) {
	const operationPlace = "repository.postgres.tender.EditTender"

	insertQuery := `
	insert into tender values (@tender_id, @name, @desc, @srv_type, @status, @org_id, @username, @version, @is_active_version)
	returning name, description, service_type, status, organization_id, creator_username`

	lastTenderVersion, err := storage.getLastTenderVersion(ctx, tenderId)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s.getLastTenderVersion: %w", operationPlace, err)
	}
	args := pgx.NamedArgs{"is_active_version": true, "version": lastTenderVersion + 1, "tender_id": tenderId}

	if newName := updateTender.TenderName; newName == nil {
		args["name"] = oldTender.TenderName
	} else {
		args["name"] = *newName
	}

	if newDesc := updateTender.Description; newDesc == nil {
		args["desc"] = oldTender.Description
	} else {
		args["desc"] = *newDesc
	}

	if newSrvType := updateTender.ServiceType; newSrvType == nil {
		args["srv_type"] = oldTender.ServiceType
	} else {
		args["srv_type"] = *newSrvType
	}

	if newStatus := updateTender.Status; newStatus == nil {
		args["status"] = oldTender.Status
	} else {
		args["status"] = *newStatus
	}

	if newOrgId := updateTender.OrganizationId; newOrgId == nil {
		args["org_id"] = oldTender.OrganizationId
	} else {
		args["org_id"] = *newOrgId
	}

	if newUsername := updateTender.CreatorUsername; newUsername == nil {
		args["username"] = oldTender.CreatorUsername
	} else {
		args["username"] = *newUsername
	}

	tx, err := storage.connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	deactivateQuery := "update tender set is_active_version = $1 where tender_id = $2"
	_, err = tx.Exec(ctx, deactivateQuery, false, tenderId)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
	}

	var tender models.Tender
	row := tx.QueryRow(ctx, insertQuery, args)
	err = row.Scan(
		&tender.TenderName,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.OrganizationId,
		&tender.CreatorUsername,
	)
	if err != nil {
		return models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
	}

	return tender, nil
}
func (storage *Storage) RollbackTender(ctx context.Context, tenderId int, toVersionRollback int) (err error) {
	const operationPlace = "repository.postgres.tender.RollbackTender"
	deactivateVersionQuery := `update tender set is_active_version = $1 where tender_id = $2`
	rollbackQuery := `update tender set is_active_version = $1 where tender_id = $2 and version = $3`

	tx, err := storage.connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("%s: %w", operationPlace, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, deactivateVersionQuery, false, tenderId)
	if err != nil {
		return fmt.Errorf("%s: %w", operationPlace, err)
	}

	_, err = tx.Exec(ctx, rollbackQuery, true, tenderId, toVersionRollback)
	if err != nil {
		return fmt.Errorf("%s: %w", operationPlace, err)
	}
	return nil
}
func (storage *Storage) GetTenderById(ctx context.Context, tenderId int) (models.Tender, error) {
	const operationPlace = "repository.postgres.tender.GetTenderById"
	query := `select name, description, service_type, status, organization_id, creator_username 
				from tender
				where tender_id = $1 and is_active_version=$2`

	var tender models.Tender

	row := storage.connection.QueryRow(ctx, query, tenderId, true)
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
			return models.Tender{}, fmt.Errorf("%s: %w", operationPlace, outerror.ErrTenderNotFound)
		} else {
			return models.Tender{}, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}

	return tender, nil
}

func (storage *Storage) FindTenderVersion(ctx context.Context, tenderId int, version int) error {
	const operationPlace = "repository.postgres.tender.getLastInsertedTenderId"
	query := `select tender_id from tender where tender_id = $1 and version = $2`
	var id int
	row := storage.connection.QueryRow(ctx, query, tenderId, version)
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", operationPlace, outerror.ErrTenderVersionNotFound)
		}
		return fmt.Errorf("%s: %w", operationPlace, err)
	}

	return nil
}

func (storage *Storage) GetTenderStatus(ctx context.Context, tenderStatus string) (string, error) {
	panic("impl me")
}

func (storage *Storage) GetLastInsertedTenderId(ctx context.Context) (int, error) {
	const operationPlace = "repository.postgres.tender.getLastInsertedTenderId"
	query := `select tender_id from tender order by tender_id desc limit 1`
	row := storage.connection.QueryRow(ctx, query)
	var tenderId int
	err := row.Scan(&tenderId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("%s: %w", operationPlace, err)
	}
	return tenderId, nil
}

func (storage *Storage) getLastTenderVersion(ctx context.Context, tenderId int) (int, error) {
	const operationPlace = "repository.postgres.tender.getLastTenderVersion"
	query := "select version from tender where tender_id = $1 order by version desc limit 1"
	row := storage.connection.QueryRow(ctx, query, tenderId)
	var version int
	err := row.Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operationPlace, err)
	}
	return version, nil
}
