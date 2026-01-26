package domain

import "errors"

var (
	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrPhoneAlreadyExists = errors.New("phone number already registered")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenNotFound      = errors.New("token not found")

	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	// Validation errors
	ErrInvalidInput         = errors.New("invalid input")
	ErrPasswordTooShort     = errors.New("password must be at least 6 characters")
	ErrRequiredFieldMissing = errors.New("required field is missing")

	// Internal errors
	ErrInternalServer    = errors.New("internal server error")
	ErrDatabaseOperation = errors.New("database operation failed")
)

// AppError represents an application error with additional context
type AppError struct {
	Err        error
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(err error, message string, statusCode int) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		StatusCode: statusCode,
	}
}
