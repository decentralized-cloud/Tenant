// Package repository implements different repository services required by the project service
package repository

import "context"

// RepositoryContract declares the repository service that can create new project, read, update
// and delete existing projects.
type RepositoryContract interface {
	// CreateProject creates a new project.
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request to create a new project
	// Returns either the result of creating new project or error if something goes wrong.
	CreateProject(
		ctx context.Context,
		request *CreateProjectRequest) (*CreateProjectResponse, error)

	// ReadProject read an existing project
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request to read an esiting project
	// Returns either the result of reading an exiting project or error if something goes wrong.
	ReadProject(
		ctx context.Context,
		request *ReadProjectRequest) (*ReadProjectResponse, error)

	// UpdateProject update an existing project
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request to update an esiting project
	// Returns either the result of updateing an exiting project or error if something goes wrong.
	UpdateProject(
		ctx context.Context,
		request *UpdateProjectRequest) (*UpdateProjectResponse, error)

	// DeleteProject delete an existing project
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request to delete an esiting project
	// Returns either the result of deleting an exiting project or error if something goes wrong.
	DeleteProject(
		ctx context.Context,
		request *DeleteProjectRequest) (*DeleteProjectResponse, error)

	// Search returns the list of projects that matched the criteria
	// ctx: Mandatory The reference to the context
	// request: Mandatory. The request contains the search criteria
	// Returns the list of projects that matched the criteria
	Search(
		ctx context.Context,
		request *SearchRequest) (*SearchResponse, error)
}
