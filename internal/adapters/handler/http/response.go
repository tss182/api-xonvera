package http

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type (
	Map map[string]interface{}

	Resp struct {
		RequestID string      `json:"request_id"`
		Errors    []string    `json:"errors,omitempty"`
		Data      interface{} `json:"data"`
		Meta      interface{} `json:"meta,omitempty"`
	}
)

func JSON(c *fiber.Ctx, httpStatus int, errs []string, data, meta interface{}) error {
	if httpStatus < 100 || httpStatus > 599 {
		httpStatus = http.StatusInternalServerError
	}

	var resp interface{}
	var requestID string

	if req, ok := c.Locals("request_id").(string); ok {
		requestID = req
	}

	if data == nil {
		data = Map{}
	}

	resp = Resp{
		RequestID: requestID,
		Errors:    errs,
		Data:      data,
		Meta:      meta,
	}

	return c.Status(httpStatus).JSON(resp)
}

func OK(c *fiber.Ctx, resp interface{}) error {
	return JSON(c, http.StatusOK, nil, resp, nil)
}

// func Pagination[T any](c *fiber.Ctx, resp *pageResp[T]) error {
// 	if resp.Data == nil {
// 		resp.Data = []T{}
// 	}
// 	return JSON(c, http.StatusOK, nil, resp.Data, resp.Meta)
// }

func NoAuth(c *fiber.Ctx) error {
	return JSON(c, http.StatusUnauthorized, []string{"unauthorized"}, nil, nil)
}

func BadRequest(c *fiber.Ctx, errs []string) error {
	return JSON(c, http.StatusBadRequest, errs, nil, nil)
}

func NotFound(c *fiber.Ctx, errs []string) error {
	return JSON(c, http.StatusNotFound, errs, nil, nil)
}

func InternalServerError(c *fiber.Ctx, errs []string, debugError bool) error {
	if errs == nil || !debugError {
		errs = []string{"internal server error"}
	}
	return JSON(c, http.StatusInternalServerError, errs, nil, nil)
}

func ErrorLimited(c *fiber.Ctx, errs []string) error {
	return JSON(c, http.StatusTooManyRequests, errs, nil, nil)
}

func HandlerErrorGlobal(c *fiber.Ctx, err error) error {
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound(c, []string{
			"data not found",
		})
	}

	var debug = false

	if val, ok := c.Locals("debug").(bool); ok {
		debug = val
	}

	errs := strings.Split(err.Error(), ":")
	if len(errs) >= 2 {
		httpCode, _ := strconv.Atoi(errs[0])
		if httpCode >= 200 && httpCode <= 599 {
			var errResp []string
			if errs[1] != "" {
				errResp = []string{
					strings.TrimSpace(errs[1]),
				}
			}
			var resp interface{}
			if len(errs) >= 3 {
				resp = strings.Split(errs[2], ",")
			}
			return JSON(c, httpCode, errResp, resp, nil)
		}
	}

	return InternalServerError(c, []string{err.Error()}, debug)
}

func RecoveryError(c *fiber.Ctx, r any) error {
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
