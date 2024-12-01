package service

import "errors"

var (
	ErrTendersNotFound                       = errors.New("no tenders found")
	ErrEmployeeNotFound                      = errors.New("employee not found")
	ErrEmployeeTendersNotFound               = errors.New("employee tenders not found")
	ErrOrganizationNotFound                  = errors.New("organization not found")
	ErrEmployeeNotResponsibleForOrganization = errors.New("employee not responsible for organization")
)
