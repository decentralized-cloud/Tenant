// Package business implements different business services required by the project service
package business

import (
	"github.com/decentralized-cloud/project/models"
	"github.com/micro-business/go-core/common"
)

// CreateProjectRequest contains the request to create a new project
type CreateProjectRequest struct {
	UserEmail string
	Project   models.Project
}

// CreateProjectResponse contains the result of creating a new project
type CreateProjectResponse struct {
	Err       error
	ProjectID string
	Project   models.Project
	Cursor    string
}

// ReadProjectRequest contains the request to read an existing project
type ReadProjectRequest struct {
	UserEmail string
	ProjectID string
}

// ReadProjectResponse contains the result of reading an existing project
type ReadProjectResponse struct {
	Err     error
	Project models.Project
}

// UpdateProjectRequest contains the request to update an existing project
type UpdateProjectRequest struct {
	UserEmail string
	ProjectID string
	Project   models.Project
}

// UpdateProjectResponse contains the result of updating an existing project
type UpdateProjectResponse struct {
	Err     error
	Project models.Project
	Cursor  string
}

// DeleteProjectRequest contains the request to delete an existing project
type DeleteProjectRequest struct {
	UserEmail string
	ProjectID string
}

// DeleteProjectResponse contains the result of deleting an existing project
type DeleteProjectResponse struct {
	Err error
}

// ListProjectsRequest contains the filter criteria to look for existing projects
type ListProjectsRequest struct {
	UserEmail      string
	Pagination     common.Pagination
	SortingOptions []common.SortingOptionPair
	ProjectIDs     []string
}

// ListProjectsResponse contains the list of the projects that matched the result
type ListProjectsResponse struct {
	Err             error
	HasPreviousPage bool
	HasNextPage     bool
	TotalCount      int64
	Projects        []models.ProjectWithCursor
}
