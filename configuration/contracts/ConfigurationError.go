// Package contracts defines Configuration Error
package contracts

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
