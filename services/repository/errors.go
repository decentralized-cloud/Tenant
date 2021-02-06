// Package repository implements different repository services required by the project service
package repository

import "fmt"

// UnknownError indicates that an unknown error has happened
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

// Unwrap returns the err if provided through NewUnknownErrorWithError function, otherwise returns nil
func (e UnknownError) Unwrap() error {
	return e.Err
}

// IsUnknownError indicates whether the error is of type UnknownError
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

// ProjectAlreadyExistsError indicates that the project with the given information already exists
type ProjectAlreadyExistsError struct {
	Err error
}

// Error returns message for the ProjectAlreadyExistsError error type
// Returns the error nessage
func (e ProjectAlreadyExistsError) Error() string {
	if e.Err == nil {
		return "Project already exists."
	}

	return fmt.Sprintf("Project already exists. Error: %s", e.Err.Error())
}

// Unwrap returns the err if provided through NewProjectAlreadyExistsErrorWithError function, otherwise returns nil
func (e ProjectAlreadyExistsError) Unwrap() error {
	return e.Err
}

// IsProjectAlreadyExistsError indicates whether the error is of type ProjectAlreadyExistsError
func IsProjectAlreadyExistsError(err error) bool {
	_, ok := err.(ProjectAlreadyExistsError)

	return ok
}

// NewProjectAlreadyExistsError creates a new ProjectAlreadyExistsError error
func NewProjectAlreadyExistsError() error {
	return ProjectAlreadyExistsError{}
}

// NewProjectAlreadyExistsErrorWithError creates a new ProjectAlreadyExistsError error
func NewProjectAlreadyExistsErrorWithError(err error) error {
	return ProjectAlreadyExistsError{
		Err: err,
	}
}

// ProjectNotFoundError indicates that the project with the given projectID does not exist
type ProjectNotFoundError struct {
	ProjectID string
	Err       error
}

// Error returns message for the ProjectNotFoundError error type
// Returns the error nessage
func (e ProjectNotFoundError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("Project not found. ProjectID: %s.", e.ProjectID)
	}

	return fmt.Sprintf("Project not found. ProjectID: %s. Error: %s", e.ProjectID, e.Err.Error())
}

// Unwrap returns the err if provided through NewProjectNotFoundErrorWithError function, otherwise returns nil
func (e ProjectNotFoundError) Unwrap() error {
	return e.Err
}

// IsProjectNotFoundError indicates whether the error is of type ProjectNotFoundError
func IsProjectNotFoundError(err error) bool {
	_, ok := err.(ProjectNotFoundError)

	return ok
}

// NewProjectNotFoundError creates a new ProjectNotFoundError error
// projectID: Mandatory. The projectID that did not match any existing project
func NewProjectNotFoundError(projectID string) error {
	return ProjectNotFoundError{
		ProjectID: projectID,
	}
}

// NewProjectNotFoundErrorWithError creates a new ProjectNotFoundError error
// projectID: Mandatory. The projectID that did not match any existing project
func NewProjectNotFoundErrorWithError(projectID string, err error) error {
	return ProjectNotFoundError{
		ProjectID: projectID,
		Err:       err,
	}
}
