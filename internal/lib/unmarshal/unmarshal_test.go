package unmarshal_test

import (
	"testing"

	"github.com/sariya23/tender/internal/api"
	"github.com/sariya23/tender/internal/domain/models"
	"github.com/sariya23/tender/internal/lib/unmarshal"
	"github.com/stretchr/testify/require"
)

// TestCreateRequest_Success
func TestCreateRequest_SuccessAllFields(t *testing.T) {
	// Arrange
	reqBody := `{
		"tender": {
			"name": "qwe",
			"description":"qwe",
			"service_type": "qwe",
			"status": "qwe",
			"organization_id": 1,
			"creator_username": "qwe"			
		}
	}`
	expectedReq := api.CreateTenderRequest{
		Tender: models.Tender{
			TenderName:      "qwe",
			Description:     "qwe",
			ServiceType:     "qwe",
			Status:          "qwe",
			OrganizationId:  1,
			CreatorUsername: "qwe",
		},
	}

	// Act
	req, err := unmarshal.CreateRequest([]byte(reqBody))

	// Assert
	require.NoError(t, err)
	require.Equal(t, req, expectedReq)
}

// TestCreateRequest_SuccessPartFields
func TestCreateRequest_SuccessPartFields(t *testing.T) {
	// Arrange
	reqBody := `{
		"tender": {
			"description":"qwe",
			"service_type": "qwe",
			"status": "qwe",
			"organization_id": 1,
			"creator_username": "qwe"			
		}
	}`
	expectedReq := api.CreateTenderRequest{
		Tender: models.Tender{
			TenderName:      "",
			Description:     "qwe",
			ServiceType:     "qwe",
			Status:          "qwe",
			OrganizationId:  1,
			CreatorUsername: "qwe",
		},
	}

	// Act
	req, err := unmarshal.CreateRequest([]byte(reqBody))

	// Assert
	require.NoError(t, err)
	require.Equal(t, req, expectedReq)
}

func TestCreateRequest_FailSyntaxError(t *testing.T) {
	reqBodyMissingComma := `{
		"tender": {
			"description":"qwe",
			"service_type": "qwe"
			"status": "qwe",
			"organization_id": 1,
			"creator_username": "qwe"			
		}
	}`
	reqBodyNoQuotes := `{
		"tender": {
			"description":qwe,
			"service_type": "qwe",
			"status": "qwe",
			"organization_id": 1,
			"creator_username": "qwe"			
		}
	}`
	reqBodyMissingBracket := `{
		"tender": {
			"description":"qwe",
			"service_type": "qwe",
			"status": "qwe",
			"organization_id": 1,
			"creator_username": "qwe"			
		}
	`
	reqBodyMissingDot := `{
		"tender": {
			"description" "qwe",
			"service_type": "qwe"
			"status": "qwe",
			"organization_id": 1,
			"creator_username": "qwe"			
		}
	}`
	cases := []struct {
		name string
		body string
	}{
		{name: "missing comma", body: reqBodyMissingComma},
		{name: "missing bracker", body: reqBodyMissingBracket},
		{name: "missing :", body: reqBodyMissingDot},
		{name: "no quotes", body: reqBodyNoQuotes},
	}
	for _, ts := range cases {
		t.Run(ts.name, func(t *testing.T) {
			expectedReq := api.CreateTenderRequest{}

			req, err := unmarshal.CreateRequest([]byte(ts.body))

			require.ErrorIs(t, err, unmarshal.ErrSyntax)
			require.Equal(t, req, expectedReq)
		})
	}
}

// TestCreateRequest_FailTypeErr
func TestCreateRequest_FailTypeErr(t *testing.T) {
	reqBodyWrongString := `{
		"tender": {
			"description":123,
			"service_type": "qwe",
			"status": "qwe",
			"organization_id": 1,
			"creator_username": "qwe"			
		}
	}`
	reqBodyWrongInt := `{
		"tender": {
			"description":"qwe",
			"service_type": "qwe",
			"status": "qwe",
			"organization_id": "asd",
			"creator_username": "qwe"			
}}
	`
	reqBodyBoolInteadOfInt := `{
		"tender": {
			"name": "qwe",
			"description":"qwe",
			"service_type": "qwe",
			"status": "qwe",
			"organization_id": false,
			"creator_username": "qwe"			
		}
	}`
	cases := []struct {
		name string
		body string
	}{
		{name: "int instead of string", body: reqBodyWrongString},
		{name: "string instead of int", body: reqBodyWrongInt},
		{name: "bool instead of int :", body: reqBodyBoolInteadOfInt},
	}
	for _, ts := range cases {
		t.Run(ts.name, func(t *testing.T) {
			expectedReq := api.CreateTenderRequest{}

			req, err := unmarshal.CreateRequest([]byte(ts.body))

			require.ErrorIs(t, err, unmarshal.ErrType)
			require.Equal(t, req, expectedReq)
		})
	}
}
