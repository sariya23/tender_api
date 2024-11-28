package tender_test

import (
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
