package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/sariya23/tender/internal/domain/models"
	"github.com/sariya23/tender/internal/lib/logger/slogdiscard"
	"github.com/sariya23/tender/internal/service/tender"
	"github.com/sariya23/tender/internal/service/tender/mocks"
	"github.com/stretchr/testify/require"
)

// TestCreateTenders_Success проверяет, что если
// тендер создает юзер который есть, и
func TestCreateTenders_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	tenderToCreate := models.Tender{
		TenderName:      "Tender 1",
		Description:     "qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: "qwe",
	}
	exptectedTender := models.Tender{
		TenderName:      "Tender 1",
		Description:     "qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: "qwe",
	}
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("CreateTender", ctx, tenderToCreate).Return(exptectedTender, nil)
	mockEmployeeRepo.On("GetEmployeeByUsername", ctx, "qwe").Return(models.Employee{}, nil)
	mockOrgRepo.On("GetOrganizationById", ctx, 1).Return(models.Organization{}, nil)
	mockResponsibler.On("CheckResponsibility", ctx, 0, 1).Return(nil)

	// Act
	tender, err := tenderService.CreateTender(ctx, tenderToCreate)

	// Assert
	require.NoError(t, err)
	require.Equal(t, exptectedTender, tender)
}

// TestCreateTenders_FailEmployeeNotFound проверяет, что тендер не создатеся,
// если не получилось проверить существует сотрудник или нет.
func TestCreateTenders_FailEmployeeNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	tenderToCreate := models.Tender{
		TenderName:      "Tender 1",
		Description:     "qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: "qwe",
	}
	repoErr := errors.New("some error")
	exptectedTender := models.Tender{}
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("CreateTender", ctx, tenderToCreate).Return(exptectedTender, nil)
	mockEmployeeRepo.On("GetEmployeeByUsername", ctx, "qwe").Return(models.Employee{}, repoErr)

	// Act
	tender, err := tenderService.CreateTender(ctx, tenderToCreate)

	// Assert
	require.Error(t, err)
	require.Empty(t, tender)
}

// TestCreateTenders_FailOrgNotFound проверяет, что тендер не создается
// если не получилось проверить существует ли организация или нет.
func TestCreateTenders_FailOrgNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	tenderToCreate := models.Tender{
		TenderName:      "Tender 1",
		Description:     "qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: "qwe",
	}
	repoErr := errors.New("some error")
	exptectedTender := models.Tender{}
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("CreateTender", ctx, tenderToCreate).Return(exptectedTender, nil)
	mockEmployeeRepo.On("GetEmployeeByUsername", ctx, "qwe").Return(models.Employee{}, nil)
	mockOrgRepo.On("GetOrganizationById", ctx, 1).Return(models.Organization{}, repoErr)

	// Act
	tender, err := tenderService.CreateTender(ctx, tenderToCreate)

	// Assert
	require.Error(t, err)
	require.Empty(t, tender)
}

// TestCreateTenders_FailEmployeeNotResponsibleForOrganization проверяет, что
// тендер не создается, если пользователь не ответсвенный за переданную организацию.
func TestCreateTenders_FailEmployeeNotResponsibleForOrganization(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	tenderToCreate := models.Tender{
		TenderName:      "Tender 1",
		Description:     "qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: "qwe",
	}
	repoErr := errors.New("some error")
	exptectedTender := models.Tender{}
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("CreateTender", ctx, tenderToCreate).Return(exptectedTender, nil)
	mockEmployeeRepo.On("GetEmployeeByUsername", ctx, "qwe").Return(models.Employee{}, nil)
	mockOrgRepo.On("GetOrganizationById", ctx, 1).Return(models.Organization{}, nil)
	mockResponsibler.On("CheckResponsibility", ctx, 0, 1).Return(repoErr)

	// Act
	tender, err := tenderService.CreateTender(ctx, tenderToCreate)

	// Assert
	require.Error(t, err)
	require.Empty(t, tender)
}
