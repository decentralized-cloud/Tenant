// Package business implements different business services required by the project service
package business

import "context"

// BusinessContract declares the service that can create new project, read, update
// and delete existing projects.
type BusinessContract interface {
	// CreateProject creates a new project.
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request to create a new project
	// Returns either the result of creating new project or error if something goes wrong.
	CreateProject(
		ctx context.Context,
		request *CreateProjectRequest) (*CreateProjectResponse, error)

	// ReadProject read an existing project
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request to read an existing project
	// Returns either the result of reading an existing project or error if something goes wrong.
	ReadProject(
		ctx context.Context,
		request *ReadProjectRequest) (*ReadProjectResponse, error)

	// UpdateProject update an existing project
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request to update an existing project
	// Returns either the result of updateing an existing project or error if something goes wrong.
	UpdateProject(
		ctx context.Context,
		request *UpdateProjectRequest) (*UpdateProjectResponse, error)

	// DeleteProject delete an existing project
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request to delete an existing project
	// Returns either the result of deleting an existing project or error if something goes wrong.
	DeleteProject(
		ctx context.Context,
		request *DeleteProjectRequest) (*DeleteProjectResponse, error)

	// ListProjects returns the list of projects that matched the criteria
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request contains the search criteria
	// Returns the list of projects that matched the criteria
	ListProjects(
		ctx context.Context,
		request *ListProjectsRequest) (*ListProjectsResponse, error)
}
