package outerror

import "errors"

var (
	ErrEmployeeNotFound                      = errors.New("employee not found")
	ErrTendersWithThisServiceTypeNotFound    = errors.New("not found tenders with this service type")
	ErrEmployeeTendersNotFound               = errors.New("not found tenders for this employee")
	ErrOrganizationNotFound                  = errors.New("organization not found")
	ErrEmployeeNotResponsibleForOrganization = errors.New("employee not responsible for organization")
	ErrTenderNotFound                        = errors.New("tender not found")
	ErrTenderVersionNotFound                 = errors.New("tender version not found")
)
