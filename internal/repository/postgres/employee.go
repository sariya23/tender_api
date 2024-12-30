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
	const operationPlace = "repository.postgres.employee.GetEmployeeByUsername"
	query := "select employee_id, username, first_name, last_name from employee where username = $1"

	var employee models.Employee

	row := storage.connection.QueryRow(ctx, query, username)
	err := row.Scan(&employee.ID, &employee.Username, &employee.FirstName, &employee.LastName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Employee{}, fmt.Errorf("%s: %w", operationPlace, outerror.ErrEmployeeNotFound)
		} else {
			return models.Employee{}, fmt.Errorf("%s: %w", operationPlace, err)
		}
	}

	return employee, nil
}

func (storage *Storage) CreateEmployee(ctx context.Context, employee models.Employee) error {
	const operationPlace = "repository.postgres.employee.CreateEmployee"
	inserEmployee := "insert into employee (username, first_name, last_name) values (@username, @first_name, @last_name)"

	_, err := storage.connection.Exec(
		ctx,
		inserEmployee,
		pgx.NamedArgs{
			"username":   employee.Username,
			"first_name": employee.FirstName,
			"last_name":  employee.LastName,
		},
	)

	if err != nil {
		return fmt.Errorf("%s: %w", operationPlace, err)
	}

	return nil
}

func (storage *Storage) GetEmployeeById(ctx context.Context, id int) (models.Employee, error) {
	panic("impl me")
}
