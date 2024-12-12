package mocks

import (
	"context"

	"github.com/sariya23/tender/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

// MockTenderRepo реализует интерфейс TenderRepository
// для целей тестирования. Он позволяет задавать ожидаемые результаты
// методов:
//
// - CreateTender
//
// - GetAllTenders
//
// - GetTendersByServiceType
//
// - GetEmployeeTendersByUsername
//
// - EditTender
//
// - RollbackTender
//
// - GetTenderById
type MockTenderRepo struct {
	mock.Mock
}

func (m *MockTenderRepo) GetLastInsertedTenderId(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockTenderRepo) CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	args := m.Called(ctx, tender)
	return args.Get(0).(models.Tender), args.Error(1)
}

func (m *MockTenderRepo) GetAllTenders(ctx context.Context) ([]models.Tender, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Tender), args.Error(1)
}

func (m *MockTenderRepo) GetTendersByServiceType(ctx context.Context, serviceType string) ([]models.Tender, error) {
	args := m.Called(ctx, serviceType)
	return args.Get(0).([]models.Tender), args.Error(1)
}

func (m *MockTenderRepo) GetEmployeeTenders(ctx context.Context, empl models.Employee) ([]models.Tender, error) {
	args := m.Called(ctx, empl)
	return args.Get(0).([]models.Tender), args.Error(1)
}

func (m *MockTenderRepo) EditTender(ctx context.Context, tenderId int, updateTender models.TenderToUpdate) (models.Tender, error) {
	args := m.Called(ctx, tenderId, updateTender)
	return args.Get(0).(models.Tender), args.Error(1)
}

func (m *MockTenderRepo) RollbackTender(ctx context.Context, tenderId int, toVersionRollback int) (models.Tender, error) {
	args := m.Called(ctx, tenderId, toVersionRollback)
	return args.Get(0).(models.Tender), args.Error(1)
}

func (m *MockTenderRepo) GetTenderById(ctx context.Context, tenderId int) (models.Tender, error) {
	args := m.Called(ctx, tenderId)
	return args.Get(0).(models.Tender), args.Error(1)
}

func (m *MockTenderRepo) FindTenderVersion(ctx context.Context, tenderId int, version int) error {
	args := m.Called(ctx, tenderId, version)
	return args.Error(0)
}

func (m *MockTenderRepo) GetTenderStatus(ctx context.Context, tenderStatus string) (string, error) {
	args := m.Called(ctx, tenderStatus)
	return args.Get(0).(string), args.Error(1)
}

// MockTenderRepo реализует интерфейс MockEmployeeRepo
// для целей тестирования. Он позволяет задавать ожидаемые результаты
// методов:
//
// - GetEmployeeByUsername
//
// - GetEmployeeById
type MockEmployeeRepo struct {
	mock.Mock
}

func (m *MockEmployeeRepo) GetEmployeeByUsername(ctx context.Context, username string) (models.Employee, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(models.Employee), args.Error(1)
}

func (m *MockEmployeeRepo) GetEmployeeById(ctx context.Context, id int) (models.Employee, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Employee), args.Error(1)
}

// MockOrgRepo реализует интерфейс OrganizationRepository
// для целей тестирования. Он позволяет задавать ожидаемые результаты
// методов:
//
// - GetOrganizationById
type MockOrgRepo struct {
	mock.Mock
}

func (m *MockOrgRepo) GetOrganizationById(ctx context.Context, orgId int) (models.Organization, error) {
	args := m.Called(ctx, orgId)
	return args.Get(0).(models.Organization), args.Error(1)
}

// MockEmployeeResponsibler реализует интерфейс EmployeeResponsibler
// для целей тестирования. Он позволяет задавать ожидаемые результаты
// методов:
//
// - CheckResponsibility
type MockEmployeeResponsibler struct {
	mock.Mock
}

func (m *MockEmployeeResponsibler) CheckResponsibility(ctx context.Context, emplId int, orgId int) error {
	args := m.Called(ctx, emplId, orgId)
	return args.Error(0)
}
