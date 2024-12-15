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
	return tender.Status == TenderCreatedStatus
}

const (
	TenderCreatedStatus   = "CREATED"
	TenderPublishedStatus = "PUBLISHED"
	TenderClosedStatus    = "CLOSED"
)

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
		return *tender.Status == TenderCreatedStatus || *tender.Status == TenderPublishedStatus || *tender.Status == TenderClosedStatus
	}
	return true
}

// CanSetThisTenderStatus проверяет, может ли
// новый статус тендера быть установлен.
//
// - нельзя перевести тендер из статуса PUBLISED в CREATED;
//
// - нельзя перевести тендер из статуса CLOSED в CREATED;
//
// Если передать также же статус, то его установить можно.
func (tender *TenderToUpdate) CanSetThisTenderStatus(currTenderStatus string) bool {
	if newStatus := tender.Status; newStatus != nil {
		if currTenderStatus == TenderPublishedStatus && *newStatus == TenderCreatedStatus {
			return false
		} else if currTenderStatus == TenderClosedStatus && *newStatus == TenderCreatedStatus {
			return false
		}
	}
	return true
}
