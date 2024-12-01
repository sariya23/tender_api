package unmarshal

import "errors"

var (
	ErrSyntax  = errors.New("syntax error")
	ErrType    = errors.New("wrong types")
	ErrUnknown = errors.New("unknown error")
)
