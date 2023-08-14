package validation

import (
	"github.com/go-playground/validator/v10"
)

// registerValidations registers multiple custom validations to the given validator.Validate instance.
func registerValidations(v *validator.Validate) {
	v.RegisterValidation("without_space", withoutSpace)
	v.RegisterValidation("without_punct", withoutPunct)
	v.RegisterValidation("without_symbol", withoutSymbol)

	v.RegisterValidation("decimal", isDecimal)
	v.RegisterValidation("decimal_gt", isDecimalGt)
}

// NewValidator func for create a new validator for api requests.
func NewValidator() *validator.Validate {
	v := validator.New()
	registerValidations(v)
	return v
}
