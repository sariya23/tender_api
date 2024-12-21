package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/sariya23/tender/internal/domain/models"
	"github.com/sariya23/tender/internal/lib/logger/slogdiscard"
	outerror "github.com/sariya23/tender/internal/out_error"
	"github.com/sariya23/tender/internal/service/tender"
	"github.com/sariya23/tender/internal/service/tender/mocks"
	"github.com/stretchr/testify/require"
)

// TestRollbackTender_Success проверяет, что
// тендер успешно откатывается на указанныую версию.
func TestRollbackTender_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	expectedTender := models.Tender{
		TenderName:      "Tender 1",
		Description:     "qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: "qwe",
	}

	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTenderById", ctx, 2).Return(models.Tender{CreatorUsername: "qwe"}, nil).Once()
	mockTenderRepo.On("GetTenderById", ctx, 2).Return(expectedTender, nil).Once()
	mockTenderRepo.On("FindTenderVersion", ctx, 2, 1).Return(nil)
	mockTenderRepo.On("RollbackTender", ctx, 2, 1).Return(nil)

	// Act
	tender, err := tenderService.RollbackTender(ctx, 2, 1, "qwe")

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedTender, tender)
}

// TestRollbackTender_FailTenderNotFound проверяет, что если
// тендера с таким id нет, то возвращается ошибка.
func TestRollbackTender_FailTenderNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTenderById", ctx, 2).Return(models.Tender{}, outerror.ErrTenderNotFound)

	// Act
	tender, err := tenderService.RollbackTender(ctx, 2, 1, "qwe")

	// Assert
	require.ErrorIs(t, err, outerror.ErrTenderNotFound)
	require.Equal(t, models.Tender{}, tender)
}

// TestRollbackTender_FailVersionNotFound проверяет, что
// если указанной версии нет, то возвращается ошибка.
func TestRollbackTender_FailVersionNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTenderById", ctx, 2).Return(models.Tender{CreatorUsername: "qwe"}, nil)
	mockTenderRepo.On("FindTenderVersion", ctx, 2, 1).Return(outerror.ErrTenderVersionNotFound)

	// Act
	tender, err := tenderService.RollbackTender(ctx, 2, 1, "qwe")

	// Assert
	require.ErrorIs(t, err, outerror.ErrTenderVersionNotFound)
	require.Equal(t, models.Tender{}, tender)
}

// TestRollbackTender_FailCannotRollbackTender проверяет, что
// в случае непредвиденной ошибки при откате тендера, возвращается эта ошибка.
func TestRollbackTender_FailCannotRollbackTender(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	someErr := errors.New("some err")
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTenderById", ctx, 2).Return(models.Tender{CreatorUsername: "qwe"}, nil)
	mockTenderRepo.On("FindTenderVersion", ctx, 2, 1).Return(nil)
	mockTenderRepo.On("RollbackTender", ctx, 2, 1).Return(someErr)

	// Act
	tender, err := tenderService.RollbackTender(ctx, 2, 1, "qwe")

	// Assert
	require.ErrorIs(t, err, someErr)
	require.Equal(t, models.Tender{}, tender)
}

// TestRollbackTender_FailEmployeeNotResponsibleForTender проверяет, что
// пользователь, который неответсвенный за тендер, который он хочет откатить, не
// сможет этого сделать.
func TestRollbackTender_FailEmployeeNotResponsibleForTender(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTenderById", ctx, 2).Return(models.Tender{CreatorUsername: "qwe"}, nil)
	mockTenderRepo.On("FindTenderVersion", ctx, 2, 1).Return(nil)

	// Act
	tender, err := tenderService.RollbackTender(ctx, 2, 1, "zxc")

	// Assert
	require.ErrorIs(t, err, outerror.ErrEmployeeNotResponsibleForTender)
	require.Equal(t, models.Tender{}, tender)
}
