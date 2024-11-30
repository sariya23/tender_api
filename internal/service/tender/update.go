package tender

import (
	"context"
	"log/slog"

	"github.com/sariya23/tender/internal/domain/models"
)

func (s *TenderService) Edit(ctx context.Context, tenderId int, updateTender models.TenderToUpdate) (models.TenderToUpdate, error) {
	const op = "internal.service.tender.update.Edit"
	logger := s.logger.With("op", op)

	var updatedEmpl models.Employee
	var updatedOrg models.Organization
	var err error

	currTender, err := s.tenderRepo.GetTenderById(ctx, tenderId)
	if err != nil {
		logger.Error(
			"cannot get tender by id",
			slog.Int("tender id", tenderId),
			slog.String("err", err.Error()),
		)
		return models.TenderToUpdate{}, err
	}

	updatedUsername := updateTender.CreatorUsername
	updatedOrgId := updateTender.OrganizationId

	// Нужно проверить, что если пользователь обновлися, то
	// он есть базе.
	if updatedUsername != nil {
		updatedEmpl, err = s.employeeRepo.GetEmployeeByUsername(ctx, *updatedUsername)
		if err != nil {
			logger.Error(
				"cannot get employee by username",
				slog.String("username", *updatedUsername),
				slog.String("err", err.Error()),
			)
			return models.TenderToUpdate{}, err
		}
	}

	// Нужно проверить, что если обновалась организация, то
	// она есть в базе.
	if updatedOrgId != nil {
		updatedOrg, err = s.orgRepo.GetOrganizationById(ctx, *updatedOrgId)
		if err != nil {
			logger.Error(
				"cannot get org by id",
				slog.Int("org id", *updatedOrgId),
				slog.String("err", err.Error()),
			)
			return models.TenderToUpdate{}, err
		}
	}

	// Нужно проверить, что если обновился и юзер, и организация,
	// то пользователь ответсвенный за новую огранизацию.
	if updatedUsername != nil && updatedOrgId != nil {
		err = s.employeeResponsibler.CheckResponsibility(ctx, updatedEmpl.ID, updatedOrg.ID)
		if err != nil {
			logger.Error(
				"new employee not responsible to new org",
				slog.Int("employee id", updatedEmpl.ID),
				slog.String("err", err.Error()),
			)
			return models.TenderToUpdate{}, err
		}
	}

	// Нужно проверить, что если поменялся только пользователь, то
	// он ответсвенный за текущую организацию.
	if updatedUsername != nil && updatedOrgId == nil {
		err = s.employeeResponsibler.CheckResponsibility(ctx, updatedEmpl.ID, currTender.OrganizationId)
		if err != nil {
			logger.Error(
				"new employee not responsible for current org",
				slog.Int("org id", currTender.OrganizationId),
				slog.String("username", *updatedUsername),
			)
			return models.TenderToUpdate{}, err
		}
	}

	var currEmpl models.Employee

	// Нужно проверить, что если поменялась только огранизация, то
	// текущий пользователь ответсвенный за новую организацию.
	if updatedUsername == nil && updatedOrgId != nil {
		currEmpl, err = s.employeeRepo.GetEmployeeByUsername(ctx, currTender.CreatorUsername)
		if err != nil {
			logger.Error("cannot get employee by id", slog.String("err", err.Error()))
			return models.TenderToUpdate{}, err
		}
		err = s.employeeResponsibler.CheckResponsibility(ctx, currEmpl.ID, *updatedOrgId)
		if err != nil {
			logger.Error(
				"curr employee not responsible",
				slog.String("curr employee", currEmpl.Username),
				slog.Int("new org id", *updatedOrgId),
			)
			return models.TenderToUpdate{}, err
		}
	}

	updatedTender, err := s.tenderRepo.EditTender(ctx, tenderId, updateTender)

	if err != nil {
		logger.Error("cannot update tender", slog.String("err", err.Error()))
		return models.TenderToUpdate{}, err
	}

	return updatedTender, nil
}
