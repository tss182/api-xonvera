package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func StringNumberOnly(fl validator.FieldLevel) bool {
	plaintext := fl.Field().String()
	return regexp.MustCompile(`^[a-zA-Z0-9\s]+$`).MatchString(plaintext)
}

func StringNumberRequired(fl validator.FieldLevel) bool {
	plaintext := fl.Field().String()
	if plaintext == "" {
		return true
	}
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(plaintext)
}
