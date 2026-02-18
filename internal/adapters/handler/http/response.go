package http

import (
	"app/xonvera-core/internal/core/domain"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type (
	Map map[string]interface{}

	Resp struct {
		RequestID string      `json:"request_id"`
		Errors    []string    `json:"errors,omitempty"`
		Meta      interface{} `json:"meta,omitempty"`
		Data      interface{} `json:"data"`
	}
)

// JSON sends a JSON response with the given status code
func JSON(c fiber.Ctx, httpStatus int, errs []string, data, meta interface{}) error {
	// Validate HTTP status code
	if httpStatus < 100 || httpStatus > 599 {
		httpStatus = http.StatusInternalServerError
	}

	// Extract request ID from context
	requestID := extractRequestID(c)

	// Default empty data to empty map
	if data == nil {
		data = Map{}
	}

	resp := Resp{
		RequestID: requestID,
		Errors:    errs,
		Meta:      meta,
		Data:      data,
	}

	return c.Status(httpStatus).JSON(resp)
}

// extractRequestID safely extracts request ID from context
func extractRequestID(c fiber.Ctx) string {
	if req, ok := c.Locals("request_id").(string); ok {
		return req
	}
	return ""
}

// OK sends a 200 OK response with data
func OK(c fiber.Ctx, resp interface{}) error {
	return JSON(c, http.StatusOK, nil, resp, nil)
}

// OK sends a 200 OK response with data
func Page(c fiber.Ctx, resp *domain.PaginationResponse) error {
	if resp.Data == nil {
		resp.Data = []any{}
	}
	return JSON(c, http.StatusOK, nil, resp.Data, resp.Meta)
}

// NoAuth sends a 401 Unauthorized response
func NoAuth(c fiber.Ctx) error {
	return JSON(c, http.StatusUnauthorized, []string{"unauthorized"}, nil, nil)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c fiber.Ctx, errs []string) error {
	return JSON(c, http.StatusBadRequest, errs, nil, nil)
}

// NotFound sends a 404 Not Found response
func NotFound(c fiber.Ctx, errs []string) error {
	return JSON(c, http.StatusNotFound, errs, nil, nil)
}

// Conflict sends a 409 Conflict response
func Conflict(c fiber.Ctx, errs []string) error {
	return JSON(c, http.StatusConflict, errs, nil, nil)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c fiber.Ctx, errs []string, debugError bool) error {
	if errs == nil || !debugError {
		errs = []string{"internal server error"}
	}
	return JSON(c, http.StatusInternalServerError, errs, nil, nil)
}

// ErrorLimited sends a 429 Too Many Requests response
func ErrorLimited(c fiber.Ctx, errs []string) error {
	return JSON(c, http.StatusTooManyRequests, errs, nil, nil)
}

// HandlerErrorGlobal is the global error handler for the application
func HandlerErrorGlobal(c fiber.Ctx, err error) error {
	// Handle record not found errors
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound(c, []string{"data not found"})
	}

	// Check if debug mode is enabled
	debug := isDebugMode(c)

	// Parse error string for HTTP code and message
	// Format: "httpCode:message:data"
	errs := strings.Split(err.Error(), ":")
	if len(errs) >= 2 {
		httpCode, err := strconv.Atoi(errs[0])
		if err == nil && httpCode >= 200 && httpCode <= 599 {
			var errResp []string
			if len(errs) > 1 && errs[1] != "" {
				errResp = []string{strings.TrimSpace(errs[1])}
			}

			var resp interface{}
			if len(errs) >= 3 && errs[2] != "" {
				resp = strings.Split(errs[2], ",")
			}

			return JSON(c, httpCode, errResp, resp, nil)
		}
	}

	// Default to internal server error
	return InternalServerError(c, []string{err.Error()}, debug)
}

// isDebugMode checks if debug mode is enabled in context
func isDebugMode(c fiber.Ctx) bool {
	if val, ok := c.Locals("debug").(bool); ok {
		return val
	}
	return false
}

// RecoveryError handles panic recovery
func RecoveryError(c fiber.Ctx, r any) error {
	var err error
	switch v := r.(type) {
	case error:
		err = v
	case string:
		err = fmt.Errorf("%s", v)
	case int:
		err = fmt.Errorf("%d", v)
	default:
		err = fmt.Errorf("unknown error panic: %v", v)
	}
	return HandlerErrorGlobal(c, err)
}
