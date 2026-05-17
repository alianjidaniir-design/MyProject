package val

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func CheckValidation(req any) error {
	validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		var errVal []string
		if castErr, ok := err.(validator.ValidationErrors); ok {
			for _, v := range castErr {
				errVal = append(errVal, SwitchValidateErr(v))
			}
			return fmt.Errorf("%v", errVal[0])

		}
		return err
	}
	return nil
}

func SwitchValidateErr(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + ": this field is required"
	case "max":
		return err.Field() + " : this field cannot be more than " + err.Param()
	case "min":
		return err.Field() + ": this field cannot be lower than " + err.Param()
	case "email":
		return err.Field() + ": this field must be a valid email address"
	case "len":
		return err.Field() + ": this field must be " + err.Param()
	case "alpha":
		return err.Field() + ": this field must be a valid alpha character "
	case "numeric":
		return err.Field() + ": this field must be a valid number "
	default:
		return err.Field() + " validation failed"

	}
}
