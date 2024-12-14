package tender

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sariya23/tender/internal/domain/models"
	outerror "github.com/sariya23/tender/internal/out_error"
)

// EditTender обновляет тендер. Если пришел запрос на обновление:
//
// - Полей НЕ юзера или организации, то проверяется только наличие тендера для обновления
//
// - Поля юзера без поля организации (и другие поля), то проверяется существует ли этот юзер и отвественный ли он за текущую организацию
//
// - Поля организации без поля юзера (и другие поля), то проверяется существует ли эта организация и ответсвенный ли за него текущий юзер
//
// - И оля юзера, и поля организации (и другие поля), то проверяется существует ли этот юзер и организация и ответсвенный ли этот юзер за новую организацию.
func (s *TenderService) EditTender(ctx context.Context, tenderId int, updateTender models.TenderToUpdate) (models.Tender, error) {
	const op = "internal.service.tender.update.Edit"
	logger := s.logger.With("op", op)

	var updatedEmpl models.Employee
	var updatedOrg models.Organization
	var err error

	currTender, err := s.tenderRepo.GetTenderById(ctx, tenderId)
	if err != nil {
		if errors.Is(err, outerror.ErrTenderNotFound) {
			logger.Warn("tender not found", slog.Int("tender id", tenderId))
			return models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrTenderNotFound)
		}
		logger.Error(
			"cannot get tender by id",
			slog.Int("tender id", tenderId),
			slog.String("err", err.Error()),
		)
		return models.Tender{}, fmt.Errorf("cannot get tender by id: %w", err)
	}

	updatedUsername := updateTender.CreatorUsername
	updatedOrgId := updateTender.OrganizationId

	// Нужно проверить, что если сотрудник обновлися, то
	// он есть базе.
	if updatedUsername != nil {
		updatedEmpl, err = s.employeeRepo.GetEmployeeByUsername(ctx, *updatedUsername)
		if err != nil {
			if errors.Is(err, outerror.ErrEmployeeNotFound) {
				logger.Warn("employee to update not found", slog.String("username", *updatedUsername))
				return models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrEmployeeNotFound)
			}
			logger.Error(
				"cannot get employee by username",
				slog.String("username", *updatedUsername),
				slog.String("err", err.Error()),
			)
			return models.Tender{}, fmt.Errorf("cannot get employee by username: %w", err)
		}
	}

	// Нужно проверить, что если обновалась организация, то
	// она есть в базе.
	if updatedOrgId != nil {
		updatedOrg, err = s.orgRepo.GetOrganizationById(ctx, *updatedOrgId)
		if err != nil {
			if errors.Is(err, outerror.ErrOrganizationNotFound) {
				logger.Warn("organization to update not found", slog.Int("org id", *updatedOrgId))
				return models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrOrganizationNotFound)
			}
			logger.Error(
				"cannot get org by id",
				slog.Int("org id", *updatedOrgId),
				slog.String("err", err.Error()),
			)
			return models.Tender{}, fmt.Errorf("cannot get organization by id: %w", err)
		}
	}

	// Нужно проверить, что если обновился и юзер, и организация,
	// то сотрудник ответсвенный за новую огранизацию.
	if updatedUsername != nil && updatedOrgId != nil {
		err = s.employeeResponsibler.CheckResponsibility(ctx, updatedEmpl.ID, updatedOrg.ID)
		if err != nil {
			if errors.Is(err, outerror.ErrEmployeeNotResponsibleForOrganization) {
				logger.Warn("updated employee not responsible for updated organization", slog.Int("employee id", updatedEmpl.ID), slog.Int("org id", updatedOrg.ID))
				return models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrEmployeeNotResponsibleForOrganization)
			}
			logger.Error(
				"cannot check new employee responsibility to new org",
				slog.Int("employee id", updatedEmpl.ID),
				slog.String("err", err.Error()),
			)
			return models.Tender{}, fmt.Errorf("cannot check new employee responsibility to new org: %w", err)
		}
	}

	// Нужно проверить, что если поменялся только сотрудник, то
	// он ответсвенный за текущую организацию.
	if updatedUsername != nil && updatedOrgId == nil {
		err = s.employeeResponsibler.CheckResponsibility(ctx, updatedEmpl.ID, currTender.OrganizationId)
		if err != nil {
			if errors.Is(err, outerror.ErrEmployeeNotResponsibleForOrganization) {
				logger.Warn("updated employee not responsible for current organization", slog.Int("employee id", updatedEmpl.ID), slog.Int("org id", currTender.OrganizationId))
				return models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrEmployeeNotResponsibleForOrganization)
			}
			logger.Error(
				"cannot check new employee responsibility to current org",
				slog.Int("employee id", updatedEmpl.ID),
				slog.String("err", err.Error()),
			)
			return models.Tender{}, fmt.Errorf("cannot check new employee responsibility to current org: %w", err)
		}
	}

	var currEmpl models.Employee

	// Нужно проверить, что если поменялась только огранизация, то
	// текущий сотрудник ответсвенный за новую организацию.
	if updatedUsername == nil && updatedOrgId != nil {
		currEmpl, err = s.employeeRepo.GetEmployeeByUsername(ctx, currTender.CreatorUsername)
		if err != nil {
			logger.Error("cannot get employee by id", slog.String("err", err.Error()))
			return models.Tender{}, err
		}
		err = s.employeeResponsibler.CheckResponsibility(ctx, currEmpl.ID, *updatedOrgId)
		if err != nil {
			if errors.Is(err, outerror.ErrEmployeeNotResponsibleForOrganization) {
				logger.Warn("current employee not responsible for updated organization", slog.Int("employee id", updatedEmpl.ID), slog.Int("org id", currTender.OrganizationId))
				return models.Tender{}, fmt.Errorf("%s: %w", op, outerror.ErrEmployeeNotResponsibleForOrganization)
			}
			logger.Error(
				"cannot check current employee responsobility for updated org",
				slog.String("curr employee", currEmpl.Username),
				slog.Int("new org id", *updatedOrgId),
			)
			return models.Tender{}, err
		}
	}

	updatedTender, err := s.tenderRepo.EditTender(ctx, currTender, tenderId, updateTender)

	if err != nil {
		logger.Error("cannot update tender", slog.String("err", err.Error()))
		return models.Tender{}, err
	}

	return updatedTender, nil
}
