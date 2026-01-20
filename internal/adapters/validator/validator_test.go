package validator

import (
	"testing"
)

// TestStructures for validation tests
type ValidUser struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Phone    string `validate:"required"`
	Password string `validate:"required,min=6"`
}

type UserWithMaxField struct {
	Name string `validate:"max=10"`
}

type UserWithLenField struct {
	Phone string `validate:"len=10"`
}

type UserWithGteField struct {
	Age int `validate:"gte=18"`
}

type UserWithLteField struct {
	Age int `validate:"lte=100"`
}

type UserWithUrlField struct {
	Website string `validate:"url"`
}

type UserWithMultipleErrors struct {
	Name     string `validate:"required,min=3"`
	Email    string `validate:"required,email"`
	Phone    string `validate:"required,len=10"`
	Password string `validate:"required,min=8"`
}

type NonStructData struct {
	Value string
}

// TestValidate_SuccessfulValidation tests successful validation
func TestValidate_SuccessfulValidation(t *testing.T) {
	user := ValidUser{
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "1234567890",
		Password: "password123",
	}

	err := Validate(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestValidate_RequiredFieldMissing tests required field validation
func TestValidate_RequiredFieldMissing(t *testing.T) {
	user := ValidUser{
		Name:     "",
		Email:    "john@example.com",
		Phone:    "1234567890",
		Password: "password123",
	}

	err := Validate(user)
	if err == nil {
		t.Error("Expected validation error for missing name")
	}
	if err != nil && err.Error() != "name is required" {
		t.Errorf("Expected 'name is required', got '%v'", err.Error())
	}
}

// TestValidate_InvalidEmail tests email validation
func TestValidate_InvalidEmail(t *testing.T) {
	user := ValidUser{
		Name:     "John Doe",
		Email:    "invalid-email",
		Phone:    "1234567890",
		Password: "password123",
	}

	err := Validate(user)
	if err == nil {
		t.Error("Expected validation error for invalid email")
	}
	if err != nil && err.Error() != "email must be a valid email address" {
		t.Errorf("Expected 'email must be a valid email address', got '%v'", err.Error())
	}
}

// TestValidate_MinLengthValidation tests min length validation
func TestValidate_MinLengthValidation(t *testing.T) {
	user := ValidUser{
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "1234567890",
		Password: "short",
	}

	err := Validate(user)
	if err == nil {
		t.Error("Expected validation error for short password")
	}
	if err != nil && err.Error() != "password must be at least 6 characters" {
		t.Errorf("Expected 'password must be at least 6 characters', got '%v'", err.Error())
	}
}

// TestValidate_MaxLengthValidation tests max length validation
func TestValidate_MaxLengthValidation(t *testing.T) {
	user := UserWithMaxField{
		Name: "This is a very long name that exceeds 10 characters",
	}

	err := Validate(user)
	if err == nil {
		t.Error("Expected validation error for max length")
	}
	if err != nil && err.Error() != "name must be at most 10 characters" {
		t.Errorf("Expected 'name must be at most 10 characters', got '%v'", err.Error())
	}
}

// TestValidate_LenValidation tests exact length validation
func TestValidate_LenValidation(t *testing.T) {
	user := UserWithLenField{
		Phone: "12345",
	}

	err := Validate(user)
	if err == nil {
		t.Error("Expected validation error for exact length")
	}
	if err != nil && err.Error() != "phone must be exactly 10 characters" {
		t.Errorf("Expected 'phone must be exactly 10 characters', got '%v'", err.Error())
	}
}

// TestValidate_GteValidation tests greater than or equal validation
func TestValidate_GteValidation(t *testing.T) {
	user := UserWithGteField{
		Age: 15,
	}

	err := Validate(user)
	if err == nil {
		t.Error("Expected validation error for gte")
	}
	if err != nil && err.Error() != "age must be greater than or equal to 18" {
		t.Errorf("Expected 'age must be greater than or equal to 18', got '%v'", err.Error())
	}
}

// TestValidate_LteValidation tests less than or equal validation
func TestValidate_LteValidation(t *testing.T) {
	user := UserWithLteField{
		Age: 150,
	}

	err := Validate(user)
	if err == nil {
		t.Error("Expected validation error for lte")
	}
	if err != nil && err.Error() != "age must be less than or equal to 100" {
		t.Errorf("Expected 'age must be less than or equal to 100', got '%v'", err.Error())
	}
}

// TestValidate_AllValidationTags tests that all validation tag cases are covered
func TestValidate_AllValidationTags(t *testing.T) {
	tests := []struct {
		name          string
		data          interface{}
		expectError   bool
		expectedError string
	}{
		{
			name: "valid_register_request",
			data: ValidUser{
				Name:     "John Doe",
				Email:    "john@example.com",
				Phone:    "1234567890",
				Password: "password123",
			},
			expectError: false,
		},
		{
			name: "invalid_required_field",
			data: ValidUser{
				Name:     "",
				Email:    "john@example.com",
				Phone:    "1234567890",
				Password: "password123",
			},
			expectError:   true,
			expectedError: "name is required",
		},
		{
			name: "invalid_email_format",
			data: ValidUser{
				Name:     "John Doe",
				Email:    "not-an-email",
				Phone:    "1234567890",
				Password: "password123",
			},
			expectError:   true,
			expectedError: "email must be a valid email address",
		},
		{
			name: "invalid_min_length",
			data: ValidUser{
				Name:     "John Doe",
				Email:    "john@example.com",
				Phone:    "1234567890",
				Password: "123",
			},
			expectError:   true,
			expectedError: "password must be at least 6 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.data)
			if (err != nil) != tt.expectError {
				t.Errorf("expectError: got %v, want %v", err != nil, tt.expectError)
			}
			if tt.expectError && err.Error() != tt.expectedError {
				t.Errorf("expected error: got '%s', want '%s'", err.Error(), tt.expectedError)
			}
		})
	}
}

// TestValidate_UrlValidation tests URL validation
func TestValidate_UrlValidation(t *testing.T) {
	user := UserWithUrlField{
		Website: "not-a-url",
	}

	err := Validate(user)
	if err == nil {
		t.Error("Expected validation error for invalid URL")
	}
	if err != nil && err.Error() != "website must be a valid URL" {
		t.Errorf("Expected 'website must be a valid URL', got '%v'", err.Error())
	}
}

// TestValidate_ValidUrl tests valid URL
func TestValidate_ValidUrl(t *testing.T) {
	user := UserWithUrlField{
		Website: "https://example.com",
	}

	err := Validate(user)
	if err != nil {
		t.Errorf("Expected no error for valid URL, got %v", err)
	}
}

// TestValidate_MultipleErrors tests validation with multiple errors
func TestValidate_MultipleErrors(t *testing.T) {
	user := UserWithMultipleErrors{
		Name:     "Jo",
		Email:    "invalid",
		Phone:    "123",
		Password: "short",
	}

	err := Validate(user)
	if err == nil {
		t.Error("Expected validation errors")
	}
	if err != nil {
		errMsg := err.Error()
		// Check that multiple errors are present
		if !containsSubstring(errMsg, "name must be") {
			t.Errorf("Expected name error in '%v'", errMsg)
		}
		if !containsSubstring(errMsg, "email must be") {
			t.Errorf("Expected email error in '%v'", errMsg)
		}
	}
}

// TestValidate_PartialErrors tests validation with some valid and some invalid fields
func TestValidate_PartialErrors(t *testing.T) {
	user := UserWithMultipleErrors{
		Name:     "John Doe",
		Email:    "invalid-email",
		Phone:    "1234567890",
		Password: "password123",
	}

	err := Validate(user)
	if err == nil {
		t.Error("Expected validation error for invalid email")
	}
	if err != nil && err.Error() != "email must be a valid email address" {
		t.Errorf("Expected 'email must be a valid email address', got '%v'", err.Error())
	}
}

// TestValidate_ValidComplexStruct tests validation with all valid complex data
func TestValidate_ValidComplexStruct(t *testing.T) {
	user := UserWithMultipleErrors{
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "1234567890",
		Password: "password123456",
	}

	err := Validate(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestValidate_NilValue tests validation with nil
func TestValidate_NilValue(t *testing.T) {
	var user *ValidUser
	err := Validate(user)
	if err == nil {
		t.Error("Expected validation error for nil pointer")
	}
}

// TestValidate_EmptyStruct tests validation with empty struct
func TestValidate_EmptyStruct(t *testing.T) {
	user := ValidUser{}
	err := Validate(user)
	if err == nil {
		t.Error("Expected validation error for empty struct")
	}
}

// TestValidate_ValidMinValue tests valid min length value
func TestValidate_ValidMinValue(t *testing.T) {
	user := ValidUser{
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "1234567890",
		Password: "123456",
	}

	err := Validate(user)
	if err != nil {
		t.Errorf("Expected no error for valid min length, got %v", err)
	}
}

// TestValidate_ValidMaxValue tests valid max length value
func TestValidate_ValidMaxValue(t *testing.T) {
	user := UserWithMaxField{
		Name: "12345678",
	}

	err := Validate(user)
	if err != nil {
		t.Errorf("Expected no error for valid max length, got %v", err)
	}
}

// TestValidate_ValidLenValue tests valid exact length value
func TestValidate_ValidLenValue(t *testing.T) {
	user := UserWithLenField{
		Phone: "1234567890",
	}

	err := Validate(user)
	if err != nil {
		t.Errorf("Expected no error for valid exact length, got %v", err)
	}
}

// TestValidate_ValidGteValue tests valid gte value
func TestValidate_ValidGteValue(t *testing.T) {
	user := UserWithGteField{
		Age: 18,
	}

	err := Validate(user)
	if err != nil {
		t.Errorf("Expected no error for valid gte value, got %v", err)
	}
}

// TestValidate_ValidLteValue tests valid lte value
func TestValidate_ValidLteValue(t *testing.T) {
	user := UserWithLteField{
		Age: 100,
	}

	err := Validate(user)
	if err != nil {
		t.Errorf("Expected no error for valid lte value, got %v", err)
	}
}

// Helper function to check if error message contains a substring
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// Benchmark tests
func BenchmarkValidate_Success(b *testing.B) {
	user := ValidUser{
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "1234567890",
		Password: "password123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Validate(user)
	}
}

func BenchmarkValidate_WithErrors(b *testing.B) {
	user := ValidUser{
		Name:     "",
		Email:    "invalid",
		Phone:    "",
		Password: "short",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Validate(user)
	}
}
