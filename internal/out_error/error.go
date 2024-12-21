package outerror

import "errors"

var (
	ErrEmployeeNotFound                           = errors.New("employee not found")
	ErrTendersWithThisServiceTypeNotFound         = errors.New("not found tenders with this service type")
	ErrEmployeeTendersNotFound                    = errors.New("not found tenders for this employee")
	ErrOrganizationNotFound                       = errors.New("organization not found")
	ErrEmployeeNotResponsibleForOrganization      = errors.New("employee not responsible for organization")
	ErrUpdatedEmployeeNotResponsibleForCurrentOrg = errors.New("updated employee not responsible for current organization")
	ErrCurrentEmployeeNotResponsibleForUpdatedOrg = errors.New("current employee not responsible for updated organization")
	ErrUpdatedEmployeeNotResponsibleForUpdatedOrg = errors.New("updated employee not responsible for updated organization")
	ErrTenderNotFound                             = errors.New("tender not found")
	ErrTenderVersionNotFound                      = errors.New("tender version not found")
	ErrUnknownTenderStatus                        = errors.New("unknown tender status")
	ErrNothingToUpdate                            = errors.New("nothing to update")
	ErrNewTenderCannotCreatedWithStatusNotCreated = errors.New("tender cannot be created with status not created")
	ErrCannotSetThisTenderStatus                  = errors.New("cannot set tender status in this cases: PUBLISED -> CREATED, CLOSED -> CREATED")
	ErrEmployeeNotResponsibleForTender            = errors.New("employee not respobsible for this tender")
)
