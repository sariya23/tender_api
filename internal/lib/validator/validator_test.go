package validator_test

import (
	"tender/internal/lib/validator"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestValidateEditRequest_Success проверяет, что при
// передаче json с известыми полями, либо с их частью,
// функция не возвращает ошибок.
func TestValidateEditRequest_Success(t *testing.T) {
	type EditTenderRequest struct {
		TenderName      string `json:"name"`
		Description     string `json:"description"`
		ServiceType     string `json:"serviceType"`
		Status          string `json:"status"`
		OrganizationId  int    `json:"organizationId"`
		CreatorUsername string `json:"creatorUsername"`
	}
	cases := []struct {
		Name string
		JSON string
	}{
		{Name: "All fields", JSON: `{"TenderName": "qwe", "Description": "qwe", "ServiceType": "qwe", "Status": "qwe", "OrganizationId": 1, "CreatorUsername": "qwe"}`},
		{Name: "part of fields", JSON: `{"TenderName": "qwe", "Description": "qwe", "ServiceType": "qwe"}`},
		{Name: "one field", JSON: `{"TenderName": "qwe"}`},
	}

	for _, ts := range cases {
		t.Run(ts.Name, func(t *testing.T) {
			var s EditTenderRequest
			err := validator.ValidateEditRequest([]byte(ts.JSON), &s)
			require.NoError(t, err)
		})
	}
}

// TestValidateEditRequest_FailUInknownField проверяет, что при
// передаче json неизвестыними полями, либо с их частью,
// функция вернет ошибку ErrUnknownField.
func TestValidateEditRequest_FailUInknownField(t *testing.T) {
	type EditTenderRequest struct {
		TenderName      string `json:"name"`
		Description     string `json:"description"`
		ServiceType     string `json:"serviceType"`
		Status          string `json:"status"`
		OrganizationId  int    `json:"organizationId"`
		CreatorUsername string `json:"creatorUsername"`
	}
	cases := []struct {
		Name string
		JSON string
	}{
		{Name: "All fields unknown", JSON: `{"ASD": "qwe", "QWE": "qwe", "ZXC": "qwe", "QWE": "qwe", "ZXCAS": 1, "ASD": "qwe"}`},
		{Name: "One field unknown", JSON: `{"TenderName": "qwe", "Description": "qwe", "ServiceType": "qwe", "ASDASD": "zxc"}`},
	}

	for _, ts := range cases {
		t.Run(ts.Name, func(t *testing.T) {
			var s EditTenderRequest
			err := validator.ValidateEditRequest([]byte(ts.JSON), &s)
			require.ErrorIs(t, err, validator.ErrUnknownField)
		})
	}
}
