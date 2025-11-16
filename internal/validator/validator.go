package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(data any) error {
	if err := v.validate.Struct(data); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var messages []string
			for _, e := range validationErrors {
				messages = append(messages, fmt.Sprintf("%s: %s", e.Field(), validationErrorMessage(e)))
			}
			return fmt.Errorf("%s", strings.Join(messages, ", "))
		}
		return err
	}
	return nil
}

func validationErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be valid email"
	case "min":
		return fmt.Sprintf("must be at least %s characters", e.Param())
	case "max":
		return fmt.Sprintf("must not exceed %s characters", e.Param())
	default:
		return "unhandled error"
	}
}
