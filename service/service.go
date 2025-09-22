package service

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrValidation       = errors.New("Invalid Input")
	ErrEmailConflict    = errors.New("Email Already Exists")
	ErrUsernameConflict = errors.New("Username Already Exists")
	ErrUnauthorized     = errors.New("Unauthorized")
	ErrEmailNotVerified = errors.New("Email Not Verified")
)

type ValidationErrorMap = validation.Errors

// ExtractValidationErrorsMap tries to convert an error into a ValidationErrorMap.
func ExtractValidationErrorsMap(err error) (ValidationErrorMap, bool) {
	var vem ValidationErrorMap
	ok := errors.As(err, &vem)
	return vem, ok
}
