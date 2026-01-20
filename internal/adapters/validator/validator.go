package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate validates a struct based on tags
func Validate(data interface{}) error {
	if err := validate.Struct(data); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatValidationErrors(validationErrors)
		}
		return err
	}
	return nil
}

// formatValidationErrors formats validation errors into a readable message
func formatValidationErrors(errs validator.ValidationErrors) error {
	var messages []string
	for _, err := range errs {
		messages = append(messages, formatFieldError(err))
	}
	return fmt.Errorf("%s", strings.Join(messages, "; "))
}

// formatFieldError formats a single field validation error
func formatFieldError(err validator.FieldError) string {
	field := strings.ToLower(err.Field())

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, err.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, err.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	default:
		return fmt.Sprintf("%s validation failed on %s", field, err.Tag())
	}
}
