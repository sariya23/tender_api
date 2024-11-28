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
