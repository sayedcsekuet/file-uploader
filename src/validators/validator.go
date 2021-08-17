package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

type AppValidator interface {
	echo.Validator
	ValidateValue(i interface{}, tag string) error
}
type appValidator struct {
	validator *validator.Validate
}

func (av *appValidator) Validate(i interface{}) error {
	if err := av.validator.Struct(i); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return validationErrors
	}
	return nil
}
func (av *appValidator) ValidateValue(i interface{}, tag string) error {
	if err := av.validator.Var(i, tag); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return validationErrors
	}
	return nil
}
func NewAppValidator() AppValidator {
	return &appValidator{validator: validator.New()}
}
