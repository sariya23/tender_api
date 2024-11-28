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
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTenderEditor struct {
	mock.Mock
}

func (m *MockTenderEditor) EditTender(ctx context.Context, tenderId int, newTender models.Tender) (models.Tender, error) {
	args := m.Called(ctx, tenderId, newTender)
	return args.Get(0).(models.Tender), args.Error(1)
}

// TestEditTender_Success проверяет, что
// тендер успешно обновляется и возвращаются обновленные данные
// и код 200.
func TestEditTender_Success(t *testing.T) {
	ctx := context.Background()
	logger := slogdiscard.NewDiscardLogger()
	mockTenderEditor := new(MockTenderEditor)
	mockTender := models.Tender{
		TenderName:      "New name",
		Description:     "qwe",
		ServiceType:     "op",
		Status:          "open",
		OrganizationId:  1,
		CreatorUsername: "sariya",
	}
	mockUpdateTender := models.Tender{
		TenderName: "New name",
	}
	h := tender.EditTender(ctx, logger, mockTenderEditor)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PATCH("/api/tender/:tender_id/edit", h)

	mockTenderEditor.On("EditTender", ctx, 1, mockUpdateTender).Return(mockTender, nil)

	reqData := tender.EditTenderRequest{
		TenderName: "New name",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqData)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, "/api/tender/1/edit", &buf)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var responseBody tender.EditTenderResponse
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.Equal(t, mockTender, responseBody.UpdatedTender)
	require.Equal(t, "ok", responseBody.Message)
}

// TestCannotEditTender_IdNotInt проверяет, что если
// tender_id в URL не число, то возвращается код 400 и сообщение.
func TestCannotEditTender_IdNotInt(t *testing.T) {
	ctx := context.Background()
	logger := slogdiscard.NewDiscardLogger()
	mockTenderEditor := new(MockTenderEditor)

	mockUpdateTender := models.Tender{
		TenderName: "New name",
	}
	mockTenderEditor.On("EditTender", ctx, 1, mockUpdateTender).Return(nil, nil)

	h := tender.EditTender(ctx, logger, mockTenderEditor)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PATCH("/api/tender/:tender_id/edit", h)

	req := httptest.NewRequest(http.MethodPatch, "/api/tender/abobus/edit", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	var responseBody tender.EditTenderResponse
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	require.Empty(t, responseBody.UpdatedTender)
	require.Equal(t, "cannot parse tender id", responseBody.Message)
}

func TestCannotEditTender_WrongFields(t *testing.T) {
	ctx := context.Background()
	logger := slogdiscard.NewDiscardLogger()
	mockTenderEditor := new(MockTenderEditor)

	reqData := struct {
		Aboba  string
		Status string
	}{
		Aboba:  "qwe",
		Status: "qwe",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqData)
	require.NoError(t, err)

	h := tender.EditTender(ctx, logger, mockTenderEditor)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PATCH("/api/tender/:tender_id/edit", h)

	req := httptest.NewRequest(http.MethodPatch, "/api/tender/1/edit", &buf)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	// var responseBody tender.EditTenderResponse
	// err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	// require.NoError(t, err)

	// require.Empty(t, responseBody.UpdatedTender)
	// require.Equal(t, "cannot parse tender id", responseBody.Message)
}
