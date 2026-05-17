package val

import "github.com/go-playground/validator/v10"

func SwitchValidateErr(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + "this field is required"
	case "max":
		return err.Field() + "this field cannot be more than " + err.Param() + " allowed"
	case "min":
		return err.Field() + "this field cannot be lower than " + err.Param() + "allowed"
	case "email":
		return err.Field() + "this field must be a valid email address"
	case "len":
		return err.Field() + "this field must be " + err.Param() + "."
	case "alpha":
		return err.Field() + "this field must be a valid alpha character "
	case "numeric":
		return err.Field() + "this field must be a valid number "
	default:
		return err.Tag() + err.Field()

	}
}
