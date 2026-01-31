package validator

import "github.com/go-playground/validator/v10"

const (
	HandlerQuery = "query"
	HandlerBody  = "body"
)

var (
	defaultLimit uint8
	maxLimit     uint8
)

func SetPaginationDefaults(defaultLim, maxLim int) {
	defaultLimit = uint8(defaultLim)
	maxLimit = uint8(maxLim)
}

func Init(v *validator.Validate) {
	_ = v.RegisterValidation("stringNumberOnly", StringNumberOnly)
	_ = v.RegisterValidation("stringNumberRequired", StringNumberRequired)
}
