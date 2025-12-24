package validator

import (
	"sync"

	val "github.com/go-playground/validator/v10"
)

type ValidatorItf interface {
	ValidateStruct(data any) error
}

type validator struct {
	validator *val.Validate
}

var (
	instance ValidatorItf
	once     sync.Once
)

func NewValidator() ValidatorItf {
	once.Do(func() {
		instance = &validator{
			validator: val.New(),
		}
	})
	return instance
}

func (v *validator) ValidateStruct(data any) error {
	return v.validator.Struct(data)
}
