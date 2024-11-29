package validator

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrSyntaxError       = errors.New("syntax error")
	ErrUnknownField      = errors.New("unknown field")
	ErrTypeMismatch      = errors.New("type mismatch")
	ErrInvalidJSONFormat = errors.New("invalid JSON format")
	ErrInvalidInput      = errors.New("len of types and fields is different")
)

func ValidateEditRequest(data []byte, v any, fields []string, types []string) error {
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			return fmt.Errorf("%w: at byte %d", ErrSyntaxError, syntaxErr.Offset)
		}
		return fmt.Errorf("%w: %v", ErrInvalidJSONFormat, err)
	}

	validFields := makeValidField(fields)

	for key := range rawMap {
		if _, valid := validFields[key]; !valid {
			return fmt.Errorf("%w: %s", ErrUnknownField, key)
		}
	}

	if err := json.Unmarshal(data, &v); err != nil {
		if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
			return fmt.Errorf("%w: for field %s, expected %s but got %s",
				ErrTypeMismatch, typeErr.Field, typeErr.Type.String(), typeErr.Value)
		}
		return fmt.Errorf("%w: %v", ErrInvalidJSONFormat, err)
	}

	m := makeExpectedTypes(fields, types)
	if m == nil {
		return ErrInvalidInput
	}
	err := validateEditRequestFieldTypes(rawMap, m)
	return err
}

func validateEditRequestFieldTypes(rawMap map[string]json.RawMessage, expectedTypes map[string]string) error {
	for field, rawValue := range rawMap {
		expectedType, exists := expectedTypes[field]
		if !exists {
			continue
		}

		var asInterface interface{}
		if err := json.Unmarshal(rawValue, &asInterface); err != nil {
			return fmt.Errorf("%w: error decoding field %s", ErrInvalidJSONFormat, field)
		}

		if !matchesExpectedType(asInterface, expectedType) {
			return fmt.Errorf("%w: type mismatch for field %s (expected %s)",
				ErrTypeMismatch, field, expectedType)
		}
	}
	return nil
}

func matchesExpectedType(value interface{}, expectedType string) bool {
	switch expectedType {
	case "string":
		_, ok := value.(string)
		return ok
	case "int":
		_, ok := value.(float64)
		return ok
	default:
		return false
	}
}

func makeValidField(fields []string) map[string]struct{} {
	m := make(map[string]struct{}, len(fields))
	for _, v := range fields {
		m[v] = struct{}{}
	}
	return m
}

func makeExpectedTypes(fields []string, types []string) map[string]string {
	if len(fields) != len(types) {
		return nil
	}
	m := make(map[string]string, len(fields))
	for i := 0; i < len(fields); i++ {
		m[fields[i]] = types[i]
	}
	return m
}
