package validations

import (
	"context"
	"errors"
	"reflect"
	"stock_broker_application/src/constants"
	"unicode"

	"gopkg.in/go-playground/validator.v9"
)

var custValidator *validator.Validate

func NewCustomValidator(ctx context.Context) {
	custValidator = validator.New()

	custValidator.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get(constants.JsonConfig)
	})
	custValidator.RegisterValidation(constants.PasswordValidation, ValidatePasswordStruct)

}

func GetCustomValidator(ctx context.Context) *validator.Validate {
	if custValidator == nil {
		NewCustomValidator(ctx)
	}
	return custValidator
}

// ValidatePasswordStruct is a custom validation function for password format.
func ValidatePasswordStruct(fl validator.FieldLevel) bool {
	input := fl.Field().String()

	if err := validateCustomPasswordFormat(input); err != nil {
		return false
	}

	return true
}

func validateCustomPasswordFormat(input string) error {
	hasAlpha := false
	hasNumeric := false

	for _, char := range input {
		if unicode.IsLetter(char) {
			hasAlpha = true
		} else if unicode.IsDigit(char) {
			hasNumeric = true
		}

		if hasAlpha && hasNumeric {
			return nil
		}
	}

	return errors.New(constants.ErrorValidatePassword)
}
