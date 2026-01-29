package domain

import "testing"

func TestAppError_ErrorMessage(t *testing.T) {
	err := NewAppError(ErrInvalidCredentials, "custom message", 400)
	if err.Error() != "custom message" {
		t.Fatalf("expected custom message, got %s", err.Error())
	}
}

func TestAppError_ErrorFallback(t *testing.T) {
	err := NewAppError(ErrInvalidCredentials, "", 401)
	if err.Error() != ErrInvalidCredentials.Error() {
		t.Fatalf("expected fallback error message, got %s", err.Error())
	}
}

func TestAppError_Unwrap(t *testing.T) {
	wrapped := NewAppError(ErrInvalidCredentials, "", 401)
	if wrapped.Unwrap() != ErrInvalidCredentials {
		t.Fatalf("expected unwrap to return original error")
	}
}
