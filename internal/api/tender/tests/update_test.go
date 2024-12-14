package tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	tenderapi "github.com/sariya23/tender/internal/api/tender"
	"github.com/sariya23/tender/internal/api/tender/mocks"
	"github.com/sariya23/tender/internal/domain/models"
	"github.com/sariya23/tender/internal/lib/logger/slogdiscard"
	outerror "github.com/sariya23/tender/internal/out_error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEditTender_SuccessAllFiedsUpdate(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	tenderName := "update Tender 1"
	description := "update qwe"
	serviceType := "update op"
	status := "update open"
	organizationId := 2
	creatorUsername := "update qwe"

	mockTender := models.Tender{
		TenderName:      "update Tender 1",
		Description:     "update qwe",
		ServiceType:     "update op",
		Status:          "update open",
		OrganizationId:  2,
		CreatorUsername: "update qwe",
	}
	tenderToUpdate := models.TenderToUpdate{
		TenderName:      &tenderName,
		Description:     &description,
		ServiceType:     &serviceType,
		Status:          &status,
		OrganizationId:  &organizationId,
		CreatorUsername: &creatorUsername,
	}
	reqBody := `
		{
			"update_tender_data": {
				"name": "update Tender 1",
				"description": "update qwe",
				"service_type": "update op",
				"status": "update open",
				"organization_id": 2,
				"creator_username": "update qwe"
			}
		}`

	expectedBody := `
		{
			"updated_tender": {
				"name": "update Tender 1",
				"description": "update qwe",
				"service_type": "update op",
				"status": "update open",
				"organization_id": 2,
				"creator_username": "update qwe"
			},
			"message": "ok"
		}`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("EditTender", ctx, 2, tenderToUpdate).Return(mockTender, nil)
	router := gin.New()
	router.PATCH("/api/tenders/:tenderId/edit", svc.EditTender(ctx))
	req := httptest.NewRequest(http.MethodPatch, "/api/tenders/2/edit", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestEditTender_SuccessPartFiedsUpdate проверяет, что
// если передать часть полей для обновления, то обновление пройдет
// успешно.
func TestEditTender_SuccessPartFiedsUpdate(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	description := "update qwe"
	creatorUsername := "update qwe"

	mockTender := models.Tender{
		TenderName:      "Tender 1",
		Description:     "update qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: "update qwe",
	}
	tenderToUpdate := models.TenderToUpdate{
		Description:     &description,
		CreatorUsername: &creatorUsername,
	}

	reqBody := `
		{
			"update_tender_data": {
				"description": "update qwe",
				"creator_username": "update qwe"
			}
		}`

	expectedBody := `
		{
			"updated_tender": {
				"name": "Tender 1",
				"description": "update qwe",
				"service_type": "op",
				"status": "open",
				"organization_id": 1,
				"creator_username": "update qwe"
			},
			"message": "ok"
		}`
	svc := tenderapi.New(logger, mockTenderService)
	mockTenderService.On("EditTender", ctx, 2, tenderToUpdate).Return(mockTender, nil)
	router := gin.New()
	router.PATCH("/api/tenders/:tenderId/edit", svc.EditTender(ctx))
	req := httptest.NewRequest(http.MethodPatch, "/api/tenders/2/edit", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestEditTender_SuccessUndefinedFields проверяет, что
// если в json положить какие-то другие поля, то ошибки при валидации
// не будет.
func TestEditTender_SuccessUndefinedFields(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	description := "update qwe"
	creatorUsername := "update qwe"

	mockTender := models.Tender{
		TenderName:      "Tender 1",
		Description:     "update qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: "update qwe",
	}
	tenderToUpdate := models.TenderToUpdate{
		Description:     &description,
		CreatorUsername: &creatorUsername,
	}

	reqBody := `
		{
			"update_tender_data": {
				"description": "update qwe",
				"creator_username": "update qwe",
				"aboba": "qwe"
			}
		}`

	expectedBody := `
		{
			"updated_tender": {
				"name": "Tender 1",
				"description": "update qwe",
				"service_type": "op",
				"status": "open",
				"organization_id": 1,
				"creator_username": "update qwe"
			},
			"message": "ok"
		}`
	svc := tenderapi.New(logger, mockTenderService)
	mockTenderService.On("EditTender", ctx, 2, tenderToUpdate).Return(mockTender, nil)
	router := gin.New()
	router.PATCH("/api/tenders/:tenderId/edit", svc.EditTender(ctx))
	req := httptest.NewRequest(http.MethodPatch, "/api/tenders/2/edit", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestEditTender_FailInvalidTenderId проверяет, что
// если в пути :tenderId не является числом, то возвращается ошибка
// и код 400.
func TestEditTender_FailInvalidTenderId(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	reqBody := `
		{
			"update_tender_data": {
				"name": "update Tender 1",
				"description": "update qwe",
				"service_type": "update op",
				"status": "update open",
				"organization_id": 2,
				"creator_username": "update qwe"
			}
		}`

	expectedBody := `
		{
			"updated_tender": {
				"name": "",
				"description": "",
				"service_type": "",
				"status": "",
				"organization_id": 0,
				"creator_username": ""
			},
			"message": "wrong path"
		}`
	svc := tenderapi.New(logger, mockTenderService)

	router := gin.New()
	router.PATCH("/api/tenders/:tenderId/edit", svc.EditTender(ctx))
	req := httptest.NewRequest(http.MethodPatch, "/api/tenders/qwe/edit", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestEditTender_FailJsonSyntaxError проверяет, что
// если запрос синтаксически неверный, то возвращается ошибка и код 400.
func TestEditTender_FailJsonSyntaxError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	reqBody := `
		{
			"update_tender_data": {
				"name": "update Tender 1",
				"description": "update qwe"
				"service_type": "update op",
				"status": "update open",
				"organization_id": 2,
				"creator_username": "update qwe"
			}
		}`

	svc := tenderapi.New(logger, mockTenderService)

	router := gin.New()
	router.PATCH("/api/tenders/:tenderId/edit", svc.EditTender(ctx))
	req := httptest.NewRequest(http.MethodPatch, "/api/tenders/2/edit", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "json syntax err")
}

// TestEditTender_FailJsonTypeError проверяет, что
// если в запросе неверные типы, то возвращается ошибка и код 400.
func TestEditTender_FailJsonTypeError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	reqBody := `
		{
			"update_tender_data": {
				"name": "update Tender 1",
				"description": "update qwe",
				"service_type": false,
				"status": 123,
				"organization_id": 2,
				"creator_username": "update qwe"
			}
		}`

	svc := tenderapi.New(logger, mockTenderService)

	router := gin.New()
	router.PATCH("/api/tenders/:tenderId/edit", svc.EditTender(ctx))
	req := httptest.NewRequest(http.MethodPatch, "/api/tenders/2/edit", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "json type err")
}

// TestEditTender_FailValidationError проверяет, что
// если organization_id будет отрицательным, то вернется ошибка и
// код 400.
func TestEditTender_FailValidationError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	reqBody := `
		{
			"update_tender_data": {
				"name": "update Tender 1",
				"description": "update qwe",
				"service_type": "qwe",
				"status": "qwe",
				"organization_id": -2,
				"creator_username": "update qwe"
			}
		}`

	svc := tenderapi.New(logger, mockTenderService)

	router := gin.New()
	router.PATCH("/api/tenders/:tenderId/edit", svc.EditTender(ctx))
	req := httptest.NewRequest(http.MethodPatch, "/api/tenders/2/edit", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "validation failed")
}

// TestEditTender_FailTenderNotFound проверяет, что
// если тендера с таким id не существует, то возвращается ошибка
// и код 400.
func TestEditTender_FailTenderNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	tenderName := "update Tender 1"
	description := "update qwe"
	serviceType := "update op"
	status := "update open"
	organizationId := 2
	creatorUsername := "update qwe"

	tenderToUpdate := models.TenderToUpdate{
		TenderName:      &tenderName,
		Description:     &description,
		ServiceType:     &serviceType,
		Status:          &status,
		OrganizationId:  &organizationId,
		CreatorUsername: &creatorUsername,
	}
	reqBody := `
		{
			"update_tender_data": {
				"name": "update Tender 1",
				"description": "update qwe",
				"service_type": "update op",
				"status": "update open",
				"organization_id": 2,
				"creator_username": "update qwe"
			}
		}`

	expectedBody := `
		{
			"updated_tender": {
				"name": "",
				"description": "",
				"service_type": "",
				"status": "",
				"organization_id": 0,
				"creator_username": ""
			},
			"message": "tender with id=\"2\" not found"
		}`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("EditTender", ctx, 2, tenderToUpdate).Return(models.Tender{}, outerror.ErrTenderNotFound)
	router := gin.New()
	router.PATCH("/api/tenders/:tenderId/edit", svc.EditTender(ctx))
	req := httptest.NewRequest(http.MethodPatch, "/api/tenders/2/edit", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestEditTender_FailEmployeeNotFound проверяет, что
// если сотрудника с таким username не существует, то возвращается
// ошибка и код 400.
func TestEditTender_FailEmployeeNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	tenderName := "update Tender 1"
	description := "update qwe"
	serviceType := "update op"
	status := "update open"
	organizationId := 2
	creatorUsername := "update qwe"

	tenderToUpdate := models.TenderToUpdate{
		TenderName:      &tenderName,
		Description:     &description,
		ServiceType:     &serviceType,
		Status:          &status,
		OrganizationId:  &organizationId,
		CreatorUsername: &creatorUsername,
	}
	reqBody := `
		{
			"update_tender_data": {
				"name": "update Tender 1",
				"description": "update qwe",
				"service_type": "update op",
				"status": "update open",
				"organization_id": 2,
				"creator_username": "update qwe"
			}
		}`

	expectedBody := `
		{
			"updated_tender": {
				"name": "",
				"description": "",
				"service_type": "",
				"status": "",
				"organization_id": 0,
				"creator_username": ""
			},
			"message": "updated employee with username=\"update qwe\" not found"
		}`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("EditTender", ctx, 2, tenderToUpdate).Return(models.Tender{}, outerror.ErrEmployeeNotFound)
	router := gin.New()
	router.PATCH("/api/tenders/:tenderId/edit", svc.EditTender(ctx))
	req := httptest.NewRequest(http.MethodPatch, "/api/tenders/2/edit", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestEditTender_FailOrgNotFound проверяет, что
// если организации с таким id не существует, то возвращается
// ошибка и код 400.
func TestEditTender_FailOrgNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	tenderName := "update Tender 1"
	description := "update qwe"
	serviceType := "update op"
	status := "update open"
	organizationId := 2
	creatorUsername := "update qwe"

	tenderToUpdate := models.TenderToUpdate{
		TenderName:      &tenderName,
		Description:     &description,
		ServiceType:     &serviceType,
		Status:          &status,
		OrganizationId:  &organizationId,
		CreatorUsername: &creatorUsername,
	}
	reqBody := `
		{
			"update_tender_data": {
				"name": "update Tender 1",
				"description": "update qwe",
				"service_type": "update op",
				"status": "update open",
				"organization_id": 2,
				"creator_username": "update qwe"
			}
		}`

	expectedBody := `
		{
			"updated_tender": {
				"name": "",
				"description": "",
				"service_type": "",
				"status": "",
				"organization_id": 0,
				"creator_username": ""
			},
			"message": "updated organization with id=\"2\" not found"
		}`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("EditTender", ctx, 2, tenderToUpdate).Return(models.Tender{}, outerror.ErrOrganizationNotFound)
	router := gin.New()
	router.PATCH("/api/tenders/:tenderId/edit", svc.EditTender(ctx))
	req := httptest.NewRequest(http.MethodPatch, "/api/tenders/2/edit", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestEditTender_FailEmployeeNotResponsibleForOrg проверяет, что
// если обновленный сотрудник не ответсвенный за новую организацию, то
// возвращается ошибка и код 400.
func TestEditTender_FailEmployeeNotResponsibleForOrg(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	tenderName := "update Tender 1"
	description := "update qwe"
	serviceType := "update op"
	status := "update open"
	organizationId := 2
	creatorUsername := "update qwe"

	tenderToUpdate := models.TenderToUpdate{
		TenderName:      &tenderName,
		Description:     &description,
		ServiceType:     &serviceType,
		Status:          &status,
		OrganizationId:  &organizationId,
		CreatorUsername: &creatorUsername,
	}
	reqBody := `
		{
			"update_tender_data": {
				"name": "update Tender 1",
				"description": "update qwe",
				"service_type": "update op",
				"status": "update open",
				"organization_id": 2,
				"creator_username": "update qwe"
			}
		}`

	expectedBody := `
		{
			"updated_tender": {
				"name": "",
				"description": "",
				"service_type": "",
				"status": "",
				"organization_id": 0,
				"creator_username": ""
			},
			"message": "employee with username=\"update qwe\" not responsible for organization with id=\"2\""
		}`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("EditTender", ctx, 2, tenderToUpdate).Return(models.Tender{}, outerror.ErrEmployeeNotResponsibleForOrganization)
	router := gin.New()
	router.PATCH("/api/tenders/:tenderId/edit", svc.EditTender(ctx))
	req := httptest.NewRequest(http.MethodPatch, "/api/tenders/2/edit", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

func TestEditTender_FailUnknownTenderStatus(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	tenderName := "update Tender 1"
	description := "update qwe"
	serviceType := "update op"
	status := "update open"
	organizationId := 2
	creatorUsername := "update qwe"

	tenderToUpdate := models.TenderToUpdate{
		TenderName:      &tenderName,
		Description:     &description,
		ServiceType:     &serviceType,
		Status:          &status,
		OrganizationId:  &organizationId,
		CreatorUsername: &creatorUsername,
	}
	reqBody := `
		{
			"update_tender_data": {
				"name": "update Tender 1",
				"description": "update qwe",
				"service_type": "update op",
				"status": "update open",
				"organization_id": 2,
				"creator_username": "update qwe"
			}
		}`

	expectedBody := `
		{
			"updated_tender": {
				"name": "",
				"description": "",
				"service_type": "",
				"status": "",
				"organization_id": 0,
				"creator_username": ""
			},
			"message": "tender status \"update open\" is unknown"
		}`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("EditTender", ctx, 2, tenderToUpdate).Return(models.Tender{}, outerror.ErrUnknownTenderStatus)
	router := gin.New()
	router.PATCH("/api/tenders/:tenderId/edit", svc.EditTender(ctx))
	req := httptest.NewRequest(http.MethodPatch, "/api/tenders/2/edit", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

func TestEditTender_FailCannotUpdateTender(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	tenderName := "update Tender 1"
	description := "update qwe"
	serviceType := "update op"
	status := "update open"
	organizationId := 2
	creatorUsername := "update qwe"

	tenderToUpdate := models.TenderToUpdate{
		TenderName:      &tenderName,
		Description:     &description,
		ServiceType:     &serviceType,
		Status:          &status,
		OrganizationId:  &organizationId,
		CreatorUsername: &creatorUsername,
	}
	reqBody := `
		{
			"update_tender_data": {
				"name": "update Tender 1",
				"description": "update qwe",
				"service_type": "update op",
				"status": "update open",
				"organization_id": 2,
				"creator_username": "update qwe"
			}
		}`

	expectedBody := `
		{
			"updated_tender": {
				"name": "",
				"description": "",
				"service_type": "",
				"status": "",
				"organization_id": 0,
				"creator_username": ""
			},
			"message": "internal error"
		}`
	svc := tenderapi.New(logger, mockTenderService)
	someErr := errors.New("some error")
	mockTenderService.On("EditTender", ctx, 2, tenderToUpdate).Return(models.Tender{}, someErr)
	router := gin.New()
	router.PATCH("/api/tenders/:tenderId/edit", svc.EditTender(ctx))
	req := httptest.NewRequest(http.MethodPatch, "/api/tenders/2/edit", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}
