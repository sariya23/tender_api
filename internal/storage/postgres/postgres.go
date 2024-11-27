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

func (s *Storage) Tenders(ctx context.Context) ([]models.Tedner, error) {
	panic("implement me")
}

func (s *Storage) TendersByServiceType(ctx context.Context, serviceType string) ([]models.Tedner, error) {
	panic("implement me")
}
