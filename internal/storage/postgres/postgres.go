package postgres

import (
	"context"
	"tender/internal/domain/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	connection *pgxpool.Pool
}

func MustNewConnection(ctx context.Context, dbURL string) *Storage {
	panic("implement me")
}

func (s *Storage) Tenders(ctx context.Context) ([]models.Tender, error) {
	panic("implement me")
}

func (s *Storage) TendersByServiceType(ctx context.Context, serviceType string) ([]models.Tender, error) {
	panic("implement me")
}

func (s *Storage) CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	panic("impl me")
}

func (s *Storage) UserByUsername(ctx context.Context, username string) error {
	panic("impl me")
}

func (s *Storage) OrganizationById(ctx context.Context, ogranizationId int) error {
	panic("impl me")
}

func (s *Storage) CheckUserResponsible(ctx context.Context, username string, organizationId int) error {
	panic("impl me")
}

func (s *Storage) UserTenders(ctx context.Context, username string) ([]models.Tender, error) {
	panic("impl me")
}
