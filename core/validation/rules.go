package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func uuidValidator(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if _, err := uuid.Parse(field); err != nil {
		return true
	}
	return false
}
