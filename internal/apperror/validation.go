package apperror

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ErrValidation(err error) *AppError {
	var errMsgs []ValError

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			var message string
			switch e.Tag() {
			case "required": message = fmt.Sprintf("%s is required", e.Field())
			case "min": message = fmt.Sprintf("%s is required", e.Field())
			case "email": message = fmt.Sprintf("%s must be a valid email", e.Field())
			default: message = fmt.Sprintf("%s is invalid", e.Field())
			}

			errMsgs = append(errMsgs, ValError{
				Field:   e.Field(),
				Message: message,
			})
		}
	} 

	return &AppError{
		Code: 400,
		Message: "Bad Request",
		Detail: errMsgs,
	}
}