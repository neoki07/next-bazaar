package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/ot07/next-bazaar/api/validations"
)

// registerValidations registers multiple custom validations to the given validator.Validate instance.
func registerValidations(v *validator.Validate) {
	v.RegisterValidation("without_space", validations.WithoutSpace)
	v.RegisterValidation("without_punct", validations.WithoutPunct)
	v.RegisterValidation("without_symbol", validations.WithoutSymbol)
}

// newValidator func for create a new validator for api requests.
func newValidator() *validator.Validate {
	v := validator.New()
	registerValidations(v)
	return v
}
