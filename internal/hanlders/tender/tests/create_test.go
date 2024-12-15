package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sariya23/tender/internal/domain/models"
	schema "github.com/sariya23/tender/internal/hanlders"
	tenderapi "github.com/sariya23/tender/internal/hanlders/tender"
	"github.com/sariya23/tender/internal/hanlders/tender/mocks"
	"github.com/sariya23/tender/internal/lib/logger/slogdiscard"
	outerror "github.com/sariya23/tender/internal/out_error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateTender_Success
func TestCreateTender_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)
	mockTender := models.Tender{
		TenderName: "Tender 1", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 1, CreatorUsername: "qwe",
	}
	reqBody := `
	{
		"tender": {
			"name": "Tender 1",
			"description": "qwe",
			"service_type": "op",
			"status": "open",
			"organization_id": 1,
			"creator_username": "qwe"
		}
	}`

	expectedBody := `
	{
		"tender": {
			"name": "Tender 1",
			"description": "qwe",
			"service_type": "op",
			"status": "open",
			"organization_id": 1,
			"creator_username": "qwe"
		},
		"message": "ok"
	}`
	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("CreateTender", ctx, mockTender).Return(mockTender, nil)
	req := httptest.NewRequest(http.MethodPost, "/tenders/new", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler := svc.CreateTender(ctx)
	handler(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, expectedBody, w.Body.String())
}

// TestCreateTender_FailUnmurshalError проверяет, что
// если прилетает невалидный json, то тендер не создается
// и возвращается ошибка.
func TestCreateTender_FailUnmurshalSyntaxError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	reqBody := `
	{
		"tender": {
			"name": "Tender 1",
			"description": "qwe",
			"service_type": "op",
			"status": "open",
			"organization_id": 1
			"creator_username": "qwe"
		}
	}`

	svc := tenderapi.New(logger, mockTenderService)

	req := httptest.NewRequest(http.MethodPost, "/tenders/new", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler := svc.CreateTender(ctx)
	handler(c)

	// Assert
	var resp schema.CreateTenderResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	require.Equal(t, resp.Tender, models.Tender{})
	require.Contains(t, resp.Message, "json syntax")
}

// TestCreateTender_FailUnmurshalTypeError проверяет,
// что если прилетит неверный тип в json, то тендер не создаться
// и вернеться ошибка.
func TestCreateTender_FailUnmurshalTypeError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	reqBody := `
	{
		"tender": {
			"name": "Tender 1",
			"description": "qwe",
			"service_type": "op",
			"status": "open",
			"organization_id": "qwe",
			"creator_username": "qwe"
		}
	}`

	svc := tenderapi.New(logger, mockTenderService)

	req := httptest.NewRequest(http.MethodPost, "/tenders/new", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler := svc.CreateTender(ctx)
	handler(c)

	// Assert
	var resp schema.CreateTenderResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	require.Equal(t, resp.Tender, models.Tender{})
	require.Contains(t, resp.Message, "json type")
}

// TestCreateTender_FailNegativeOrgID проверяет, что
// если прилетит отрицательный id организации, то тендер не создаться
// и вернеться ошибка.
func TestCreateTender_FailNegativeOrgID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	reqBody := `
	{
		"tender": {
			"name": "Tender 1",
			"description": "qwe",
			"service_type": "op",
			"status": "open",
			"organization_id": -1000,
			"creator_username": "qwe"
		}
	}`

	svc := tenderapi.New(logger, mockTenderService)

	req := httptest.NewRequest(http.MethodPost, "/tenders/new", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler := svc.CreateTender(ctx)
	handler(c)

	// Assert
	var resp schema.CreateTenderResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	require.Equal(t, resp.Tender, models.Tender{})
	require.Contains(t, resp.Message, "validation failed")
}

// TestCreateTender_FailEmployeeNotFound проверяет, что
// если пользователя не существует, то тендер не создается и возвращается
// сообщение с ошибкой.
func TestCreateTender_FailEmployeeNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)
	creatorUsername := "qwe"
	mockTender := models.Tender{
		TenderName:      "Tender 1",
		Description:     "qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: creatorUsername,
	}
	reqBody := `
	{
		"tender": {
			"name": "Tender 1",
			"description": "qwe",
			"service_type": "op",
			"status": "open",
			"organization_id": 1,
			"creator_username": "qwe"
		}
	}`

	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("CreateTender", ctx, mockTender).Return(models.Tender{}, outerror.ErrEmployeeNotFound)
	req := httptest.NewRequest(http.MethodPost, "/tenders/new", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler := svc.CreateTender(ctx)
	handler(c)

	// Assert
	var resp schema.CreateTenderResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, models.Tender{}, resp.Tender)
	require.Equal(t, "employee with username=\"qwe\" not found", resp.Message)
}

// TestCreateTender_FailOrganizationNotFound проверяет, что
// если организации не существует, что тендер не создается и возвращается
// сообщение с ошибкой.
func TestCreateTender_FailOrganizationNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)
	orgId := 1
	mockTender := models.Tender{
		TenderName:      "Tender 1",
		Description:     "qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  orgId,
		CreatorUsername: "qwe",
	}
	reqBody := `
	{
		"tender": {
			"name": "Tender 1",
			"description": "qwe",
			"service_type": "op",
			"status": "open",
			"organization_id": 1,
			"creator_username": "qwe"
		}
	}`

	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("CreateTender", ctx, mockTender).Return(models.Tender{}, outerror.ErrOrganizationNotFound)
	req := httptest.NewRequest(http.MethodPost, "/tenders/new", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler := svc.CreateTender(ctx)
	handler(c)

	// Assert
	var resp schema.CreateTenderResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, models.Tender{}, resp.Tender)
	require.Equal(t, "organization with id=1 not found", resp.Message)
}

// TestCreateTender_FailUserNotResponsibleForOrganization
func TestCreateTender_FailEmployeerNotResponsibleForOrganization(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	mockTender := models.Tender{
		TenderName:      "Tender 1",
		Description:     "qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: "qwe",
	}
	reqBody := `
	{
		"tender": {
			"name": "Tender 1",
			"description": "qwe",
			"service_type": "op",
			"status": "open",
			"organization_id": 1,
			"creator_username": "qwe"
		}
	}`

	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("CreateTender", ctx, mockTender).Return(models.Tender{}, outerror.ErrEmployeeNotResponsibleForOrganization)
	req := httptest.NewRequest(http.MethodPost, "/tenders/new", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler := svc.CreateTender(ctx)
	handler(c)

	// Assert
	var resp schema.CreateTenderResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, models.Tender{}, resp.Tender)
	require.Equal(t, "employee \"qwe\" not responsible for organization with id=1", resp.Message)
}

func TestCreateTender_FailWrongTenderStatus(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	logger := slogdiscard.NewDiscardLogger()
	mockTenderService := new(mocks.MockTenderServiceProvider)

	mockTender := models.Tender{
		TenderName:      "Tender 1",
		Description:     "qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: "qwe",
	}
	reqBody := `
	{
		"tender": {
			"name": "Tender 1",
			"description": "qwe",
			"service_type": "op",
			"status": "open",
			"organization_id": 1,
			"creator_username": "qwe"
		}
	}`

	svc := tenderapi.New(logger, mockTenderService)

	mockTenderService.On("CreateTender", ctx, mockTender).Return(models.Tender{}, outerror.ErrNewTenderCannotCreatedWithStatusNotCreated)
	req := httptest.NewRequest(http.MethodPost, "/tenders/new", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Act
	handler := svc.CreateTender(ctx)
	handler(c)

	// Assert
	var resp schema.CreateTenderResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	require.Equal(t, models.Tender{}, resp.Tender)
	require.Equal(t, "cannot create tender with status \"open\"", resp.Message)
}
