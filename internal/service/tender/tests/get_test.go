package tests

import (
	"context"
	"testing"

	"github.com/sariya23/tender/internal/domain/models"
	"github.com/sariya23/tender/internal/lib/logger/slogdiscard"
	outerror "github.com/sariya23/tender/internal/out_error"
	"github.com/sariya23/tender/internal/service/tender"
	"github.com/sariya23/tender/internal/service/tender/mocks"
	"github.com/stretchr/testify/require"
)

// TestGetAllTenders_Success проверяет, что
// если передать serviceType = all, то вернутся
// все тендеры, которые есть.
func TestGetAllTenders_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	expectedTenders := []models.Tender{
		{TenderName: "Tender 1", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 1, CreatorUsername: "qwe"},
		{TenderName: "Tender 2", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 2, CreatorUsername: "zxc"},
	}
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetAllTenders", ctx).Return(expectedTenders, nil)

	// Act
	tenders, err := tenderService.GetTenders(ctx, "all")

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedTenders, tenders)
}

// TestGetAllTenders_FailGetAllTenders проверяет, что в случае
// какой-то ошибки возвращается пустой список тендеров и ошибка.
func TestGetAllTenders_FailGetAllTenders(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetAllTenders", ctx).Return([]models.Tender{}, outerror.ErrTendersWithThisServiceTypeNotFound)

	// Act
	tenders, err := tenderService.GetTenders(ctx, "all")

	// Assert
	require.ErrorIs(t, err, outerror.ErrTendersWithThisServiceTypeNotFound)
	require.Empty(t, tenders)
}

// TestGetAllTenders_Success проверяет, что
// если передать serviceType = all, то вернутся
// все тендеры, которые есть.
func TestGetTendersByServiceType_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	expectedTenders := []models.Tender{
		{TenderName: "Tender 1", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 1, CreatorUsername: "qwe"},
		{TenderName: "Tender 2", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 2, CreatorUsername: "zxc"},
	}
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTendersByServiceType", ctx, "qwe").Return(expectedTenders, nil)

	// Act
	tenders, err := tenderService.GetTenders(ctx, "qwe")

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedTenders, tenders)
}

// TestGetEmployeeTenders_Success проверяет, что
// если сотрудник существует в базе и у него есть связанные тендеры,
// то возвращается список этих тендеров.
func TestGetEmployeeTenders_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	usermame := "qwe"
	expectedTenders := []models.Tender{
		{TenderName: "Tender 1", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 1, CreatorUsername: usermame},
		{TenderName: "Tender 2", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 2, CreatorUsername: usermame},
	}
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockEmployeeRepo.On("GetEmployeeByUsername", ctx, usermame).Return(models.Employee{Username: usermame}, nil)
	mockTenderRepo.On("GetEmployeeTendersByUsername", ctx, usermame).Return(expectedTenders, nil)

	// Act
	tenders, err := tenderService.GetEmployeeTendersByUsername(ctx, usermame)

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedTenders, tenders)
}

// TestGetEmployeeTenders_FailEmployeeNotFound проверяет, что
// если сотрудника не существует, то возвращается ошибка и пустой список
// тендеров.
func TestGetEmployeeTenders_FailEmployeeNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	usermame := "qwe"
	expectedTenders := []models.Tender{}
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockEmployeeRepo.On("GetEmployeeByUsername", ctx, usermame).Return(models.Employee{}, outerror.ErrEmployeeNotFound)

	// Act
	tenders, err := tenderService.GetEmployeeTendersByUsername(ctx, usermame)

	// Assert
	require.ErrorIs(t, err, outerror.ErrEmployeeNotFound)
	require.Equal(t, expectedTenders, tenders)
}

func TestGetEmployeeTenders_FailEmployeeTendersNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	usermame := "qwe"
	expectedTenders := []models.Tender{}
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockEmployeeRepo.On("GetEmployeeByUsername", ctx, usermame).Return(models.Employee{Username: usermame}, nil)
	mockTenderRepo.On("GetEmployeeTendersByUsername", ctx, usermame).Return(expectedTenders, outerror.ErrEmployeeTendersNotFound)

	// Act
	tenders, err := tenderService.GetEmployeeTendersByUsername(ctx, usermame)

	// Assert
	require.ErrorIs(t, err, outerror.ErrEmployeeTendersNotFound)
	require.Equal(t, expectedTenders, tenders)
}
