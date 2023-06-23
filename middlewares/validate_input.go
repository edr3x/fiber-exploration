package middlewares

import (
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Field  string `json:"field"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

var validate = validator.New()

func ValidateInput(schema interface{}) []*ErrorResponse {
	var errors []*ErrorResponse

	if err := validate.Struct(schema); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.Field()
			element.Error = err.Tag()
			element.Reason = err.Error()
			errors = append(errors, &element)
		}
	}
	return errors
}
