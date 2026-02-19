package validator

import (
	"app/xonvera-core/internal/core/domain"
	"app/xonvera-core/internal/utils"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var fieldNameCache = make(map[string]map[string]string)
var cacheMutex sync.RWMutex

func HandlerBindingError(c fiber.Ctx, obj any, shouldType string, skips ...string) (res []string) {
	var err error
	switch shouldType {
	case HandlerQuery:
		err = c.Bind().Query(obj)
	case HandlerBody:
		err = c.Bind().Body(obj)
	}
	if err != nil {
		return []string{fmt.Sprintf("payload: %s", err.Error())}
	}

	// Handle PaginationRequest offset calculation
	if paginationReq, ok := obj.(*domain.PaginationRequest); ok {
		if paginationReq.Page < 1 {
			paginationReq.Page = 1
		}
		if paginationReq.Limit < 1 {
			paginationReq.Limit = defaultLimit
		} else if paginationReq.Limit > maxLimit {
			paginationReq.Limit = maxLimit
		}

		paginationReq.Offset = uint64((paginationReq.Page - 1) * paginationReq.Limit)
	}

	return Validation(obj, skips...)
}

func Validation(obj interface{}, skips ...string) []string {
	v := validator.New()
	Init(v)
	err := v.Struct(obj)
	return ErrorHandle(obj, err, skips...)
}

// Extract field name from struct field tags
func extractFieldName(structField reflect.StructField) string {
	jsonTag := structField.Tag.Get("json")
	if jsonTag != "" && jsonTag != "-" {
		return strings.Split(jsonTag, ",")[0]
	}

	queryTag := structField.Tag.Get("query")
	if queryTag != "" && queryTag != "-" {
		return strings.Split(queryTag, ",")[0]
	}

	formTag := structField.Tag.Get("form")
	if formTag != "" && formTag != "-" {
		return strings.Split(formTag, ",")[0]
	}

	return structField.Name
}

// Map field names to their JSON tags.
func getJSONFieldName(obj interface{}) map[string]string {
	objValue := reflect.ValueOf(obj)
	objType := objValue.Type()

	cacheMutex.RLock()
	if cached, ok := fieldNameCache[objType.String()]; ok {
		cacheMutex.RUnlock()
		return cached
	}
	cacheMutex.RUnlock()

	// Handle a pointer to struct if necessary
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}

	var resp = make(map[string]string, objType.NumField())

	// Iterate over struct fields
	for i := 0; i < objType.NumField(); i++ {
		structField := objType.Field(i)
		field := extractFieldName(structField)

		// Handle embedded (anonymous) structs
		if structField.Anonymous {
			embeddedField := objValue.Field(i)
			if embeddedField.IsValid() {
				embeddedResp := getJSONFieldName(embeddedField.Interface())
				maps.Copy(resp, embeddedResp)
			}
		}

		resp[structField.Name] = field
	}

	cacheMutex.Lock()
	fieldNameCache[objType.String()] = resp
	cacheMutex.Unlock()

	return resp // Return the original field name by default
}

func ErrorHandle(obj interface{}, err error, skips ...string) []string {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		var required []string
		var errorsNew []string
		var fieldNames = getJSONFieldName(obj)
		for _, f := range errs {
			fieldName := fieldNames[f.Field()]
			if utils.InArray(skips, fieldName) {
				continue
			}
			var str string
			switch f.Tag() {
			case "required":
				required = append(required, fieldName)
			case "min":
				i, _ := strconv.Atoi(f.Param())
				str = fmt.Sprintf("minimum field '%s' is %s", fieldName, utils.DecimalSeparator(i))
			case "max":
				i, _ := strconv.Atoi(f.Param())
				str = fmt.Sprintf("maximum field '%s' is %s", fieldName, utils.DecimalSeparator(i))
			case "stringNumberOnly":
				str = fmt.Sprintf("'%s' just support letter and number", fieldName)
			case "oneof":
				value := strings.Join(strings.Split(f.Param(), " "), " or ")
				str = fmt.Sprintf("field '%s' value must be one of character '%s'", fieldName, value)
			case "email":
				str = fmt.Sprintf("'%s' invalid email format", fieldName)
			default:
				str = fmt.Sprintf("error '%s' is '%s'", fieldName, f.Value())
			}

			if str != "" {
				errorsNew = append(errorsNew, str)
			}
		}

		for _, v := range required {
			errorsNew = append(errorsNew, fmt.Sprintf("'%s' is required", v))
		}

		return errorsNew
	}
	return nil
}
