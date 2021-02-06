// Package repository implements different repository services required by the project service
package repository

import (
	"github.com/decentralized-cloud/project/models"
	"github.com/micro-business/go-core/common"
)

// CreateProjectRequest contains the request to create a new project
type CreateProjectRequest struct {
	Project models.Project
}

// CreateProjectResponse contains the result of creating a new project
type CreateProjectResponse struct {
	ProjectID string
	Project   models.Project
	Cursor    string
}

// ReadProjectRequest contains the request to read an existing project
type ReadProjectRequest struct {
	ProjectID string
}

// ReadProjectResponse contains the result of reading an existing project
type ReadProjectResponse struct {
	Project models.Project
}

// UpdateProjectRequest contains the request to update an existing project
type UpdateProjectRequest struct {
	ProjectID string
	Project   models.Project
}

// UpdateProjectResponse contains the result of updating an existing project
type UpdateProjectResponse struct {
	Project models.Project
	Cursor  string
}

// DeleteProjectRequest contains the request to delete an existing project
type DeleteProjectRequest struct {
	ProjectID string
}

// DeleteProjectResponse contains the result of deleting an existing project
type DeleteProjectResponse struct {
}

// SearchRequest contains the filter criteria to look for existing projects
type SearchRequest struct {
	Pagination     common.Pagination
	SortingOptions []common.SortingOptionPair
	ProjectIDs     []string
}

// SearchResponse contains the list of the projects that matched the result
type SearchResponse struct {
	HasPreviousPage bool
	HasNextPage     bool
	TotalCount      int64
	Projects        []models.ProjectWithCursor
}
