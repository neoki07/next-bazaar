package validation

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

// withoutSpace checks if the field value does not contain any spaces
func withoutSpace(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()

	for _, char := range fieldValue {
		if unicode.IsSpace(char) {
			return false
		}
	}

	return true
}

// withoutPunct checks if the field value does not contain any punctuation marks
func withoutPunct(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()

	for _, char := range fieldValue {
		if unicode.IsPunct(char) {
			return false
		}
	}

	return true
}

// withoutSymbol checks if the field value does not contain any symbol characters
func withoutSymbol(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()

	for _, char := range fieldValue {
		if unicode.IsSymbol(char) {
			return false
		}
	}

	return true
}
