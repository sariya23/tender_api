package models

type Tender struct {
	TenderName      string `json:"name" validate:"required"`
	Description     string `json:"description" validate:"required"`
	ServiceType     string `json:"service_type" validate:"required"`
	Status          string `json:"status" validate:"required"`
	OrganizationId  int    `json:"organization_id" validate:"required,gte=0"`
	CreatorUsername string `json:"creator_username" validate:"required"`
}

func (tender *Tender) IsNewTenderHasStatusCreated() bool {
	return tender.Status == "CREATED"
}

type TenderToUpdate struct {
	TenderName      *string `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	ServiceType     *string `json:"service_type,omitempty"`
	Status          *string `json:"status,omitempty"`
	OrganizationId  *int    `json:"organization_id,omitempty" validate:"omitempty,gte=0"`
	CreatorUsername *string `json:"creator_username,omitempty"`
}

func (tender *TenderToUpdate) IsTenderStatusKnown() bool {
	if tender.Status != nil {
		return *tender.Status == "CREATED" || *tender.Status == "PUBLISHED" || *tender.Status == "CLOSED"
	}

	return true
}
