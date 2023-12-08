package validation

import (
	goerrors "errors"
	"github.com/go-playground/validator/v10"
)

// ===========================================================================================================
// 											Validator
// ===========================================================================================================

// instance A singleton Validator instance.
var instance *validator.Validate

// ValidatorInstance func for create a new validator for model fields.
func ValidatorInstance() *validator.Validate {
	if instance != nil {
		return instance
	}

	// Create a new validator for a Book model.
	instance := validator.New()

	// Custom validation for uuid.UUID fields. Use `validate:"uuid"`
	_ = instance.RegisterValidation("uuid", uuidValidator)

	return instance
}

// ===========================================================================================================
// 											Validation functions
// ===========================================================================================================

// Validate verify a data struct type.
func Validate(structData interface{}, msgForTag MsgForTagFunc) (map[string][]string, error) {
	validatorInstance := ValidatorInstance()
	var out map[string][]string

	// Validate data
	err := validatorInstance.Struct(structData)
	if err != nil {
		// Determine error type ValidationErrors.
		var ve validator.ValidationErrors
		if goerrors.As(err, &ve) {
			out = make(map[string][]string, len(ve))
			// Parse error to build custom message
			for _, fe := range ve {
				messages := out[fe.Field()]
				message := msgForTag(fe)

				out[fe.Field()] = append(messages, message)
			}
		}
	}

	return out, err
}
