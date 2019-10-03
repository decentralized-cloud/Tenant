// Package contract defines the different tenant business contracts

package contract

import "fmt"

// UnknownError indicates that the tenant with the given information already exists
type UnknownError struct {
	Message string
	Err     error
}

// Error returns message for the UnknownError error type
// Returns the error nessage
func (e UnknownError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("Unknown error occurred. Error message: %s.", e.Message)
	}

	return fmt.Sprintf("Unknown error occurred. Error message: %s. Error: %s", e.Message, e.Err.Error())
}

func (e UnknownError) Unwrap() error {
	return e.Err
}

func IsUnknownError(err error) bool {
	_, ok := err.(UnknownError)

	return ok
}

// NewUnknownError creates a new UnknownError error
func NewUnknownError(message string) error {
	return UnknownError{
		Message: message,
	}
}

// NewUnknownErrorWithError creates a new UnknownError error
func NewUnknownErrorWithError(message string, err error) error {
	return UnknownError{
		Message: message,
		Err:     err,
	}
}

// TenantAlreadyExistsError indicates that the tenant with the given information already exists
type TenantAlreadyExistsError struct {
	Err error
}

// Error returns message for the TenantAlreadyExistsError error type
// Returns the error nessage
func (e TenantAlreadyExistsError) Error() string {
	if e.Err == nil {
		return "Tenant already exists."
	}

	return fmt.Sprintf("Tenant already exists. Error: %s", e.Err.Error())
}

func (e TenantAlreadyExistsError) Unwrap() error {
	return e.Err
}

func IsTenantAlreadyExistsError(err error) bool {
	_, ok := err.(TenantAlreadyExistsError)

	return ok
}

// NewTenantAlreadyExistsError creates a new TenantAlreadyExistsError error
func NewTenantAlreadyExistsError() error {
	return TenantAlreadyExistsError{}
}

// NewTenantAlreadyExistsErrorWithError creates a new TenantAlreadyExistsError error
func NewTenantAlreadyExistsErrorWithError(err error) error {
	return TenantAlreadyExistsError{
		Err: err,
	}
}

// TenantNotFoundError indicates that the tenant with the given tenantID does not exist
type TenantNotFoundError struct {
	TenantID string
	Err      error
}

// Error returns message for the TenantNotFoundError error type
// Returns the error nessage
func (e TenantNotFoundError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("Tenant not found. TenantID: %s.", e.TenantID)
	}

	return fmt.Sprintf("Tenant not found. TenantID: %s. Error: %s", e.TenantID, e.Err.Error())
}

func (e TenantNotFoundError) Unwrap() error {
	return e.Err
}

func IsTenantNotFoundError(err error) bool {
	_, ok := err.(TenantNotFoundError)

	return ok
}

// NewTenantNotFoundError creates a new TenantNotFoundError error
// tenantID: Mandatory. The tenantID that did not match any existing tenant
func NewTenantNotFoundError(tenantID string) error {
	return TenantNotFoundError{
		TenantID: tenantID,
	}
}

// NewTenantNotFoundErrorWithError creates a new TenantNotFoundError error
// tenantID: Mandatory. The tenantID that did not match any existing tenant
func NewTenantNotFoundErrorWithError(tenantID string, err error) error {
	return TenantNotFoundError{
		TenantID: tenantID,
		Err:      err,
	}
}
