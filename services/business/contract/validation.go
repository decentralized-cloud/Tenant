// Package contract defines the different tenant business contracts
package contract

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// Validate validates the CreateTenantRequest model and return error if the validation failes
// Returns error if validation failes
func (val CreateTenantRequest) Validate() error {
	return validation.ValidateStruct(&val,
		// Validate Tenant using its own validation rules
		validation.Field(&val.Tenant),
	)
}

// Validate validates the ReadTenantRequest model and return error if the validation failes
// Returns error if validation failes
func (val ReadTenantRequest) Validate() error {
	return validation.ValidateStruct(&val,
		// TenantID cannot be empty
		validation.Field(&val.TenantID, validation.Required),
	)
}

// Validate validates the UpdateTenantRequest model and return error if the validation failes
// Returns error if validation failes
func (val UpdateTenantRequest) Validate() error {
	return validation.ValidateStruct(&val,
		// TenantID cannot be empty
		validation.Field(&val.TenantID, validation.Required),
		// Validate Tenant using its own validation rules
		validation.Field(&val.Tenant),
	)
}

// Validate validates the DeleteTenantRequest model and return error if the validation failes
// Returns error if validation failes
func (val DeleteTenantRequest) Validate() error {
	return validation.ValidateStruct(&val,
		// TenantID cannot be empty
		validation.Field(&val.TenantID, validation.Required),
	)
}
