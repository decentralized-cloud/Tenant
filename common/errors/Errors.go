// Package contracts defines the different tenant repository contracts
package errors

import "fmt"

// ArgumentError indicates that the provided input argument is invalid.
type ArgumentError struct {
	argumentName string
	errorMessage string
	message      string
}

// Error returns message for the TenantAlreadyExistsError error type
// Returns the error nessage
func (e ArgumentError) Error() string {
	return e.message
}

// NewArgumentError creates a new ArgumentError error
func NewArgumentError(argumentName, errorMessage string) error {
	return ArgumentError{
		argumentName: argumentName,
		errorMessage: errorMessage,
		message:      fmt.Sprintf("Argument \"%s\" is invalid. Error message: %s", argumentName, errorMessage),
	}
}
