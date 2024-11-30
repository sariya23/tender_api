package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sariya23/tender/internal/domain/models"
)

type Storage struct {
	connection *pgxpool.Pool
}

func (s *Storage) MustNewConnection(dbURL string) *Storage {
	panic("impl me")
}

func (s *Storage) CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	panic("impl me")
}
func (s *Storage) GetAllTendets(ctx context.Context) ([]models.Tender, error) {
	panic("impl me")
}

func (s *Storage) GetTendersByServiceType(ctx context.Context, serviceType string) ([]models.Tender, error) {
	panic("impl me")
}
func (s *Storage) GetUserTenders(ctx context.Context, username string) ([]models.Tender, error) {
	panic("impl me")
}
func (s *Storage) EditTender(ctx context.Context, tenderId int, updateTender models.TenderToUpdate) (models.TenderToUpdate, error) {
	panic("impl me")
}
func (s *Storage) RollbackTender(ctx context.Context, tenderId int, toVersionRollback int) (models.Tender, error) {
	panic("impl me")
}
func (s *Storage) GetTenderById(ctx context.Context, tenderId int) (models.Tender, error) {
	panic("impl me")
}
func (s *Storage) GetEmployeeByUsername(ctx context.Context, username string) (models.Employee, error) {
	panic("impl me")
}
func (s *Storage) GetEmployeeById(ctx context.Context, id int) (models.Employee, error) {
	panic("impl me")
}

func (s Storage) GetOrganizationById(ctx context.Context, orgId int) (models.Organization, error) {
	panic("impl me")
}

func (s Storage) CheckResponsibility(ctx context.Context, emplId int, orgId int) error {
	panic("impl me")
}
