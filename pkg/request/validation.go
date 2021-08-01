package request

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidateInput(input interface{}) map[string]string {
	validate := validator.New()

	if err := validate.Struct(input); err != nil {
		errs := err.(validator.ValidationErrors)
		errorMap := make(map[string]string)
		for _, fieldError := range errs {
			switch {
			case fieldError.Tag() == "required":
				errorMap[strings.ToLower(fieldError.Field())] = "must be provided"
			case fieldError.Tag() == "unique":
				errorMap[strings.ToLower(fieldError.Field())] = "must not contain duplicate values"
			}
		}
		return errorMap
	}
	return nil
}
