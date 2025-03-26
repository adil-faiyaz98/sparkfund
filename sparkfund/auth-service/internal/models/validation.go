package models

import "github.com/go-playground/validator/v10"

// ValidateStruct validates a struct using validator tags
func ValidateStruct(s interface{}) error {
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		return err
	}
	return nil
}