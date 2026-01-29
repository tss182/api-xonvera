package utils

import (
	"testing"
	"time"
)

func TestTernary(t *testing.T) {
	if got := Ternary(true, "yes", "no"); got != "yes" {
		t.Fatalf("expected yes, got %s", got)
	}
	if got := Ternary(false, 1, 2); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestGetTimeDuration(t *testing.T) {
	if got := GetTimeDuration("minute", 2); got != 2*time.Minute {
		t.Fatalf("expected 2 minutes, got %s", got)
	}
	if got := GetTimeDuration("hour", 3); got != 3*time.Hour {
		t.Fatalf("expected 3 hours, got %s", got)
	}
	if got := GetTimeDuration("day", 1); got != 24*time.Hour {
		t.Fatalf("expected 24 hours, got %s", got)
	}
	if got := GetTimeDuration("unknown", 5); got != 0 {
		t.Fatalf("expected 0 duration, got %s", got)
	}
}

func TestHashSha256(t *testing.T) {
	plaintext := "hello"
	hash := HashSha256(plaintext)
	if len(hash) == 0 {
		t.Fatalf("expected hash to be non-empty")
	}
	if !VerifyHashSha256(plaintext, hash) {
		t.Fatalf("expected hash to verify for plaintext")
	}
	if VerifyHashSha256("other", hash) {
		t.Fatalf("expected hash to fail for different plaintext")
	}
}

func TestValidateEmail(t *testing.T) {
	if !ValidateEmail("user@example.com") {
		t.Fatalf("expected valid email")
	}
	if ValidateEmail("not-an-email") {
		t.Fatalf("expected invalid email")
	}
}

func TestPhoneNumberFormat(t *testing.T) {
	if got := PhoneNumberFormat("0812-345 678"); got != "62812345678" {
		t.Fatalf("expected 62812345678, got %s", got)
	}
	if got := PhoneNumberFormat("8123"); got != "628123" {
		t.Fatalf("expected 628123, got %s", got)
	}
	if got := PhoneNumberFormat("628123"); got != "628123" {
		t.Fatalf("expected 628123, got %s", got)
	}
}

func TestInArray(t *testing.T) {
	values := []string{"a", "b", "c"}
	if !InArray(values, "b") {
		t.Fatalf("expected to find b in array")
	}
	if InArray(values, "d") {
		t.Fatalf("expected to not find d in array")
	}
}

func TestDecimalSeparator(t *testing.T) {
	if got := DecimalSeparator(1000); got != "1,000" {
		t.Fatalf("expected 1,000, got %s", got)
	}
	if got := DecimalSeparator(10); got != "10" {
		t.Fatalf("expected 10, got %s", got)
	}
}
