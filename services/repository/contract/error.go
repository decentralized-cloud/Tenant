// Package contract defines the different tenant repository contracts
package contract

import "fmt"

// UnknownError indicates that the tenant with the given information already exists
type UnknownError struct {
	errorMessage string
	message      string
}

// Error returns message for the UnknownError error type
// Returns the error nessage
func (e UnknownError) Error() string {
	return e.message
}

// NewUnknownError creates a new UnknownError error
func NewUnknownError(errorMessage string) error {
	return UnknownError{
		errorMessage: errorMessage,
		message:      fmt.Sprintf("Unknow error occurs. Error message is: %s", errorMessage),
	}
}

// TenantAlreadyExistsError indicates that the tenant with the given information already exists
type TenantAlreadyExistsError struct {
	message string
}

// Error returns message for the TenantAlreadyExistsError error type
// Returns the error nessage
func (e TenantAlreadyExistsError) Error() string {
	return e.message
}

// NewTenantAlreadyExistsError creates a new TenantAlreadyExistsError error
func NewTenantAlreadyExistsError() error {
	return TenantAlreadyExistsError{
		message: fmt.Sprintf("Tenant already exists"),
	}
}

// TenantNotFoundError indicates that the tenant with the given tenantID does not exist
type TenantNotFoundError struct {
	TenantID string
	message  string
}

// Error returns message for the TenantNotFoundError error type
// Returns the error nessage
func (e TenantNotFoundError) Error() string {
	return e.message
}

// NewTenantNotFoundError creates a new TenantNotFoundError error
// tenantID: Mandatory. The tenantID that did not match any existing tenant
func NewTenantNotFoundError(tenantID string) error {
	return TenantNotFoundError{
		TenantID: tenantID,
		message:  fmt.Sprintf("Tenant with tenantID: %s not found", tenantID),
	}
}
