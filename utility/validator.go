package utility

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

func ExtractValidationErrors(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, fmt.Sprintf("Error in field '%s': %s", e.Field(), e.Tag()))
	}
	return errors
}

func PhoneNumberValidation(fl validator.FieldLevel) bool {
	fieldVal := fl.Field().String()
	match, _ := regexp.MatchString("^[0-9+-]+$", fieldVal)
	return match
}

func EmailValidation(fl validator.FieldLevel) bool {
	fieldVal := fl.Field().String()
	match, _ := regexp.MatchString("^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", fieldVal)
	return match
}

func AlphaSpace(fl validator.FieldLevel) bool {
	fieldVal := fl.Field().String()
	match, _ := regexp.MatchString("^[a-zA-Z\\s]+$", fieldVal)
	return match
}
