package repository

import "errors"

var (
	ErrEmployeeNotFound             = errors.New("employee not found")
	ErrNoTendersWithThisServiceType = errors.New("not found tenders with this service type")
)
