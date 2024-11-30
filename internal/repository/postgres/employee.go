package postgres

import (
	"context"

	"github.com/sariya23/tender/internal/domain/models"
)

func (s *Storage) GetEmployeeByUsername(ctx context.Context, username string) (models.Employee, error) {
	panic("impl me")
}
func (s *Storage) GetEmployeeById(ctx context.Context, id int) (models.Employee, error) {
	panic("impl me")
}
