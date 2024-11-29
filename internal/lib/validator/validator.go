package validator

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrUnknownField = errors.New("unknown filed")
)

func ValidateEditRequest(data []byte, v any) error {
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	validFields := map[string]struct{}{
		"TenderName":      {},
		"Description":     {},
		"ServiceType":     {},
		"Status":          {},
		"OrganizationId":  {},
		"CreatorUsername": {},
	}

	for key := range rawMap {
		if _, valid := validFields[key]; !valid {
			return fmt.Errorf("%w: %s", ErrUnknownField, key)
		}
	}

	return json.Unmarshal(data, v)
}
