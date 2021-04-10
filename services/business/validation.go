// Package business implements different business services required by the project service
package business

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Validate validates the CreateProjectRequest model and return error if the validation failes
// Returns error if validation failes
func (val CreateProjectRequest) Validate() error {
	return validation.ValidateStruct(&val,
		// Email must be provided
		validation.Field(&val.UserEmail, validation.Required, is.Email),
		// Validate Project using its own validation rules
		validation.Field(&val.Project),
	)
}

// Validate validates the ReadProjectRequest model and return error if the validation failes
// Returns error if validation failes
func (val ReadProjectRequest) Validate() error {
	return validation.ValidateStruct(&val,
		// Email must be provided
		validation.Field(&val.UserEmail, validation.Required, is.Email),
		// ProjectID cannot be empty
		validation.Field(&val.ProjectID, validation.Required),
	)
}

// Validate validates the UpdateProjectRequest model and return error if the validation failes
// Returns error if validation failes
func (val UpdateProjectRequest) Validate() error {
	return validation.ValidateStruct(&val,
		// Email must be provided
		validation.Field(&val.UserEmail, validation.Required, is.Email),
		// ProjectID cannot be empty
		validation.Field(&val.ProjectID, validation.Required),
		// Validate Project using its own validation rules
		validation.Field(&val.Project),
	)
}

// Validate validates the DeleteProjectRequest model and return error if the validation failes
// Returns error if validation failes
func (val DeleteProjectRequest) Validate() error {
	return validation.ValidateStruct(&val,
		// Email must be provided
		validation.Field(&val.UserEmail, validation.Required, is.Email),
		// ProjectID cannot be empty
		validation.Field(&val.ProjectID, validation.Required),
	)
}

// Validate validates the ListProjectsRequest model and return error if the validation failes
// Returns error if validation failes
func (val ListProjectsRequest) Validate() error {
	return nil
}
