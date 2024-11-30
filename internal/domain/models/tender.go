package models

type Tender struct {
	TenderName      string `json:"name" validate:"required"`
	Description     string `json:"description" validate:"required"`
	ServiceType     string `json:"service_type" validate:"required"`
	Status          string `json:"status" validate:"required"`
	OrganizationId  int    `json:"organization_id" validate:"required"`
	CreatorUsername string `json:"creator_username" validate:"required"`
}

type TenderToUpdate struct {
	TenderName      *string `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	ServiceType     *string `json:"serviceType,omitempty"`
	Status          *string `json:"status,omitempty"`
	OrganizationId  *int    `json:"organizationId,omitempty"`
	CreatorUsername *string `json:"creatorUsername,omitempty"`
}
