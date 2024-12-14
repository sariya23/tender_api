package unmarshal

import (
	"encoding/json"
	"errors"
	"fmt"

	schema "github.com/sariya23/tender/internal/api"
)

func CreateRequest(body []byte) (schema.CreateTenderRequest, error) {
	var req schema.CreateTenderRequest
	err := json.Unmarshal(body, &req)

	if err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError

		if errors.As(err, &syntaxErr) {
			return schema.CreateTenderRequest{}, fmt.Errorf("%s: %w", err.Error(), ErrSyntax)
		} else if errors.As(err, &typeErr) {
			return schema.CreateTenderRequest{}, fmt.Errorf("%s: %w", err.Error(), ErrType)
		} else {
			return schema.CreateTenderRequest{}, fmt.Errorf("%s: %w", err.Error(), ErrUnknown)
		}
	}

	return req, nil
}

func EditRequest(body []byte) (schema.EditTenderRequest, error) {
	var req schema.EditTenderRequest
	err := json.Unmarshal(body, &req)

	if err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError

		if errors.As(err, &syntaxErr) {
			return schema.EditTenderRequest{}, fmt.Errorf("%s: %w", err.Error(), ErrSyntax)
		} else if errors.As(err, &typeErr) {
			return schema.EditTenderRequest{}, fmt.Errorf("%s: %w", err.Error(), ErrType)
		} else {
			return schema.EditTenderRequest{}, fmt.Errorf("%s: %w", err.Error(), ErrUnknown)
		}
	}

	return req, nil
}
