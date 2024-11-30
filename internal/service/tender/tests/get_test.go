package tests

import (
	"context"
	"errors"
	"tender/internal/domain/models"
	"tender/internal/lib/logger/slogdiscard"
	"tender/internal/service/tender"
	"tender/internal/service/tender/mocks"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGetAllTenders_Success проверяет, что
// если передать serviceType = all, то вернутся
// все тендеры, которые есть.
func TestGetAllTenders_Success(t *testing.T) {
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
	mockTenderRepo.On("GetAll", ctx).Return(expectedTenders, nil)

	tenders, err := tenderService.GetTenders(ctx, "all")

	require.NoError(t, err)
	require.Equal(t, expectedTenders, tenders)
}

// TestGetAllTenders_FailGetAllTenders проверяет, что в случае
// какой-то ошибки возвращается пустой список тендеров и ошибка.
func TestGetAllTenders_FailGetAllTenders(t *testing.T) {
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	expectedTenders := []models.Tender{}
	repoErr := errors.New("some error")

	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetAll", ctx).Return(expectedTenders, repoErr)
	tenders, err := tenderService.GetTenders(ctx, "all")
	require.Error(t, err)
	require.Equal(t, expectedTenders, tenders)
}

// TestGetAllTenders_Success проверяет, что
// если передать serviceType = all, то вернутся
// все тендеры, которые есть.
func TestGetTendersByServiceType_Success(t *testing.T) {
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
	mockTenderRepo.On("GetByServiceType", ctx, "qwe").Return(expectedTenders, nil)
	tenders, err := tenderService.GetTenders(ctx, "qwe")
	require.NoError(t, err)
	require.Equal(t, expectedTenders, tenders)
}
