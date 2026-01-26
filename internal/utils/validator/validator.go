package validator

import (
	"app/xonvera-core/internal/utils"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func HandlerBindingError(c *fiber.Ctx, obj any, shouldType string, skips ...string) (res []string) {
	var err error
	switch shouldType {
	case HandlerQuery:
		err = c.QueryParser(obj)
	case HandlerBody:
		err = c.BodyParser(obj)
	}
	if err != nil {
		return []string{fmt.Sprintf("payload: %s", err.Error())}
	}

	return Validation(obj, skips...)
}

func Validation(obj interface{}, skips ...string) []string {
	v := validator.New()
	Init(v)
	err := v.Struct(obj)
	return ErrorHandle(obj, err, skips...)
}

// Map field names to their JSON tags.
func getJSONFieldName(obj interface{}) map[string]string {
	objValue := reflect.ValueOf(obj)
	objType := objValue.Type()

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
		var field string
		jsonTag := structField.Tag.Get("json")
		queryTag := structField.Tag.Get("query")
		formTag := structField.Tag.Get("form")
		if jsonTag != "" && jsonTag != "-" {
			field = strings.Split(jsonTag, ",")[0] // Handle response like `json:"field,omitempty"`
		} else if queryTag != "" && queryTag != "-" {
			field = strings.Split(queryTag, ",")[0]
		} else if formTag != "" && formTag != "-" {
			field = strings.Split(formTag, ",")[0]
		} else {
			field = structField.Name
		}

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
