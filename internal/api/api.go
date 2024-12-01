package api

import "github.com/sariya23/tender/internal/domain/models"

type GetTendersResponse struct {
	Tenders []models.Tender `json:"tenders,omitempty"`
	Message string          `json:"message"`
}

type CreateTenderRequest struct {
	Tender models.Tender `json:"tender"`
}

type CreateTenderResponse struct {
	Tender  models.Tender `json:"tender,omitempty"`
	Message string        `json:"message"`
}
