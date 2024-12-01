package tests

import (
	"context"
	"testing"

	"github.com/sariya23/tender/internal/domain/models"
	"github.com/sariya23/tender/internal/lib/logger/slogdiscard"
	"github.com/sariya23/tender/internal/service"
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
	mockTenderRepo.On("GetAllTenders", ctx).Return([]models.Tender{}, service.ErrTendersNotFound)

	// Act
	tenders, err := tenderService.GetTenders(ctx, "all")

	// Assert
	require.ErrorIs(t, err, service.ErrTendersNotFound)
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
