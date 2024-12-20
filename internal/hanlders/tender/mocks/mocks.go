package mocks

import (
	"context"

	"github.com/sariya23/tender/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

// MockTenderServiceProvider реализует интерфейс TenderServiceProvider
// для целей тестирования. Он позволяет задавать ожидаемые результаты
// методов:
//
// - CreateTender
//
// - GetTenders
//
// - GetEmployeeTendersByUsername
//
// - EditTender
type MockTenderServiceProvider struct {
	mock.Mock
}

func (m *MockTenderServiceProvider) GetTenders(ctx context.Context, serviceType string) ([]models.Tender, error) {
	args := m.Called(ctx, serviceType)
	return args.Get(0).([]models.Tender), args.Error(1)
}

func (m *MockTenderServiceProvider) GetEmployeeTendersByUsername(ctx context.Context, username string) ([]models.Tender, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]models.Tender), args.Error(1)
}

func (m *MockTenderServiceProvider) CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	args := m.Called(ctx, tender)
	return args.Get(0).(models.Tender), args.Error(1)
}

func (m *MockTenderServiceProvider) EditTender(ctx context.Context, tenderId int, updateTender models.TenderToUpdate, username string) (models.Tender, error) {
	args := m.Called(ctx, tenderId, updateTender, username)
	return args.Get(0).(models.Tender), args.Error(1)
}

func (m *MockTenderServiceProvider) RollbackTender(ctx context.Context, tenderId int, version int, username string) (models.Tender, error) {
	args := m.Called(ctx, tenderId, version, username)
	return args.Get(0).(models.Tender), args.Error(1)
}
