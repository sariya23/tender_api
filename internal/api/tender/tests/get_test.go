package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sariya23/tender/internal/api/tender"
	"github.com/sariya23/tender/internal/api/tender/mocks"
	"github.com/sariya23/tender/internal/domain/models"
	"github.com/sariya23/tender/internal/lib/logger/slogdiscard"
	"github.com/stretchr/testify/assert"
)

// TestGetAllTenders_Success проверяет
// успешный сценарий вызова хендлера GetTenders.
//
// Возвращается код 200 и тело со списков тендеров.
func TestGetAllTenders_Success(t *testing.T) {
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
	svc := tender.New(logger, mockTenderService)

	mockTenderService.On("GetTenders", ctx, "all").Return(mockTenders, nil)
	req := httptest.NewRequest(http.MethodGet, "/tenders?srv_type=all", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := svc.GetTenders(ctx)
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expectedBody, w.Body.String())

}
