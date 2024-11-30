package tests

import (
	"context"
	"tender/internal/domain/models"
	"tender/internal/lib/logger/slogdiscard"
	"tender/internal/service/tender"
	"tender/internal/service/tender/mocks"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCreateTenders_Success проверяет, что если
// тендер создает юзер который есть, и
func TestCreateTenders_Success(t *testing.T) {
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
	mockTenderRepo.On("Create", ctx, tenderToCreate).Return(exptectedTender, nil)
	mockEmployeeRepo.On("GetByUsername", ctx, "qwe").Return(models.Employee{}, nil)
	mockOrgRepo.On("GetById", ctx, 1).Return(models.Organization{}, nil)

	tender, err := tenderService.CreateTender(ctx, tenderToCreate)

	require.NoError(t, err)
	require.Equal(t, exptectedTender, tender)
}
