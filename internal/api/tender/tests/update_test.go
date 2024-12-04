package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	tenderapi "github.com/sariya23/tender/internal/api/tender"
	"github.com/sariya23/tender/internal/api/tender/mocks"
	"github.com/sariya23/tender/internal/domain/models"
	"github.com/sariya23/tender/internal/lib/logger/slogdiscard"
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
	router.PATCH("/api/tenders/:tenderId/edit", svc.EditTedner(ctx))
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
