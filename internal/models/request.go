package models

import (
	"github.com/go-playground/validator"
)

type CreateLocationRequestBody struct {
	Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180"`
}

// type UpdateLocationRequestBody struct {
// 	Latitude  *float64 `json:"latitude" validate:"omitempty,min=-90,max=90"`
// 	Longitude *float64 `json:"longitude" validate:"omitempty,min=-180,max=180"`
// }

func (b *CreateLocationRequestBody) Validate() error {
	validate := validator.New()
	return validate.Struct(b)
}

// func (b *UpdateLocationRequestBody) Validate() error {
// 	validate := validator.New()
// 	return validate.Struct(b)
// }
