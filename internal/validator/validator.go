package validator

import "github.com/go-playground/validator/v10"

type Validator struct {
	Vld *validator.Validate
}

func New() *Validator {
	return &Validator{
		Vld: validator.New(),
	}
}
