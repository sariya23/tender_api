package models

type Tender struct {
	TenderName      string `json:"name"`
	Description     string `json:"creatorUsername"`
	ServiceType     string `json:"serviceType"`
	Status          string `json:"status"`
	OrganizationId  int    `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}
