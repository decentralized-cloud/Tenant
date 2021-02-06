// Package business implements different business services required by the project service
package business

import (
	"context"

	"github.com/decentralized-cloud/project/services/repository"
	commonErrors "github.com/micro-business/go-core/system/errors"
)

type businessService struct {
	repositoryService repository.RepositoryContract
}

// NewBusinessService creates new instance of the BusinessService, setting up all dependencies and returns the instance
// repositoryService: Mandatory. Reference to the repository service that can persist the project related data
// Returns the new service or error if something goes wrong
func NewBusinessService(
	repositoryService repository.RepositoryContract) (BusinessContract, error) {
	if repositoryService == nil {
		return nil, commonErrors.NewArgumentNilError("repositoryService", "repositoryService is required")
	}

	return &businessService{
		repositoryService: repositoryService,
	}, nil
}

// CreateProject creates a new project.
// ctx: Mandatory The reference to the context
// request: Mandatory. The request to create a new project
// Returns either the result of creating new project or error if something goes wrong.
func (service *businessService) CreateProject(
	ctx context.Context,
	request *CreateProjectRequest) (*CreateProjectResponse, error) {
	response, err := service.repositoryService.CreateProject(ctx, &repository.CreateProjectRequest{
		Project: request.Project,
	})

	if err != nil {
		return &CreateProjectResponse{
			Err: mapRepositoryError(err, ""),
		}, nil
	}

	return &CreateProjectResponse{
		ProjectID: response.ProjectID,
		Project:   response.Project,
		Cursor:    response.Cursor,
	}, nil
}

// ReadProject read an existing project
// ctx: Mandatory The reference to the context
// request: Mandatory. The request to read an existing project
// Returns either the result of reading an existing project or error if something goes wrong.
func (service *businessService) ReadProject(
	ctx context.Context,
	request *ReadProjectRequest) (*ReadProjectResponse, error) {
	response, err := service.repositoryService.ReadProject(ctx, &repository.ReadProjectRequest{
		ProjectID: request.ProjectID,
	})

	if err != nil {
		return &ReadProjectResponse{
			Err: mapRepositoryError(err, request.ProjectID),
		}, nil
	}

	return &ReadProjectResponse{
		Project: response.Project,
	}, nil
}

// UpdateProject update an existing project
// ctx: Mandatory The reference to the context
// request: Mandatory. The request to update an existing project
// Returns either the result of updateing an existing project or error if something goes wrong.
func (service *businessService) UpdateProject(
	ctx context.Context,
	request *UpdateProjectRequest) (*UpdateProjectResponse, error) {
	response, err := service.repositoryService.UpdateProject(ctx, &repository.UpdateProjectRequest{
		ProjectID: request.ProjectID,
		Project:   request.Project,
	})

	if err != nil {
		return &UpdateProjectResponse{
			Err: mapRepositoryError(err, request.ProjectID),
		}, nil
	}

	return &UpdateProjectResponse{
		Project: response.Project,
		Cursor:  response.Cursor,
	}, nil
}

// DeleteProject delete an existing project
// ctx: Mandatory The reference to the context
// request: Mandatory. The request to delete an existing project
// Returns either the result of deleting an existing project or error if something goes wrong.
func (service *businessService) DeleteProject(
	ctx context.Context,
	request *DeleteProjectRequest) (*DeleteProjectResponse, error) {
	_, err := service.repositoryService.DeleteProject(ctx, &repository.DeleteProjectRequest{
		ProjectID: request.ProjectID,
	})

	if err != nil {
		return &DeleteProjectResponse{
			Err: mapRepositoryError(err, request.ProjectID),
		}, nil
	}

	return &DeleteProjectResponse{}, nil
}

// Search returns the list of projects that matched the criteria
// ctx: Mandatory The reference to the context
// request: Mandatory. The request contains the search criteria
// Returns the list of projects that matched the criteria
func (service *businessService) Search(
	ctx context.Context,
	request *SearchRequest) (*SearchResponse, error) {
	result, err := service.repositoryService.Search(ctx, &repository.SearchRequest{
		Pagination:     request.Pagination,
		SortingOptions: request.SortingOptions,
		ProjectIDs:     request.ProjectIDs,
	})

	if err != nil {
		return &SearchResponse{
			Err: mapRepositoryError(err, ""),
		}, nil
	}

	return &SearchResponse{
		HasPreviousPage: result.HasPreviousPage,
		HasNextPage:     result.HasNextPage,
		TotalCount:      result.TotalCount,
		Projects:        result.Projects,
	}, nil
}

func mapRepositoryError(err error, projectID string) error {
	if repository.IsProjectAlreadyExistsError(err) {
		return NewProjectAlreadyExistsErrorWithError(err)
	}

	if repository.IsProjectNotFoundError(err) {
		return NewProjectNotFoundErrorWithError(projectID, err)
	}

	return NewUnknownErrorWithError("", err)
}
