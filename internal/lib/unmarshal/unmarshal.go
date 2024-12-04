package unmarshal

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sariya23/tender/internal/api"
)

func CreateRequest(body []byte) (api.CreateTenderRequest, error) {
	var req api.CreateTenderRequest
	err := json.Unmarshal(body, &req)

	if err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError

		if errors.As(err, &syntaxErr) {
			return api.CreateTenderRequest{}, fmt.Errorf("%s: %w", err.Error(), ErrSyntax)
		} else if errors.As(err, &typeErr) {
			return api.CreateTenderRequest{}, fmt.Errorf("%s: %w", err.Error(), ErrType)
		} else {
			return api.CreateTenderRequest{}, fmt.Errorf("%s: %w", err.Error(), ErrUnknown)
		}
	}

	return req, nil
}

func EditRequest(body []byte) (api.EditTenderRequest, error) {
	var req api.EditTenderRequest
	err := json.Unmarshal(body, &req)

	if err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError

		if errors.As(err, &syntaxErr) {
			return api.EditTenderRequest{}, fmt.Errorf("%s: %w", err.Error(), ErrSyntax)
		} else if errors.As(err, &typeErr) {
			return api.EditTenderRequest{}, fmt.Errorf("%s: %w", err.Error(), ErrType)
		} else {
			return api.EditTenderRequest{}, fmt.Errorf("%s: %w", err.Error(), ErrUnknown)
		}
	}

	return req, nil
}
