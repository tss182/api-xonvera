package middleware

import (
	"encoding/json"
	"strings"

	"app/xonvera-core/internal/infrastructure/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// BodyLogger logs request and response bodies with sensitive field masking
func BodyLogger(env string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get request ID from context
		requestID := GetRequestID(c)
		// Log request
		logRequest(c, requestID)
		//debug mode only
		debug := false
		if env == "development" {
			debug = true
		}
		c.Locals("debug", debug)

		// Process request
		err := c.Next()

		// Log response
		logResponse(c, requestID)

		return err
	}
}

func logRequest(c *fiber.Ctx, requestID string) {
	method := c.Method()
	path := c.Path()
	query := c.Query("*")

	// Log basic request info
	logger.Info("Request received",
		zap.String("request_id", requestID),
		zap.String("method", method),
		zap.String("path", path),
	)

	// Log request body for POST/PUT/PATCH
	if method == "POST" || method == "PUT" || method == "PATCH" {
		body := c.Body()
		if len(body) > 0 {
			maskedBody := maskSensitiveFields(body)
			logger.Debug("Request body",
				zap.String("request_id", requestID),
				zap.String("body", string(maskedBody)),
			)
		}
	}

	// Log query params if present
	if query != "" {
		logger.Debug("Query parameters",
			zap.String("request_id", requestID),
			zap.String("query", query),
		)
	}
}

func logResponse(c *fiber.Ctx, requestID string) {
	status := c.Response().StatusCode()
	contentType := c.Get("Content-Type")

	// Log basic response info
	logger.Info("Response sent",
		zap.String("request_id", requestID),
		zap.Int("status_code", status),
		zap.String("content_type", contentType),
	)

	// Log response body for JSON responses
	if strings.Contains(contentType, "application/json") {
		body := c.Response().Body()
		if len(body) > 0 {
			maskedBody := maskSensitiveFields(body)
			logger.Debug("Response body",
				zap.String("request_id", requestID),
				zap.String("body", string(maskedBody)),
			)
		}
	}
}

// maskSensitiveFields masks sensitive fields in JSON data
func maskSensitiveFields(data []byte) []byte {
	var obj map[string]interface{}

	// Try to unmarshal as JSON
	if err := json.Unmarshal(data, &obj); err != nil {
		// Not valid JSON, return as-is
		return data
	}

	// Mask top-level sensitive fields
	maskFields(obj)

	// Marshal back to JSON
	maskedData, err := json.Marshal(obj)
	if err != nil {
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
				v[key] = "***MASKED***"
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
	sensitiveFields := map[string]bool{
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

	lowerField := strings.ToLower(fieldName)
	return sensitiveFields[lowerField]
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
