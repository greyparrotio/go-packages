package validate


import (
	"fmt"
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
