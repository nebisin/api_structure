package request

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidateInput(input interface{}) map[string]string {
	validate := validator.New()

	if err := validate.Struct(input); err != nil {
		errs := err.(validator.ValidationErrors)
		errorMap := make(map[string]string)
		for _, fieldError := range errs {
			key := strings.ToLower(fieldError.Field())
			switch {
			case fieldError.Tag() == "required":
				errorMap[key] = "must be provided"
			case fieldError.Tag() == "unique":
				errorMap[key] = "must not contain duplicate values"
			case fieldError.Tag() == "gt":
				errorMap[key] = fmt.Sprintf("must be greater than %s", fieldError.Param())
			case fieldError.Tag() == "lt":
				errorMap[key] = fmt.Sprintf("must be less than %s", fieldError.Param())
			case fieldError.Tag() == "oneof":
				errorMap[key] = fmt.Sprintf("must be one of %s", fieldError.Param())
			case fieldError.Tag() == "max":
				errorMap[key] = fmt.Sprintf("length must be not more then %s", fieldError.Param())
			case fieldError.Tag() == "min":
				errorMap[key] = fmt.Sprintf("length must be minimum %s long", fieldError.Param())
			case fieldError.Tag() == "email":
				errorMap[key] = "must be a valid email"
			default:
				errorMap[key] = fmt.Sprint(fieldError.Error())
			}
		}
		return errorMap
	}
	return nil
}
