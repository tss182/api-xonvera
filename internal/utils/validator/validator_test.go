package validator

import (
	"strings"
	"testing"
)

type EmbeddedStruct struct {
	Code string `json:"code" validate:"required"`
}

type jsonTagStruct struct {
	EmbeddedStruct
	Search string `query:"search" validate:"required"`
	Form   string `form:"form_value" validate:"required"`
	Name   string `json:"name" validate:"required"`
}

func TestGetJSONFieldName(t *testing.T) {
	payload := jsonTagStruct{}
	fields := getJSONFieldName(payload)

	assertFieldMap(t, fields, "Search", "search")
	assertFieldMap(t, fields, "Form", "form_value")
	assertFieldMap(t, fields, "Name", "name")
	assertFieldMap(t, fields, "Code", "code")
}

type validationPayload struct {
	Name     string `json:"name" validate:"required,min=3,max=5"`
	Title    string `json:"title" validate:"max=3"`
	Code     string `json:"code" validate:"stringNumberOnly"`
	Optional string `json:"optional" validate:"stringNumberRequired"`
	Choice   string `json:"choice" validate:"oneof=small medium"`
	Email    string `json:"email" validate:"email"`
	Missing  string `json:"missing" validate:"required"`
}

func TestValidationErrors(t *testing.T) {
	payload := validationPayload{
		Name:     "ab",
		Title:    "abcd",
		Code:     "abc-123",
		Optional: "abc-123",
		Choice:   "large",
		Email:    "invalid",
	}

	errs := Validation(payload)
	if len(errs) == 0 {
		t.Fatalf("expected validation errors")
	}

	assertContains(t, errs, "minimum field 'name' is 3")
	assertContains(t, errs, "maximum field 'title' is 3")
	assertContains(t, errs, "'code' just support letter and number")
	assertContains(t, errs, "field 'choice' value must be one of character 'small or medium'")
	assertContains(t, errs, "'email' invalid email format")
	assertContains(t, errs, "'missing' is required")
}

func TestValidationSkips(t *testing.T) {
	payload := validationPayload{
		Name:     "ab",
		Title:    "abcd",
		Code:     "abc-123",
		Optional: "abc-123",
		Choice:   "large",
		Email:    "invalid",
	}

	errs := Validation(payload, "code")
	for _, err := range errs {
		if strings.Contains(err, "code") {
			t.Fatalf("expected code errors to be skipped, got %v", errs)
		}
	}
}

func TestValidationEmbeddedStruct(t *testing.T) {
	type outerStruct struct {
		EmbeddedStruct
		Name string `json:"name" validate:"required"`
	}

	payload := outerStruct{}
	errs := Validation(payload)
	if len(errs) == 0 {
		t.Fatalf("expected validation errors")
	}
	assertContains(t, errs, "'code' is required")
	assertContains(t, errs, "'name' is required")
}

func assertContains(t *testing.T, haystack []string, needle string) {
	t.Helper()
	for _, item := range haystack {
		if item == needle {
			return
		}
	}
	t.Fatalf("expected to find %q in %v", needle, haystack)
}

func assertFieldMap(t *testing.T, fields map[string]string, key, expected string) {
	t.Helper()
	if fields[key] != expected {
		t.Fatalf("expected %s to map to %s, got %s", key, expected, fields[key])
	}
}
