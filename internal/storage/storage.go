package storage

import "errors"

var (
	ErrNoTenderPresence        = errors.New("no tenders found")
	ErrUserNotFound            = errors.New("user not found")
	ErrOrganizationNotFound    = errors.New("organization not found")
	ErrUserNotReponsibleForOrg = errors.New("this user not reponsible for organization")
	ErrNoTenderForThisUser     = errors.New("no tenders for this user")
)
