package middleware

import (
	"encoding/json"
	"strings"

	"app/xonvera-core/internal/infrastructure/logger"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// Sensitive field names to mask in logs
var sensitiveFields = map[string]bool{
	"password":        true,
	"pwd":             true,
	"pass":            true,
	"secret":          true,
	"token":           true,
	"access_token":    true,
	"refresh_token":   true,
	"api_key":         true,
	"apikey":          true,
	"authorization":   true,
	"auth":            true,
	"bearer":          true,
	"ssn":             true,
	"social_security": true,
	"credit_card":     true,
	"creditcard":      true,
	"cvv":             true,
	"cvc":             true,
	"pin":             true,
}

const (
	maskedValue     = "***MASKED***"
	maxBodyLogSize  = 2048 // Log max 2KB of body
	jsonContentType = "application/json"
)

// BodyLogger logs request and response bodies with sensitive field masking
func BodyLogger(env string, enabled bool) fiber.Handler {
	if !enabled {
		return func(c fiber.Ctx) error {
			return c.Next()
		}
	}

	return func(c fiber.Ctx) error {

		// Log request
		logRequest(c)

		// Set debug mode based on environment
		debug := env == "development"
		c.Locals("debug", debug)

		// Process request
		err := c.Next()

		// Log response
		logResponse(c)

		return err
	}
}

// logRequest logs incoming request details
func logRequest(c fiber.Ctx) {
	method := c.Method()
	path := c.Path()
	query := c.Query("*")

	// Log basic request info
	logger.ContextInfo(c, "Request received",
		zap.String("method", method),
		zap.String("path", path),
	)

	// Log request body for mutation methods
	if shouldLogBody(method) {
		body := c.Body()
		if len(body) > 0 {
			logBodyData(c, body, "request")
		}
	}

	// Log query parameters if present
	if query != "" {
		logger.ContextDebug(c, "Query parameters",
			zap.String("query", query),
		)
	}
}

// logResponse logs outgoing response details
func logResponse(c fiber.Ctx) {
	status := c.Response().StatusCode()
	contentType := c.Get("Content-Type")

	// Log basic response info
	logger.ContextInfo(c, "Response sent",
		zap.Int("status_code", status),
		zap.String("content_type", contentType),
	)

	// Log response body for JSON responses
	if strings.Contains(contentType, jsonContentType) {
		body := c.Response().Body()
		if len(body) > 0 {
			logBodyData(c, body, "response")
		}
	}
}

// logBodyData logs body data with masking
func logBodyData(c fiber.Ctx, body []byte, bodyType string) {
	if len(body) > maxBodyLogSize {
		body = body[:maxBodyLogSize]
	}

	maskedBody := maskSensitiveFields(body)
	logger.ContextDebug(c, "Body data",
		zap.String("type", bodyType),
		zap.String("data", string(maskedBody)),
	)
}

// shouldLogBody determines if request body should be logged
func shouldLogBody(method string) bool {
	return method == fiber.MethodPost ||
		method == fiber.MethodPut ||
		method == fiber.MethodPatch
}

// maskSensitiveFields masks sensitive fields in JSON data
func maskSensitiveFields(data []byte) []byte {
	var obj map[string]interface{}

	// Try to unmarshal as JSON
	if err := json.Unmarshal(data, &obj); err != nil {
		// Not valid JSON, return as-is (truncated if needed)
		return data
	}

	// Mask sensitive fields
	maskFields(obj)

	// Marshal back to JSON
	maskedData, err := json.Marshal(obj)
	if err != nil {
		// If marshaling fails, return original
		return data
	}

	return maskedData
}

// maskFields recursively masks sensitive fields in a map
func maskFields(obj interface{}) {
	switch v := obj.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if isSensitiveField(key) {
				v[key] = maskedValue
			} else if isNestedObject(value) {
				maskFields(value)
			}
		}
	case []interface{}:
		for i, item := range v {
			if isNestedObject(item) {
				maskFields(item)
				v[i] = item
			}
		}
	}
}

// isSensitiveField checks if a field should be masked
func isSensitiveField(fieldName string) bool {
	return sensitiveFields[strings.ToLower(fieldName)]
}

// isNestedObject checks if a value is a nested object or array
func isNestedObject(value interface{}) bool {
	switch value.(type) {
	case map[string]interface{}, []interface{}:
		return true
	default:
		return false
	}
}
