package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitializeValidator() {
	validate = validator.New()
}

func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)
	var errorMessages []string
	for _, e := range validationErrors {
		errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' failed validation: %s", e.Field(), e.ActualTag()))
	}
	return fmt.Errorf("validation failed: %s", errorMessages)
}
