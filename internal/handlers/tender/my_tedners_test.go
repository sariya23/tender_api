package tender_test

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"tender/internal/domain/models"
// 	"tender/internal/handlers/tender"
// 	"tender/internal/lib/logger/slogdiscard"
// 	"tender/internal/storage"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// )

// type MockUserTenderGetter struct {
// 	mock.Mock
// }

// func (m *MockUserTenderGetter) UserTenders(ctx context.Context, username string) ([]models.Tender, error) {
// 	args := m.Called(ctx, username)
// 	return args.Get(0).([]models.Tender), args.Error(1)
// }

// // TestGetUserTenders_Success проверяет, что
// // пользователь, указав свой username в query-параметре
// // получит только те тендреры, которые он создавал.
// //
// // Ожидаемое поведение:
// //
// // - Код 200;
// //
// // - Список тенедеров и сообщенеи ok.
// func TestGetUserTenders_Success(t *testing.T) {
// 	ctx := context.Background()
// 	mockUserTenderGetter := new(MockUserTenderGetter)
// 	logger := slogdiscard.NewDiscardLogger()
// 	mockTenders := []models.Tender{
// 		{TenderName: "Tender 1", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 1, CreatorUsername: "qwe"},
// 		{TenderName: "Tender 2", Description: "qwe", ServiceType: "op", Status: "open", OrganizationId: 2, CreatorUsername: "qwe"},
// 	}
// 	h := tender.GetUserTenders(ctx, logger, mockUserTenderGetter)
// 	gin.SetMode(gin.TestMode)
// 	router := gin.New()
// 	router.GET("/api/tender/my", h)
// 	mockUserTenderGetter.On("UserTenders", ctx, "qwe").Return(mockTenders, nil)

// 	req := httptest.NewRequest(http.MethodGet, "/api/tender/my?username=qwe", nil)
// 	resp := httptest.NewRecorder()
// 	router.ServeHTTP(resp, req)

// 	assert.Equal(t, http.StatusOK, resp.Code)
// 	var responseBody tender.GetUserTendersResponse
// 	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
// 	require.NoError(t, err)
// 	require.Equal(t, "ok", responseBody.Message)
// 	require.Equal(t, mockTenders, responseBody.Tenders)
// }

// // TestGetUserTenders_Redirect проверяет, что в случае
// // отсутсвия квери параметра username происходит редирект на /api/tender/, где
// // отобразятся все тендеры.
// func TestGetUserTenders_Redirect(t *testing.T) {
// 	ctx := context.Background()
// 	mockUserTenderGetter := new(MockUserTenderGetter)
// 	logger := slogdiscard.NewDiscardLogger()
// 	h := tender.GetUserTenders(ctx, logger, mockUserTenderGetter)
// 	gin.SetMode(gin.TestMode)
// 	router := gin.New()
// 	router.GET("/api/tender/my", h)

// 	req := httptest.NewRequest(http.MethodGet, "/api/tender/my", nil)
// 	resp := httptest.NewRecorder()
// 	router.ServeHTTP(resp, req)

// 	assert.Equal(t, http.StatusMovedPermanently, resp.Code)

// 	location := resp.Header().Get("Location")
// 	require.Equal(t, "/api/tender/", location)
// }

// // TestGetUserTenders_NoTendersForUser проверяет, что
// // если для указанного пользователя нет тендеров, то возвращается
// // код 200 и сообщение.
// func TestGetUserTenders_NoTendersForUser(t *testing.T) {
// 	ctx := context.Background()
// 	mockUserTenderGetter := new(MockUserTenderGetter)
// 	logger := slogdiscard.NewDiscardLogger()
// 	mockTenders := []models.Tender{}
// 	h := tender.GetUserTenders(ctx, logger, mockUserTenderGetter)
// 	gin.SetMode(gin.TestMode)
// 	router := gin.New()
// 	router.GET("/api/tender/my", h)
// 	mockUserTenderGetter.On("UserTenders", ctx, "qwe").Return(mockTenders, storage.ErrNoTenderForThisUser)

// 	req := httptest.NewRequest(http.MethodGet, "/api/tender/my?username=qwe", nil)
// 	resp := httptest.NewRecorder()
// 	router.ServeHTTP(resp, req)

// 	assert.Equal(t, http.StatusOK, resp.Code)
// 	var responseBody tender.GetUserTendersResponse
// 	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
// 	require.NoError(t, err)
// 	require.Equal(t, "no tenders for this user", responseBody.Message)
// 	require.Empty(t, responseBody.Tenders)
// }
