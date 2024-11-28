package tender_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"tender/internal/domain/models"
	"tender/internal/handlers/tender"
	"tender/internal/lib/logger/slogdiscard"
	"tender/internal/storage"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTenderCreater struct {
	mock.Mock
}

func (m *MockTenderCreater) CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	args := m.Called(ctx, tender)
	return args.Get(0).(models.Tender), args.Error(1)
}

type MockUserProvider struct {
	mock.Mock
}

func (m *MockUserProvider) UserByUsername(ctx context.Context, username string) error {
	args := m.Called(ctx, username)
	return args.Error(0)
}

type MockOrganizationProvider struct {
	mock.Mock
}

func (m *MockOrganizationProvider) OrganizationById(ctx context.Context, ogranizationId int) error {
	args := m.Called(ctx, ogranizationId)
	return args.Error(0)
}

type MockUserResponsibler struct {
	mock.Mock
}

func (m *MockUserResponsibler) CheckUserResponsible(ctx context.Context, username string, organizationId int) error {
	args := m.Called(ctx, username, organizationId)
	return args.Error(0)
}

// TestCreateTender_Success проверяет, что
// хендлер создает тенедер. Ожидаемое поведение:
//
// - Код 200
//
// - Возвращает json с созданным тендером и сообщение ok.
func TestCreateTender_Success(t *testing.T) {
	ctx := context.Background()
	mockTenderCreater := new(MockTenderCreater)
	mockUserProvider := new(MockUserProvider)
	mockOrganizationProvider := new(MockOrganizationProvider)
	mockUserResponsibler := new(MockUserResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	mockTender := models.Tender{
		TenderName:      "Test",
		Description:     "qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: "sariya",
	}
	h := tender.CreateTender(ctx, logger, mockTenderCreater, mockUserProvider, mockOrganizationProvider, mockUserResponsibler)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/tender/new", h)

	mockTenderCreater.On("CreateTender", ctx, mockTender).Return(mockTender, nil)
	mockUserProvider.On("UserByUsername", ctx, "sariya").Return(nil)
	mockOrganizationProvider.On("OrganizationById", ctx, 1).Return(nil)
	mockUserResponsibler.On("CheckUserResponsible", ctx, "sariya", 1).Return(nil)

	reqData := tender.CreateTenderRequest{Tender: mockTender}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqData)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/tender/new", &buf)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var responseBody tender.CreateTenderResponse
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.Equal(t, mockTender, responseBody.Tender)
	require.Equal(t, "ok", responseBody.Message)
}

// TestCannotCreateTender_BlankBody проверяет, что при
// отправке пустого тела запроса возвращается код 400 и
// сообщение empty json body.
func TestCannotCreateTender_BlankBody(t *testing.T) {
	ctx := context.Background()
	mockTenderCreater := new(MockTenderCreater)
	mockUserProvider := new(MockUserProvider)
	mockOrganizationProvider := new(MockOrganizationProvider)
	mockUserResponsibler := new(MockUserResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	h := tender.CreateTender(ctx, logger, mockTenderCreater, mockUserProvider, mockOrganizationProvider, mockUserResponsibler)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/tender/new", h)

	req := httptest.NewRequest(http.MethodPost, "/api/tender/new", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var responseBody tender.CreateTenderResponse
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.Equal(t, models.Tender{}, responseBody.Tender)
	require.Equal(t, "empty json body", responseBody.Message)
}

// TestCannotCreateTender_InvalidTypes проверяет, что если
// прилетят неверные типы, то вернется код 400 и сообщение
// с ошибкой.
func TestCannotCreateTender_InvalidTypes(t *testing.T) {
	ctx := context.Background()
	mockTenderCreater := new(MockTenderCreater)
	mockUserProvider := new(MockUserProvider)
	mockOrganizationProvider := new(MockOrganizationProvider)
	mockUserResponsibler := new(MockUserResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	h := tender.CreateTender(ctx, logger, mockTenderCreater, mockUserProvider, mockOrganizationProvider, mockUserResponsibler)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/tender/new", h)

	reqData := struct {
		TenderName      int    `json:"name"`
		Description     string `json:"description"`
		ServiceType     string `json:"serviceType"`
		Status          string `json:"status"`
		OrganizationId  int    `json:"organizationId"`
		CreatorUsername string `json:"creatorUsername"`
	}{
		TenderName:      1,
		Description:     "qwe",
		ServiceType:     "qwe",
		Status:          "qwe",
		OrganizationId:  2,
		CreatorUsername: "qwe",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqData)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/tender/new", &buf)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var responseBody tender.CreateTenderResponse
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.Equal(t, models.Tender{}, responseBody.Tender)
	require.Equal(t, "invalid fields", responseBody.Message)
}

// TestCannotCreateTender_InvalidRequest проверяет, что
// при передаче в теле json с некорректными полями, возвращается
// код 400 и сообщение с ошибкой.
func TestCannotCreateTender_InvalidRequest(t *testing.T) {
	ctx := context.Background()
	mockTenderCreater := new(MockTenderCreater)
	mockUserProvider := new(MockUserProvider)
	mockOrganizationProvider := new(MockOrganizationProvider)
	mockUserResponsibler := new(MockUserResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	h := tender.CreateTender(ctx, logger, mockTenderCreater, mockUserProvider, mockOrganizationProvider, mockUserResponsibler)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/tender/new", h)

	reqData := struct {
		Name string
		Age  int
	}{
		Name: "qwer",
		Age:  20,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqData)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/tender/new", &buf)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var responseBody tender.CreateTenderResponse
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.Equal(t, models.Tender{}, responseBody.Tender)
	require.Equal(t, "invalid fields", responseBody.Message)
}

// TestCannotCreateTender_NoKeyFields проверяет, что если в
// запросе отсутвуют какие-то ключевые поля, то возвращает код 400
// и сообщение с ошибкой.
func TestCannotCreateTender_NoKeyFields(t *testing.T) {
	ctx := context.Background()
	mockTenderCreater := new(MockTenderCreater)
	mockUserProvider := new(MockUserProvider)
	mockOrganizationProvider := new(MockOrganizationProvider)
	mockUserResponsibler := new(MockUserResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	h := tender.CreateTender(ctx, logger, mockTenderCreater, mockUserProvider, mockOrganizationProvider, mockUserResponsibler)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/tender/new", h)

	reqData := struct {
		TenderName  string `json:"name"`
		Description string `json:"description"`
		ServiceType string `json:"serviceType"`
		Status      string `json:"status"`
	}{
		TenderName:  "qwe",
		Description: "qwe",
		ServiceType: "qwe",
		Status:      "qwe",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqData)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/tender/new", &buf)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var responseBody tender.CreateTenderResponse
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.Equal(t, models.Tender{}, responseBody.Tender)
	require.Equal(t, "invalid fields", responseBody.Message)
}

// TestCannotCreateTenderNoUser проверяет, что если
// пользователя не существует, то создать тенедер не получится.
//
// Возвращает код 400 и сообщение с ошибкой.
func TestCannotCreateTender_NoUser(t *testing.T) {
	ctx := context.Background()
	mockTenderCreater := new(MockTenderCreater)
	mockUserProvider := new(MockUserProvider)
	mockOrganizationProvider := new(MockOrganizationProvider)
	mockUserResponsibler := new(MockUserResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	mockTender := models.Tender{
		TenderName:      "Test",
		Description:     "qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: "sariya",
	}
	h := tender.CreateTender(ctx, logger, mockTenderCreater, mockUserProvider, mockOrganizationProvider, mockUserResponsibler)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/tender/new", h)

	mockUserProvider.On("UserByUsername", ctx, "sariya").Return(storage.ErrUserNotFound)

	reqData := tender.CreateTenderRequest{Tender: mockTender}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqData)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/tender/new", &buf)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var responseBody tender.CreateTenderResponse
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.Equal(t, models.Tender{}, responseBody.Tender)
	require.Equal(t, "user not found", responseBody.Message)
}

// TestCannotCreateTender_NoOrganization проверяет, что
// если организации не существует, то тендер создать не получится.
//
// Возвращается код 400 и сообщение с ошибкой.
func TestCannotCreateTender_NoOrganization(t *testing.T) {
	ctx := context.Background()
	mockTenderCreater := new(MockTenderCreater)
	mockUserProvider := new(MockUserProvider)
	mockOrganizationProvider := new(MockOrganizationProvider)
	mockUserResponsibler := new(MockUserResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	mockTender := models.Tender{
		TenderName:      "Test",
		Description:     "qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: "sariya",
	}
	h := tender.CreateTender(ctx, logger, mockTenderCreater, mockUserProvider, mockOrganizationProvider, mockUserResponsibler)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/tender/new", h)

	mockUserProvider.On("UserByUsername", ctx, "sariya").Return(nil)
	mockOrganizationProvider.On("OrganizationById", ctx, 1).Return(storage.ErrOrganizationNotFound)

	reqData := tender.CreateTenderRequest{Tender: mockTender}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqData)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/tender/new", &buf)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var responseBody tender.CreateTenderResponse
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.Equal(t, models.Tender{}, responseBody.Tender)
	require.Equal(t, "organization not found", responseBody.Message)
}

// TestCannotCreateTender_NoUserResponsobility проверяет, что
// если юзер и организация существует, но они не связаны, то возвращается
// код 400 и сообщенеи об ошибке.
func TestCannotCreateTender_NoUserResponsobility(t *testing.T) {
	ctx := context.Background()
	mockTenderCreater := new(MockTenderCreater)
	mockUserProvider := new(MockUserProvider)
	mockOrganizationProvider := new(MockOrganizationProvider)
	mockUserResponsibler := new(MockUserResponsibler)
	logger := slogdiscard.NewDiscardLogger()
	mockTender := models.Tender{
		TenderName:      "Test",
		Description:     "qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: "sariya",
	}
	h := tender.CreateTender(ctx, logger, mockTenderCreater, mockUserProvider, mockOrganizationProvider, mockUserResponsibler)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/tender/new", h)

	mockUserProvider.On("UserByUsername", ctx, "sariya").Return(nil)
	mockOrganizationProvider.On("OrganizationById", ctx, 1).Return(nil)
	mockUserResponsibler.On("CheckUserResponsible", ctx, "sariya", 1).Return(storage.ErrUserNotReponsibleForOrg)

	reqData := tender.CreateTenderRequest{Tender: mockTender}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqData)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/tender/new", &buf)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var responseBody tender.CreateTenderResponse
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.Equal(t, models.Tender{}, responseBody.Tender)
	require.Equal(t, "user not responsible for organization", responseBody.Message)
}
