package mocks

import (
	"context"

	"github.com/sariya23/tender/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

type MockTenderServiceProvider struct {
	mock.Mock
}

func (m *MockTenderServiceProvider) GetTenders(ctx context.Context, serviceType string) ([]models.Tender, error) {
	args := m.Called(ctx, serviceType)
	return args.Get(0).([]models.Tender), args.Error(1)
}

func (m *MockTenderServiceProvider) GetUserTenders(ctx context.Context, username string) ([]models.Tender, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]models.Tender), args.Error(1)
}

func (m *MockTenderServiceProvider) CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	args := m.Called(ctx, tender)
	return args.Get(0).(models.Tender), args.Error(1)
}

func (m *MockTenderServiceProvider) Edit(ctx context.Context, tenderId int, updateTender models.TenderToUpdate) (models.TenderToUpdate, error) {
	args := m.Called(ctx, tenderId, updateTender)
	return args.Get(0).(models.TenderToUpdate), args.Error(1)
}
