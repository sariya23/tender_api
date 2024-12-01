package service

import "errors"

var (
	ErrNoTendersFound      = errors.New("no tenders found")
	ErrEmployeeNotFound    = errors.New("employee not found")
	ErrUserTendersNotFound = errors.New("user tenders not found")
)
