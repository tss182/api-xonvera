package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

type stringValidationPayload struct {
	Code     string `validate:"stringNumberOnly"`
	Optional string `validate:"stringNumberRequired"`
}

func TestStringValidations(t *testing.T) {
	v := validator.New()
	Init(v)

	valid := stringValidationPayload{Code: "abc123", Optional: ""}
	if err := v.Struct(valid); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	invalid := stringValidationPayload{Code: "abc-123", Optional: "abc-123"}
	if err := v.Struct(invalid); err == nil {
		t.Fatalf("expected validation error")
	}
}
