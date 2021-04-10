// Package endpoint implements different endpoint services required by the project service
package endpoint

import (
	"context"

	"github.com/decentralized-cloud/project/models"
	"github.com/decentralized-cloud/project/services/business"
	"github.com/go-kit/kit/endpoint"
	commonErrors "github.com/micro-business/go-core/system/errors"
)

type endpointCreatorService struct {
	businessService business.BusinessContract
}

// NewEndpointCreatorService creates new instance of the EndpointCreatorService, setting up all dependencies and returns the instance
// businessService: Mandatory. Reference to the instance of the Project  service
// Returns the new service or error if something goes wrong
func NewEndpointCreatorService(
	businessService business.BusinessContract) (EndpointCreatorContract, error) {
	if businessService == nil {
		return nil, commonErrors.NewArgumentNilError("businessService", "businessService is required")
	}

	return &endpointCreatorService{
		businessService: businessService,
	}, nil
}

// CreateProjectEndpoint creates Create Project endpoint
// Returns the Create Project endpoint
func (service *endpointCreatorService) CreateProjectEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if ctx == nil {
			return &business.CreateProjectResponse{
				Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
			}, nil
		}

		if request == nil {
			return &business.CreateProjectResponse{
				Err: commonErrors.NewArgumentNilError("request", "request is required"),
			}, nil
		}

		castedRequest := request.(*business.CreateProjectRequest)
		parsedToken := ctx.Value(models.ContextKeyParsedToken).(models.ParsedToken)
		castedRequest.UserEmail = parsedToken.Email

		if err := castedRequest.Validate(); err != nil {
			return &business.CreateProjectResponse{
				Err: commonErrors.NewArgumentErrorWithError("request", "", err),
			}, nil
		}

		return service.businessService.CreateProject(ctx, castedRequest)
	}
}

// ReadProjectEndpoint creates Read Project endpoint
// Returns the Read Project endpoint
func (service *endpointCreatorService) ReadProjectEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if ctx == nil {
			return &business.ReadProjectResponse{
				Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
			}, nil
		}

		if request == nil {
			return &business.ReadProjectResponse{
				Err: commonErrors.NewArgumentNilError("request", "request is required"),
			}, nil
		}

		castedRequest := request.(*business.ReadProjectRequest)
		parsedToken := ctx.Value(models.ContextKeyParsedToken).(models.ParsedToken)
		castedRequest.UserEmail = parsedToken.Email

		if err := castedRequest.Validate(); err != nil {
			return &business.ReadProjectResponse{
				Err: commonErrors.NewArgumentErrorWithError("request", "", err),
			}, nil
		}

		return service.businessService.ReadProject(ctx, castedRequest)
	}
}

// UpdateProjectEndpoint creates Update Project endpoint
// Returns the Update Project endpoint
func (service *endpointCreatorService) UpdateProjectEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if ctx == nil {
			return &business.UpdateProjectResponse{
				Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
			}, nil
		}

		if request == nil {
			return &business.UpdateProjectResponse{
				Err: commonErrors.NewArgumentNilError("request", "request is required"),
			}, nil
		}

		castedRequest := request.(*business.UpdateProjectRequest)
		parsedToken := ctx.Value(models.ContextKeyParsedToken).(models.ParsedToken)
		castedRequest.UserEmail = parsedToken.Email

		if err := castedRequest.Validate(); err != nil {
			return &business.UpdateProjectResponse{
				Err: commonErrors.NewArgumentErrorWithError("request", "", err),
			}, nil
		}

		return service.businessService.UpdateProject(ctx, castedRequest)
	}
}

// DeleteProjectEndpoint creates Delete Project endpoint
// Returns the Delete Project endpoint
func (service *endpointCreatorService) DeleteProjectEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if ctx == nil {
			return &business.DeleteProjectResponse{
				Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
			}, nil
		}

		if request == nil {
			return &business.DeleteProjectResponse{
				Err: commonErrors.NewArgumentNilError("request", "request is required"),
			}, nil
		}

		castedRequest := request.(*business.DeleteProjectRequest)
		parsedToken := ctx.Value(models.ContextKeyParsedToken).(models.ParsedToken)
		castedRequest.UserEmail = parsedToken.Email

		if err := castedRequest.Validate(); err != nil {
			return &business.DeleteProjectResponse{
				Err: commonErrors.NewArgumentErrorWithError("request", "", err),
			}, nil
		}

		return service.businessService.DeleteProject(ctx, castedRequest)
	}
}

// ListProjectsEndpoint creates ListProjects Project endpoint
// Returns the ListProjects Project endpoint
func (service *endpointCreatorService) ListProjectsEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if ctx == nil {
			return &business.ListProjectsResponse{
				Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
			}, nil
		}

		if request == nil {
			return &business.ListProjectsResponse{
				Err: commonErrors.NewArgumentNilError("request", "request is required"),
			}, nil
		}

		castedRequest := request.(*business.ListProjectsRequest)
		parsedToken := ctx.Value(models.ContextKeyParsedToken).(models.ParsedToken)
		castedRequest.UserEmail = parsedToken.Email

		if err := castedRequest.Validate(); err != nil {
			return &business.ListProjectsResponse{
				Err: commonErrors.NewArgumentErrorWithError("request", "", err),
			}, nil
		}

		return service.businessService.ListProjects(ctx, castedRequest)
	}
}
