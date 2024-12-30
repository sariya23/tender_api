package testdata

import "github.com/sariya23/tender/internal/domain/models"

var TestOrganization = models.Organization{
	ID:          1,
	Name:        "Test Organization",
	Description: "Test Organization",
	Type:        "OSI",
}

var TestEmployee = models.Employee{
	ID:        1,
	Username:  "sariya",
	FirstName: "Test",
	LastName:  "Testovisch",
}

var TestTender = models.Tender{
	TenderName:      "Test Tender",
	Description:     "Test Tender",
	ServiceType:     "testing",
	Status:          "OPEN",
	OrganizationId:  TestOrganization.ID,
	CreatorUsername: TestEmployee.Username,
}
