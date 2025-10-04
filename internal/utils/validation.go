package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func Validator() *validator.Validate {
	validate := validator.New()

	if err := validate.RegisterValidation("password", Password); err != nil {
		return nil
	}

	return validate
}

func Password(field validator.FieldLevel) bool {
	value, ok := field.Field().Interface().(string)
	if ok {
		hasDigit := regexp.MustCompile(`[0-9]`).MatchString(value)
		hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(value)

		if !hasDigit || !hasLetter {
			return false
		}
	}

	return true
}
