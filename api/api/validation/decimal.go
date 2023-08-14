package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

func isDecimal(fl validator.FieldLevel) bool {
	field := fl.Field()
	_, err := decimal.NewFromString(field.String())
	return err == nil
}

func isDecimalGt(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	fieldDecimal, err := decimal.NewFromString(field.String())
	if err != nil {
		return false
	}

	paramDecimal, err := decimal.NewFromString(param)
	if err != nil {
		return false
	}

	return fieldDecimal.GreaterThan(paramDecimal)
}
