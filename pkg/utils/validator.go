package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Validate(i interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := v.validator.Struct(i)
	if err != nil {
		fmt.Println(err.Error())
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}