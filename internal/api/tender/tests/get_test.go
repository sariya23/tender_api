package tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
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

// TestGetAllTenders_Success проверяет
// успешный сценарий вызова хендлера GetTenders.
//
// Возвращается код 200 и тело со списком тендеров.
func TestGetAllTenders_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)
	mockTenders := []models.Tender{
		{TenderName: "Tender 1", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 1, CreatorUsername: "qwe"},
		{TenderName: "Tender 1", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 1, CreatorUsername: "qwe"},
	}
	expectedBody := `
	{
		"tenders":[
			{"name":"Tender 1", "description": "qwe", "service_type": "op", "status": "open", "organization_id": 1, "creator_username": "qwe"},
			{"name":"Tender 1", "description": "qwe", "service_type": "op", "status": "open", "organization_id": 1, "creator_username": "qwe"}
		],"message":"ok"
	}
	`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("GetTenders", ctx, "all").Return(mockTenders, nil)
	req := httptest.NewRequest(http.MethodGet, "/tenders?srv_type=all", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler := svc.GetTenders(ctx)
	handler(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestGetAllTenders_FailTendersNotFound проверяет, что
// если нет тендеров с указанным service type, то возвращается
// код 400 и сообщение с ошибкой.
func TestGetAllTenders_FailTendersNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)
	mockTenders := []models.Tender{}
	expectedBody := `
	{
		"message":"no tenders found with service type: qwe"
	}
	`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("GetTenders", ctx, "qwe").Return(mockTenders, outerror.ErrTendersWithThisServiceTypeNotFound)
	req := httptest.NewRequest(http.MethodGet, "/tenders?srv_type=qwe", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler := svc.GetTenders(ctx)
	handler(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestGetAllTenders_FailinternalError проверяет, что
// если произошла какая-то внутренняя ошибка, то возвращается
// код 500 и сообщение
func TestGetAllTenders_FailinternalError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)
	mockTenders := []models.Tender{}
	someErr := errors.New("some error")
	expectedBody := `
	{
		"message":"internal error"
	}
	`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("GetTenders", ctx, "qwe").Return(mockTenders, someErr)
	req := httptest.NewRequest(http.MethodGet, "/tenders?srv_type=qwe", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler := svc.GetTenders(ctx)
	handler(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestGetEmployeeTenders_Success проверяет, что
// в случае успешного получения тендеров сотрудника,
// возвращается список этих тендеров.
func TestGetEmployeeTenders_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)
	username := "qwe"
	mockTenders := []models.Tender{
		{TenderName: "Tender 1", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 1, CreatorUsername: username},
		{TenderName: "Tender 1", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 1, CreatorUsername: username},
	}
	expectedBody := `
	{
		"tenders":[
			{"name":"Tender 1", "description": "qwe", "service_type": "op", "status": "open", "organization_id": 1, "creator_username": "qwe"},
			{"name":"Tender 1", "description": "qwe", "service_type": "op", "status": "open", "organization_id": 1, "creator_username": "qwe"}
		],"message":"ok"
	}
	`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("GetEmployeeTendersByUsername", ctx, username).Return(mockTenders, nil)
	req := httptest.NewRequest(http.MethodGet, "/tenders/my?username=qwe", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler := svc.GetEmployeeTendersByUsername(ctx)
	handler(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestGetEmployeeTenders_SuccessRedirect проверяет, что
// в случае, когда username пустой, то происходит редирект на
// /api/tenders.
func TestGetEmployeeTenders_SuccessRedirect(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)
	svc := tenderapi.New(logger, mockTenderService)

	req := httptest.NewRequest(http.MethodGet, "/tenders/my?username=", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler := svc.GetEmployeeTendersByUsername(ctx)
	handler(c)

	// Assert
	assert.Equal(t, http.StatusMovedPermanently, w.Code)
	require.Equal(t, w.Header().Values("location"), []string{"/api/tenders"})
}

// TestGetEmployeeTenders_FailEmployeeNotFound проверяет, что
// если сотрудника с переданным username нет, то возвращается ошибка.
func TestGetEmployeeTenders_FailEmployeeNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)
	username := "qwe"
	expectedBody := `
	{
		"message": "employee with username \"qwe\" not found"
	}
	`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("GetEmployeeTendersByUsername", ctx, username).Return([]models.Tender{}, outerror.ErrEmployeeNotFound)
	req := httptest.NewRequest(http.MethodGet, "/tenders/my?username=qwe", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler := svc.GetEmployeeTendersByUsername(ctx)
	handler(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestGetEmployeeTenders_FailEmployeeTendersNotFound проверяет, что
// если у переданного сотрудника нет тендеров, то возвращается ошибка.
func TestGetEmployeeTenders_FailEmployeeTendersNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)
	username := "qwe"
	expectedBody := `
	{
		"message": "not found tenders for employee with username \"qwe\""
	}
	`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("GetEmployeeTendersByUsername", ctx, username).Return([]models.Tender{}, outerror.ErrEmployeeTendersNotFound)
	req := httptest.NewRequest(http.MethodGet, "/tenders/my?username=qwe", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler := svc.GetEmployeeTendersByUsername(ctx)
	handler(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}
