package postgres

import (
	"context"

	"github.com/sariya23/tender/internal/domain/models"
)

func (s *Storage) CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	panic("impl me")
}
func (s *Storage) GetAllTenders(ctx context.Context) ([]models.Tender, error) {
	panic("impl me")
}

func (s *Storage) GetTendersByServiceType(ctx context.Context, serviceType string) ([]models.Tender, error) {
	panic("impl me")
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
