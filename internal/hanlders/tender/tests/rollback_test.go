package tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sariya23/tender/internal/domain/models"
	tenderapi "github.com/sariya23/tender/internal/hanlders/tender"
	"github.com/sariya23/tender/internal/hanlders/tender/mocks"
	"github.com/sariya23/tender/internal/lib/logger/slogdiscard"
	outerror "github.com/sariya23/tender/internal/out_error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRollbackTender_Success
func TestRollbackTender_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	mockTender := models.Tender{
		TenderName:      "qwe",
		Description:     "qwe",
		Status:          "qwe",
		ServiceType:     "qwe",
		OrganizationId:  1,
		CreatorUsername: "qwe",
	}

	expectedBody := `
		{
			"rollback_tender": {
				"name": "qwe",
				"description": "qwe",
				"service_type": "qwe",
				"status": "qwe",
				"organization_id": 1,
				"creator_username": "qwe"
			},
			"message": "ok"
		}`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("RollbackTender", ctx, 2, 3).Return(mockTender, nil)
	router := gin.New()
	router.PUT("/api/tenders/:tenderId/rollback/:version", svc.RollbackTender(ctx))
	req := httptest.NewRequest(http.MethodPut, "/api/tenders/2/rollback/3", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestRollbackTender_FailTenderIdIsNotInt проверяет, что
// если параметр tenderId в пути не получается конвертировать в int, то
// возвращается ошибка и код 400.
func TestRollbackTender_FailTenderIdIsNotInt(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	expectedBody := `
		{
			"rollback_tender": {
				"name": "",
				"description": "",
				"service_type": "",
				"status": "",
				"organization_id": 0,
				"creator_username": ""
			},
			"message": "tenderId must be positive integer number"
		}`
	svc := tenderapi.New(logger, mockTenderService)

	router := gin.New()
	router.PUT("/api/tenders/:tenderId/rollback/:version", svc.RollbackTender(ctx))
	req := httptest.NewRequest(http.MethodPut, "/api/tenders/2.34/rollback/3", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestRollbackTender_FailVersionIsNotInt проверяет, что
// если параметр version в пути не получается конвертировать в int,
// то возвращается ошибка и код 400.
func TestRollbackTender_FailVersionIsNotInt(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	expectedBody := `
		{
			"rollback_tender": {
				"name": "",
				"description": "",
				"service_type": "",
				"status": "",
				"organization_id": 0,
				"creator_username": ""
			},
			"message": "version must be positive integer number"
		}`
	svc := tenderapi.New(logger, mockTenderService)

	router := gin.New()
	router.PUT("/api/tenders/:tenderId/rollback/:version", svc.RollbackTender(ctx))
	req := httptest.NewRequest(http.MethodPut, "/api/tenders/2/rollback/qwe", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestRollbackTender_FailTenderNotFound проверяет, что
// если нет тендера с id, который указан в пути, то возвращается
// ошибка и код 400.
func TestRollbackTender_FailTenderNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	expectedBody := `
		{
			"rollback_tender": {
				"name": "",
				"description": "",
				"service_type": "",
				"status": "",
				"organization_id": 0,
				"creator_username": ""
			},
			"message": "tender with id=<2> not found"
		}`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("RollbackTender", ctx, 2, 3).Return(models.Tender{}, outerror.ErrTenderNotFound)
	router := gin.New()
	router.PUT("/api/tenders/:tenderId/rollback/:version", svc.RollbackTender(ctx))
	req := httptest.NewRequest(http.MethodPut, "/api/tenders/2/rollback/3", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestRollbackTender_FailVersionNotFound проверяет, что
// если у тендера нет версии, которая указана в пути, то
// возвращается ошибка и код 400.
func TestRollbackTender_FailVersionNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	expectedBody := `
		{
			"rollback_tender": {
				"name": "",
				"description": "",
				"service_type": "",
				"status": "",
				"organization_id": 0,
				"creator_username": ""
			},
			"message": "tender with id=<2> doesnt have version=<3>"
		}`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("RollbackTender", ctx, 2, 3).Return(models.Tender{}, outerror.ErrTenderVersionNotFound)
	router := gin.New()
	router.PUT("/api/tenders/:tenderId/rollback/:version", svc.RollbackTender(ctx))
	req := httptest.NewRequest(http.MethodPut, "/api/tenders/2/rollback/3", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestRollbackTender_FailInternalError проверяет, что
// если произошла какая-то внутренняя ошибка, то возвращается
// код 500.
func TestRollbackTender_FailInternalError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	expectedBody := `
		{
			"rollback_tender": {
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
	someErr := errors.New("some err")
	mockTenderService.On("RollbackTender", ctx, 2, 3).Return(models.Tender{}, someErr)
	router := gin.New()
	router.PUT("/api/tenders/:tenderId/rollback/:version", svc.RollbackTender(ctx))
	req := httptest.NewRequest(http.MethodPut, "/api/tenders/2/rollback/3", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}
