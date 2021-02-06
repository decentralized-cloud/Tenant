// Package models defines the different object models used in Project
package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// Validate validates the Project and return error if the validation failes
// Returns error if validation failes
func (val Project) Validate() error {
	return validation.ValidateStruct(&val,
		// Name cannot be empty
		validation.Field(&val.Name, validation.Required),
	)
}
