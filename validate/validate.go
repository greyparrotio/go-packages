package validate

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidateObject validates a provided object
func ValidateObject(object any) (string, error) {
	val := validator.New()
	val.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})

	// Rejects any string containing markup-enabling characters
	val.RegisterValidation("nohtml", func(fl validator.FieldLevel) bool {
		return !strings.ContainsAny(fl.Field().String(), "<>")
	})

	// This is for fields that ARE supposed to be URLs: restrict to http/https,
	// blocking javascript.
	val.RegisterValidation("safeurl", func(fl validator.FieldLevel) bool {
		raw := fl.Field().String()
		if raw == "" {
			return true // pair with `omitempty` or `required` for presence
		}
		u, err := url.Parse(raw)
		if err != nil {
			return false
		}
		return u.Scheme == "http" || u.Scheme == "https"
	})

	err := val.Struct(object)
	if err != nil {
		var messages string
		for _, validationError := range err.(validator.ValidationErrors) {
			message := fmt.Sprintf("Invalid Field Value : [%s] . %v not valid",
				validationError.Field(), validationError.Value())
			messages += ", " + message
		}
		return strings.TrimLeft(messages, ", "), err
	}
	return "", nil
}
