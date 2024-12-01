package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Tender struct {
	TenderName      string `json:"name" validate:"required"`
	Description     string `json:"description" validate:"required"`
	ServiceType     string `json:"service_type" validate:"required"`
	Status          string `json:"status" validate:"required"`
	OrganizationId  int    `json:"organization_id" validate:"required,gte=0"`
	CreatorUsername string `json:"creator_username" validate:"required"`
}

func main() {
	var t Tender
	b := `{
		"asdasd":"qwe",
		"zxcwaes":"qwe",
		"qwe":"qwe"
	}`
	err := json.Unmarshal([]byte(b), &t)
	var syErr *json.SyntaxError
	var qwe *json.UnmarshalTypeError
	fmt.Println(errors.As(err, &qwe))
	fmt.Println(errors.As(err, &syErr))
	fmt.Println(err)
	fmt.Println(t)
	validate := validator.New(validator.WithRequiredStructEnabled())
	fmt.Println(validate.Struct(t))
}
