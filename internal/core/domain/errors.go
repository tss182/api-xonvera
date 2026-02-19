package domain

// HTTP Error Response Codes
const (
	// 400 Bad Request Errors
	ErrPhoneAlreadyRegistered = "400:phone number already registered"
	ErrInvalidRegisterRequest = "400:invalid register request"
	ErrInvalidCredentials     = "400:invalid credentials"
	ErrInvalidRefreshToken    = "400:invalid refresh token"
	ErrRefreshTokenExpired    = "400:refresh token has expired"
	ErrInvoiceIDRequired      = "400:invoice ID is required for update"
	ErrInvalidPackage         = "400:invalid package"

	// 404 Not Found Errors
	ErrNotFoundInvoice = "404:not found invoice"
	ErrNotFoundToken   = "404:not found token"

	// 401 Unauthorized Errors
	ErrUnauthorized = "401:unauthorized"
)
