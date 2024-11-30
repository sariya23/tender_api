package mocks

import (
	"context"
	"tender/internal/domain/models"

	"github.com/stretchr/testify/mock"
)

type MockTenderRepo struct {
	mock.Mock
}

func (m *MockTenderRepo) Create(ctx context.Context, tender models.Tender) (models.Tender, error) {
	args := m.Called(ctx, tender)
	return args.Get(0).(models.Tender), args.Error(1)
}

func (m *MockTenderRepo) GetAll(ctx context.Context) ([]models.Tender, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Tender), args.Error(1)
}

func (m *MockTenderRepo) GetByServiceType(ctx context.Context, serviceType string) ([]models.Tender, error) {
	args := m.Called(ctx, serviceType)
	return args.Get(0).([]models.Tender), args.Error(1)
}

func (m *MockTenderRepo) GetUserTenders(ctx context.Context, username string) ([]models.Tender, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]models.Tender), args.Error(1)
}

func (m *MockTenderRepo) Edit(ctx context.Context, tenderId int, updateTender models.TenderToUpdate) (models.TenderToUpdate, error) {
	args := m.Called(ctx, tenderId, updateTender)
	return args.Get(0).(models.TenderToUpdate), args.Error(1)
}

func (m *MockTenderRepo) Rollback(ctx context.Context, tenderId int, toVersionRollback int) (models.Tender, error) {
	args := m.Called(ctx, tenderId, toVersionRollback)
	return args.Get(0).(models.Tender), args.Error(1)
}

func (m *MockTenderRepo) GetById(ctx context.Context, tenderId int) (models.Tender, error) {
	args := m.Called(ctx, tenderId)
	return args.Get(0).(models.Tender), args.Error(1)
}

type MockEmployeeRepo struct {
	mock.Mock
}

func (m *MockEmployeeRepo) GetByUsername(ctx context.Context, username string) (models.Employee, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(models.Employee), args.Error(1)
}

func (m *MockEmployeeRepo) GetById(ctx context.Context, id int) (models.Employee, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Employee), args.Error(1)
}

type MockOrgRepo struct {
	mock.Mock
}

func (m *MockOrgRepo) GetById(ctx context.Context, orgId int) (models.Organization, error) {
	args := m.Called(ctx, orgId)
	return args.Get(0).(models.Organization), args.Error(1)
}

type MockEmployeeResponsibler struct {
	mock.Mock
}

func (m *MockEmployeeResponsibler) CheckResponsibility(ctx context.Context, emplId int, orgId int) error {
	args := m.Called(ctx, emplId, orgId)
	return args.Error(0)
}
