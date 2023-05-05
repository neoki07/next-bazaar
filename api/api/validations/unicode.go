package validations

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

// WithoutSpace checks if the field value does not contain any spaces
func WithoutSpace(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()

	for _, char := range fieldValue {
		if unicode.IsSpace(char) {
			return false
		}
	}

	return true
}

// WithoutPunct checks if the field value does not contain any punctuation marks
func WithoutPunct(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()

	for _, char := range fieldValue {
		if unicode.IsPunct(char) {
			return false
		}
	}

	return true
}

// WithoutSymbol checks if the field value does not contain any symbol characters
func WithoutSymbol(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()

	for _, char := range fieldValue {
		if unicode.IsSymbol(char) {
			return false
		}
	}

	return true
}
