package services

import (
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrUnauthorized = errors.New("Unauthorized Error")
	ErrValidation   = errors.New("Validation Error")
	ErrNotFound     = errors.New("NotFound Error")
)

func ExtractValidationErrors(err error) validation.Errors {
	var e validation.Errors
	if errors.As(err, &e) {
		return ConvertValidationErrorsKeysToCamelCase(e)
	}
	return nil
}

func ConvertValidationErrorsKeysToCamelCase(errors validation.Errors) validation.Errors {
	camelCaseErrors := make(validation.Errors)
	for field, err := range errors {
		if field != "" {
			camelCaseField := strings.ToLower(field[:1]) + field[1:]
			camelCaseErrors[camelCaseField] = err
		}
	}
	return camelCaseErrors
}
