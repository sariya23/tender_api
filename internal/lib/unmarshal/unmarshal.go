package unmarshal

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sariya23/tender/internal/api/tender"
)

func CreateRequest(body []byte, req tender.CreateTenderRequest) (tender.CreateTenderRequest, error) {
	err := json.Unmarshal(body, &req)

	if err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError

		if errors.As(err, &syntaxErr) {
			return tender.CreateTenderRequest{}, fmt.Errorf("%s: %w", err.Error(), ErrSyntax)
		} else if errors.As(err, &typeErr) {
			return tender.CreateTenderRequest{}, fmt.Errorf("%s: %w", err.Error(), ErrType)
		} else {
			return tender.CreateTenderRequest{}, fmt.Errorf("%s: %w", err.Error(), ErrUnknown)
		}
	}

	return req, nil
}
