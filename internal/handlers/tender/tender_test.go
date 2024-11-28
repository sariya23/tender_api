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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTenderGetter struct {
	mock.Mock
}

func (m *MockTenderGetter) Tenders(ctx context.Context) ([]models.Tender, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Tender), args.Error(1)
}

func (m *MockTenderGetter) TendersByServiceType(ctx context.Context, serviceType string) ([]models.Tender, error) {
	args := m.Called(ctx, serviceType)
	return args.Get(0).([]models.Tender), args.Error(1)
}

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

type MockUserTenderGetter struct {
	mock.Mock
}

func (m *MockUserTenderGetter) UserTenders(ctx context.Context, username string) ([]models.Tender, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]models.Tender), args.Error(1)
}

// TestGetAllTendersSuccess проверяет, что
// хендлер отдает все тендеры по запросу на
// api/tender/.
// Ожидаемое поведение:
//
// - Код ответа 200;
//
// - В теле ответа json: {"message": "ok", "tenders": <список тендеров>, "service_type": "all"}
func TestGetAllTendersSuccess(t *testing.T) {
	ctx := context.Background()
	mockTenderGetter := new(MockTenderGetter)
	logger := slogdiscard.NewDiscardLogger()
	mockTenders := []models.Tender{
		{TenderName: "Tender 1", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 1, CreatorUsername: "qwe"},
		{TenderName: "Tender 2", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 2, CreatorUsername: "zxc"},
	}
	h := tender.GetTenders(ctx, logger, mockTenderGetter)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/tender/", h)

	mockTenderGetter.On("Tenders", ctx).Return(mockTenders, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/tender/", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var responseBody tender.GetTendersResponse
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)
	require.Equal(t, mockTenders, responseBody.Tenders)
	require.Equal(t, "ok", responseBody.Message)
	require.Equal(t, "all", responseBody.ServiceType)

	mockTenderGetter.AssertCalled(t, "Tenders", ctx)
}

// TestGetAllTendersByServiceTypeSuccess - проверяет, что
// при наличии квери параметра srv_type возвращаются тендеры
// с только указанной сферой услуг.
//
// Ожидаемое поведение:
//
// - Код 200;
//
// - В теле ответа json {"message": "ok", "tenders": <список тендеров>, "service_type": "op"}
func TestGetAllTendersByServiceTypeSuccess(t *testing.T) {
	ctx := context.Background()
	mockTenderGetter := new(MockTenderGetter)
	logger := slogdiscard.NewDiscardLogger()
	mockTenders := []models.Tender{
		{TenderName: "Tender 1", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 1, CreatorUsername: "qwe"},
		{TenderName: "Tender 2", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 2, CreatorUsername: "zxc"},
	}
	h := tender.GetTenders(ctx, logger, mockTenderGetter)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/tender/", h)

	mockTenderGetter.On("TendersByServiceType", ctx, "op").Return(mockTenders, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/tender/?srv_type=op", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var responseBody tender.GetTendersResponse
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)
	require.Equal(t, mockTenders, responseBody.Tenders)
	require.Equal(t, "ok", responseBody.Message)
	require.Equal(t, "op", responseBody.ServiceType)

	mockTenderGetter.AssertCalled(t, "TendersByServiceType", ctx, "op")
}

// TestNoTendersFound проверяет случай, когда не нашлось подходящих
// тендеров.
//
// Ожидаемое поведение:
//
// - Код 200;
//
// - В теле ответа json {"message": "no tenders found", "tenders": <empty>, "service_type": "all"}
func TestNoTendersFound(t *testing.T) {
	ctx := context.Background()
	mockTenderGetter := new(MockTenderGetter)
	logger := slogdiscard.NewDiscardLogger()
	mockTenders := []models.Tender{}
	h := tender.GetTenders(ctx, logger, mockTenderGetter)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/tender/", h)

	mockTenderGetter.On("Tenders", mock.Anything).Return(mockTenders, storage.ErrNoTenderPresence)
	req := httptest.NewRequest(http.MethodGet, "/api/tender/", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var responseBody tender.GetTendersResponse
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)
	require.Empty(t, mockTenders, responseBody.Tenders)
	require.Equal(t, "no tenders found", responseBody.Message)
	require.Equal(t, "all", responseBody.ServiceType)
	mockTenderGetter.AssertCalled(t, "Tenders", ctx)
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

	require.Equal(t, http.StatusOK, resp.Code)

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

	require.Equal(t, http.StatusBadRequest, resp.Code)

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

	require.Equal(t, http.StatusBadRequest, resp.Code)

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

	require.Equal(t, http.StatusBadRequest, resp.Code)

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

	require.Equal(t, http.StatusBadRequest, resp.Code)

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

	require.Equal(t, http.StatusBadRequest, resp.Code)

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

	require.Equal(t, http.StatusBadRequest, resp.Code)

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

	require.Equal(t, http.StatusBadRequest, resp.Code)

	var responseBody tender.CreateTenderResponse
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.Equal(t, models.Tender{}, responseBody.Tender)
	require.Equal(t, "user not responsible for organization", responseBody.Message)
}

// TestGetUserTenders_Success проверяет, что
// пользователь, указав свой username в query-параметре
// получит только те тендреры, которые он создавал.
//
// Ожидаемое поведение:
//
// - Код 200;
//
// - Список тенедеров и сообщенеи ok.
func TestGetUserTenders_Success(t *testing.T) {
	ctx := context.Background()
	mockUserTenderGetter := new(MockUserTenderGetter)
	logger := slogdiscard.NewDiscardLogger()
	mockTenders := []models.Tender{
		{TenderName: "Tender 1", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 1, CreatorUsername: "qwe"},
		{TenderName: "Tender 2", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 2, CreatorUsername: "qwe"},
	}
	h := tender.GetUserTenders(ctx, logger, mockUserTenderGetter)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/tender/my", h)
	mockUserTenderGetter.On("UserTenders", ctx, "qwe").Return(mockTenders, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/tender/my?username=qwe", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var responseBody tender.GetUserTendersResponse
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)
	require.Equal(t, "ok", responseBody.Message)
	require.Equal(t, mockTenders, responseBody.Tenders)
}

// TestGetUserTenders_Redirect проверяет, что в случае
// отсутсвия квери параметра username происходит редирект на /api/tender/, где
// отобразятся все тендеры.
func TestGetUserTenders_Redirect(t *testing.T) {
	ctx := context.Background()
	mockUserTenderGetter := new(MockUserTenderGetter)
	logger := slogdiscard.NewDiscardLogger()
	h := tender.GetUserTenders(ctx, logger, mockUserTenderGetter)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/tender/my", h)

	req := httptest.NewRequest(http.MethodGet, "/api/tender/my", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusMovedPermanently, resp.Code)

	location := resp.Header().Get("Location")
	require.Equal(t, "/api/tender/", location)
}

// TestGetUserTenders_NoTendersForUser проверяет, что
// если для указанного пользователя нет тендеров, то возвращается
// код 200 и сообщение.
func TestGetUserTenders_NoTendersForUser(t *testing.T) {
	ctx := context.Background()
	mockUserTenderGetter := new(MockUserTenderGetter)
	logger := slogdiscard.NewDiscardLogger()
	mockTenders := []models.Tender{}
	h := tender.GetUserTenders(ctx, logger, mockUserTenderGetter)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/tender/my", h)
	mockUserTenderGetter.On("UserTenders", ctx, "qwe").Return(mockTenders, storage.ErrNoTenderForThisUser)

	req := httptest.NewRequest(http.MethodGet, "/api/tender/my?username=qwe", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	var responseBody tender.GetUserTendersResponse
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)
	require.Equal(t, "no tenders for this user", responseBody.Message)
	require.Empty(t, responseBody.Tenders)
}
