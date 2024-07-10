package middleware

import (
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

var validatorV9 *validator.Validate

func NotBlank(fl validator.FieldLevel) bool {
	if len(strings.Trim(fl.Field().String(), " ")) == 0 {
		return false
	}
	return true
}

func NewValidator() *validator.Validate {
	validatorV9 = validator.New()
	validatorV9.RegisterValidation("not-blank", NotBlank, false)
	return validatorV9
}
