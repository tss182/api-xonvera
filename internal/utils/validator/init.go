package validator

import "github.com/go-playground/validator/v10"
const (
	HandlerQuery = "query"
	HandlerBody  = "body"
)

func Init(v *validator.Validate) {
	_ = v.RegisterValidation("stringNumberOnly", StringNumberOnly)
	_ = v.RegisterValidation("stringNumberRequired", StringNumberRequired)
}
