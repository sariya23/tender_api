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

// TestUpdateTender_SuccessChangeDesc проверяет, что
// если обновить описание или другие поля, кроме юзера и организации, то
// никаких проверок на юзера и организацию не будет и тендер обновится.
func TestUpdateTender_SuccessChangeDesc(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	name := "qwe"
	desc := "qwe"
	srvType := "qwer"
	status := "qwe"
	orgId := 1
	user := "qwe"
	currTender := models.Tender{
		TenderName:      name,
		Description:     "zxc",
		ServiceType:     srvType,
		Status:          status,
		OrganizationId:  orgId,
		CreatorUsername: user,
	}
	updateTender := models.TenderToUpdate{Description: &desc}
	exptectedTender := models.Tender{
		TenderName:      name,
		Description:     desc,
		ServiceType:     srvType,
		Status:          status,
		OrganizationId:  orgId,
		CreatorUsername: user,
	}

	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTenderById", ctx, 1).Return(currTender, nil)
	mockTenderRepo.On("EditTender", ctx, 1, updateTender).Return(exptectedTender, nil)

	// Act
	tender, err := tenderService.EditTender(ctx, 1, updateTender)

	// Assert
	require.NoError(t, err)
	require.Equal(t, exptectedTender, tender)
}

// TestUpdateTender_SuccessChangeOnlyUser проверяет, что если при
// обновлении тендера меняется юзер, то должны быть вызваны проверки,
// что такой сотрудник есть, и что он отвественный за организацию, которая
// указана в тендере.
func TestUpdateTender_SuccessChangeOnlyUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	name := "qwe"
	desc := "qwe"
	srvType := "qwer"
	status := "qwe"
	orgId := 1
	user := "qwe"
	currTender := models.Tender{
		TenderName:      name,
		Description:     desc,
		ServiceType:     srvType,
		Status:          status,
		OrganizationId:  orgId,
		CreatorUsername: "zxc",
	}
	updateTender := models.TenderToUpdate{CreatorUsername: &user}
	exptectedTender := models.Tender{
		TenderName:      name,
		Description:     desc,
		ServiceType:     srvType,
		Status:          status,
		OrganizationId:  orgId,
		CreatorUsername: user,
	}

	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTenderById", ctx, 1).Return(currTender, nil)
	mockEmployeeRepo.On("GetEmployeeByUsername", ctx, user).Return(models.Employee{ID: 2, Username: user}, nil)
	mockResponsibler.On("CheckResponsibility", ctx, 2, 1).Return(nil)
	mockTenderRepo.On("EditTender", ctx, 1, updateTender).Return(exptectedTender, nil)

	// Act
	tender, err := tenderService.EditTender(ctx, 1, updateTender)

	// Assert
	require.NoError(t, err)
	require.Equal(t, exptectedTender, tender)
}

// TestUpdateTender_SuccessChangeOnlyOrg проверяет, что если при
// обновлении тендера поменялась только организация, но не юзер, то
// будут выполнены проверки, что такая организация есть в базе, и что
// новый юзер ответсвенный за эту новую организацию.
func TestUpdateTender_SuccessChangeOnlyOrg(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	name := "qwe"
	desc := "qwe"
	srvType := "qwer"
	status := "qwe"
	orgId := 2
	user := "qwe"
	currTender := models.Tender{
		TenderName:      name,
		Description:     desc,
		ServiceType:     srvType,
		Status:          status,
		OrganizationId:  2,
		CreatorUsername: user,
	}
	updateTender := models.TenderToUpdate{CreatorUsername: &user}
	exptectedTender := models.Tender{
		TenderName:      name,
		Description:     desc,
		ServiceType:     srvType,
		Status:          status,
		OrganizationId:  orgId,
		CreatorUsername: user,
	}

	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTenderById", ctx, 1).Return(currTender, nil)
	mockOrgRepo.On("GetOrganizationById", ctx, orgId).Return(models.Organization{ID: 2}, nil)
	mockEmployeeRepo.On("GetEmployeeByUsername", ctx, user).Return(models.Employee{ID: 2, Username: "qwe"}, nil)
	mockResponsibler.On("CheckResponsibility", ctx, 2, 2).Return(nil)
	mockTenderRepo.On("EditTender", ctx, 1, updateTender).Return(exptectedTender, nil)

	// Act
	tender, err := tenderService.EditTender(ctx, 1, updateTender)

	// Assert
	require.NoError(t, err)
	require.Equal(t, exptectedTender, tender)
}

// TestUpdateTender_SuccessChangeOrgAndUser проверяет, что
// если поменялся и юзер, и организаця, то будут вызваны проверки
// по существованию юзера, организации и, что юзер ответсвенный
// за организацию.
func TestUpdateTender_SuccessChangeOrgAndUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	name := "qwe"
	desc := "qwe"
	srvType := "qwer"
	status := "qwe"
	orgId := 1
	user := "qwe"
	currTender := models.Tender{
		TenderName:      name,
		Description:     desc,
		ServiceType:     srvType,
		Status:          status,
		OrganizationId:  2,
		CreatorUsername: "zxc",
	}
	updateTender := models.TenderToUpdate{CreatorUsername: &user, OrganizationId: &orgId}
	exptectedTender := models.Tender{
		TenderName:      name,
		Description:     desc,
		ServiceType:     srvType,
		Status:          status,
		OrganizationId:  orgId,
		CreatorUsername: user,
	}

	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTenderById", ctx, 1).Return(currTender, nil)
	mockOrgRepo.On("GetOrganizationById", ctx, orgId).Return(models.Organization{ID: 1}, nil)
	mockEmployeeRepo.On("GetEmployeeByUsername", ctx, user).Return(models.Employee{ID: 2, Username: user}, nil)
	mockResponsibler.On("CheckResponsibility", ctx, 2, 1).Return(nil)
	mockTenderRepo.On("EditTender", ctx, 1, updateTender).Return(exptectedTender, nil)

	// Act
	tender, err := tenderService.EditTender(ctx, 1, updateTender)

	// Assert
	require.NoError(t, err)
	require.Equal(t, exptectedTender, tender)
}

// TestUpdateTender_FailTenderNotFound проверяет, что
// если тендер, который хотят обновить, не существует, то
// возвращается пустой тендер и ошибка.
func TestUpdateTender_FailTenderNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	desc := "qwe"
	tenderToUpdate := models.TenderToUpdate{Description: &desc}
	logger := slogdiscard.NewDiscardLogger()
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTenderById", ctx, 1).Return(models.Tender{}, outerror.ErrTenderNotFound)

	// Act
	tender, err := tenderService.EditTender(ctx, 1, tenderToUpdate)

	// Assert
	require.ErrorIs(t, err, outerror.ErrTenderNotFound)
	require.Equal(t, tender, models.Tender{})
}

// TestUpdateTender_FailNewUserNotFound проверяет, что
// если новый юзер не существует, то вернется пустой тендер и ошибка.
func TestUpdateTender_FailNewUserNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	user := "qwe"
	tenderToUpdate := models.TenderToUpdate{CreatorUsername: &user}
	logger := slogdiscard.NewDiscardLogger()
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTenderById", ctx, 1).Return(models.Tender{}, nil)
	mockEmployeeRepo.On("GetEmployeeByUsername", ctx, user).Return(models.Employee{}, outerror.ErrEmployeeNotFound)

	// Act
	tender, err := tenderService.EditTender(ctx, 1, tenderToUpdate)

	// Assert
	require.ErrorIs(t, err, outerror.ErrEmployeeNotFound)
	require.Equal(t, tender, models.Tender{})
}

// TestUpdateTender_FailNewOrgNotFound проверяет, что если
// новой организации не существует, то возвращается пустой тендер и ошибка.
func TestUpdateTender_FailNewOrgNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	orgId := 1
	tenderToUpdate := models.TenderToUpdate{OrganizationId: &orgId}
	logger := slogdiscard.NewDiscardLogger()
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTenderById", ctx, 1).Return(models.Tender{}, nil)
	mockOrgRepo.On("GetOrganizationById", ctx, orgId).Return(models.Organization{}, outerror.ErrOrganizationNotFound)

	// Act
	tender, err := tenderService.EditTender(ctx, 1, tenderToUpdate)

	// Assert
	require.ErrorIs(t, err, outerror.ErrOrganizationNotFound)
	require.Equal(t, tender, models.Tender{})
}

// TestUpdateTender_FailNewNewUserNotResponsibleForNewOrg проверяет,
// что если новый юзер не ответсвенный за новую организацию, то возвращается
// пустой тендер и ошибка.
func TestUpdateTender_FailNewNewUserNotResponsibleForNewOrg(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	newOrgId := 1
	newUser := "qwe"
	tenderToUpdate := models.TenderToUpdate{OrganizationId: &newOrgId, CreatorUsername: &newUser}
	logger := slogdiscard.NewDiscardLogger()
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTenderById", ctx, 1).Return(models.Tender{}, nil)
	mockEmployeeRepo.On("GetEmployeeByUsername", ctx, newUser).Return(models.Employee{}, nil)
	mockOrgRepo.On("GetOrganizationById", ctx, newOrgId).Return(models.Organization{ID: 1}, nil)
	mockResponsibler.On("CheckResponsibility", ctx, 0, 1).Return(outerror.ErrEmployeeNotResponsibleForOrganization)

	// Act
	tender, err := tenderService.EditTender(ctx, 1, tenderToUpdate)

	// Assert
	require.ErrorIs(t, err, outerror.ErrEmployeeNotResponsibleForOrganization)
	require.Equal(t, tender, models.Tender{})
}

// TestUpdateTender_FailNewNewUserNotResponsibleForCurrOrg проверяет,
// что если новый юзер не ответсвенный за текущую организацию, то возвращается пустой
// тендер и ошибка.
func TestUpdateTender_FailNewNewUserNotResponsibleForCurrOrg(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	newUser := "qwe"
	tenderToUpdate := models.TenderToUpdate{CreatorUsername: &newUser}
	logger := slogdiscard.NewDiscardLogger()
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTenderById", ctx, 1).Return(models.Tender{}, nil)
	mockEmployeeRepo.On("GetEmployeeByUsername", ctx, newUser).Return(models.Employee{}, nil)
	mockResponsibler.On("CheckResponsibility", ctx, 0, 0).Return(outerror.ErrEmployeeNotResponsibleForOrganization)

	// Act
	tender, err := tenderService.EditTender(ctx, 1, tenderToUpdate)

	// Assert
	require.ErrorIs(t, err, outerror.ErrEmployeeNotResponsibleForOrganization)
	require.Equal(t, tender, models.Tender{})
}

// TestUpdateTender_FailCurrUserNotResponsibleForNewOrg проверяет,
// что если текущий сотрудник не ответсвенный за новую организацию,
// то возвращается пустой тендер и ошибка.
func TestUpdateTender_FailCurrUserNotResponsibleForNewOrg(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockTenderRepo := new(mocks.MockTenderRepo)
	mockEmployeeRepo := new(mocks.MockEmployeeRepo)
	mockOrgRepo := new(mocks.MockOrgRepo)
	mockResponsibler := new(mocks.MockEmployeeResponsibler)
	newOrg := 1
	tenderToUpdate := models.TenderToUpdate{OrganizationId: &newOrg}
	logger := slogdiscard.NewDiscardLogger()
	tenderService := tender.New(logger, mockTenderRepo, mockEmployeeRepo, mockOrgRepo, mockResponsibler)
	mockTenderRepo.On("GetTenderById", ctx, 1).Return(models.Tender{}, nil)
	mockOrgRepo.On("GetOrganizationById", ctx, 1).Return(models.Organization{}, nil)
	mockEmployeeRepo.On("GetEmployeeByUsername", ctx, "").Return(models.Employee{}, nil)
	mockResponsibler.On("CheckResponsibility", ctx, 0, 1).Return(outerror.ErrEmployeeNotResponsibleForOrganization)

	// Act
	tender, err := tenderService.EditTender(ctx, 1, tenderToUpdate)

	// Assert
	require.ErrorIs(t, err, outerror.ErrEmployeeNotResponsibleForOrganization)
	require.Equal(t, tender, models.Tender{})
}
