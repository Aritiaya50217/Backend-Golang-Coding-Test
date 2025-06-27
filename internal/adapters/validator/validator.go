package validator

import "github.com/go-playground/validator"

type DefaultValidator struct {
	validate *validator.Validate
}

func NewValidator() *DefaultValidator {
	return &DefaultValidator{
		validate: validator.New(),
	}
}

func (v *DefaultValidator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}
