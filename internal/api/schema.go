package schema

import "github.com/sariya23/tender/internal/domain/models"

type GetTendersResponse struct {
	Tenders []models.Tender `json:"tenders"`
	Message string          `json:"message"`
}

type CreateTenderRequest struct {
	Tender models.Tender `json:"tender"`
}

type CreateTenderResponse struct {
	Tender  models.Tender `json:"tender,omitempty"`
	Message string        `json:"message"`
}

type GetEmployeeTendersResponse struct {
	Tenders []models.Tender `json:"tenders"`
	Message string          `json:"message"`
}

type EditTenderRequest struct {
	UpdateTenderData models.TenderToUpdate `json:"update_tender_data"`
}

type EditTenderResponse struct {
	UpdatedTender models.Tender `json:"updated_tender,omitempty"`
	Message       string        `json:"message"`
}

type RollbackTenderResponse struct {
	RollbackTender models.Tender `json:"rollback_tender,omitempty"`
	Message        string        `json:"message"`
}
