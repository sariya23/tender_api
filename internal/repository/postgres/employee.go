package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sariya23/tender/internal/domain/models"
	outerror "github.com/sariya23/tender/internal/out_error"
)

func (storage *Storage) GetEmployeeByUsername(ctx context.Context, username string) (models.Employee, error) {
	const op = "repository.postgres.employee.GetEmployeeByUsername"
	query := "select employee_id, username, first_name, last_name from employee where username = $1"

	var employee models.Employee

	row := storage.connection.QueryRow(ctx, query, username)
	err := row.Scan(&employee.ID, &employee.Username, &employee.FirstName, &employee.LastName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Employee{}, fmt.Errorf("%s: %w", op, outerror.ErrEmployeeNotFound)
		} else {
			return models.Employee{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	return employee, nil
}

func (storage *Storage) GetEmployeeById(ctx context.Context, id int) (models.Employee, error) {
	panic("impl me")
}
