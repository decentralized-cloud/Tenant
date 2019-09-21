// Package contract defines the different tenant business contracts
package contract

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// Validate validates the CreateTenantRequest model and return error if the validation failes
// Returns error if validation failes
func (val *CreateTenantRequest) Validate() error {
	// TODO: mortezaalizadeh: 16/09/2019: Should replace following code with nested validation
	return val.Tenant.Validate()
}

// Validate validates the CreateTenantResponse model and return error if the validation failes
// Returns error if validation failes
func (val *CreateTenantResponse) Validate() error {
	return validation.ValidateStruct(val,
		// TenantID cannot be empty
		validation.Field(&val.TenantID, validation.Required),
	)
}

// Validate validates the ReadTenantRequest model and return error if the validation failes
// Returns error if validation failes
func (val *ReadTenantRequest) Validate() error {
	return validation.ValidateStruct(val,
		// TenantID cannot be empty
		validation.Field(&val.TenantID, validation.Required),
	)
}

// Validate validates the ReadTenantResponse model and return error if the validation failes
// Returns error if validation failes
func (val *ReadTenantResponse) Validate() error {
	return validation.ValidateStruct(val,
		// Validate Tenant using its own validation rules
		validation.Field(&val.Tenant),
	)
}

// Validate validates the UpdateTenantRequest model and return error if the validation failes
// Returns error if validation failes
func (val *UpdateTenantRequest) Validate() error {
	return validation.ValidateStruct(val,
		// TenantID cannot be empty
		validation.Field(&val.TenantID, validation.Required),
		// Validate Tenant using its own validation rules
		validation.Field(&val.Tenant),
	)
}

// Validate validates the UpdateTenantResponse model and return error if the validation failes
// Returns error if validation failes
func (val *UpdateTenantResponse) Validate() error {
	return nil
}

// Validate validates the DeleteTenantRequest model and return error if the validation failes
// Returns error if validation failes
func (val *DeleteTenantRequest) Validate() error {
	return validation.ValidateStruct(val,
		// TenantID cannot be empty
		validation.Field(&val.TenantID, validation.Required),
	)
}

// Validate validates the DeleteTenantResponse model and return error if the validation failes
// Returns error if validation failes
func (val *DeleteTenantResponse) Validate() error {
	return nil
}
