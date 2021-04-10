// Package grpc implements functions to expose project service endpoint using GRPC protocol.
package grpc

import (
	"context"

	projectGRPCContract "github.com/decentralized-cloud/project/contract/grpc/go"
	"github.com/decentralized-cloud/project/models"
	"github.com/decentralized-cloud/project/services/business"
	"github.com/micro-business/go-core/common"
	commonErrors "github.com/micro-business/go-core/system/errors"
	"github.com/thoas/go-funk"
)

// decodeCreateProjectRequest decodes CreateProject request message from GRPC object to business object
// context: Mandatory The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeCreateProjectRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*projectGRPCContract.CreateProjectRequest)

	return &business.CreateProjectRequest{
		Project: models.Project{
			Name: castedRequest.Project.Name,
		}}, nil
}

// encodeCreateProjectResponse encodes CreateProject response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeCreateProjectResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.CreateProjectResponse)

	if castedResponse.Err == nil {
		return &projectGRPCContract.CreateProjectResponse{
			Error:     projectGRPCContract.Error_NO_ERROR,
			ProjectID: castedResponse.ProjectID,
			Project: &projectGRPCContract.Project{
				Name: castedResponse.Project.Name,
			},
			Cursor: castedResponse.Cursor,
		}, nil
	}

	return &projectGRPCContract.CreateProjectResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

// decodeReadProjectRequest decodes ReadProject request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeReadProjectRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*projectGRPCContract.ReadProjectRequest)

	return &business.ReadProjectRequest{
		ProjectID: castedRequest.ProjectID,
	}, nil
}

// encodeReadProjectResponse encodes ReadProject response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeReadProjectResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.ReadProjectResponse)

	if castedResponse.Err == nil {
		return &projectGRPCContract.ReadProjectResponse{
			Error: projectGRPCContract.Error_NO_ERROR,
			Project: &projectGRPCContract.Project{
				Name: castedResponse.Project.Name,
			},
		}, nil
	}

	return &projectGRPCContract.ReadProjectResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

// decodeUpdateProjectRequest decodes UpdateProject request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeUpdateProjectRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*projectGRPCContract.UpdateProjectRequest)

	return &business.UpdateProjectRequest{
		ProjectID: castedRequest.ProjectID,
		Project: models.Project{
			Name: castedRequest.Project.Name,
		}}, nil
}

// encodeUpdateProjectResponse encodes UpdateProject response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeUpdateProjectResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.UpdateProjectResponse)

	if castedResponse.Err == nil {
		return &projectGRPCContract.UpdateProjectResponse{
			Error: projectGRPCContract.Error_NO_ERROR,
			Project: &projectGRPCContract.Project{
				Name: castedResponse.Project.Name,
			},
			Cursor: castedResponse.Cursor,
		}, nil
	}

	return &projectGRPCContract.UpdateProjectResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

// decodeDeleteProjectRequest decodes DeleteProject request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeDeleteProjectRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*projectGRPCContract.DeleteProjectRequest)

	return &business.DeleteProjectRequest{
		ProjectID: castedRequest.ProjectID,
	}, nil
}

// encodeDeleteProjectResponse encodes DeleteProject response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeDeleteProjectResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.DeleteProjectResponse)
	if castedResponse.Err == nil {
		return &projectGRPCContract.DeleteProjectResponse{
			Error: projectGRPCContract.Error_NO_ERROR,
		}, nil
	}

	return &projectGRPCContract.DeleteProjectResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

// decodeListProjectsRequest decodes ListProjects request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrongw
func decodeListProjectsRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*projectGRPCContract.ListProjectsRequest)
	sortingOptions := []common.SortingOptionPair{}

	if len(castedRequest.SortingOptions) > 0 {
		sortingOptions = funk.Map(
			castedRequest.SortingOptions,
			func(sortingOption *projectGRPCContract.SortingOptionPair) common.SortingOptionPair {
				direction := common.Ascending

				if sortingOption.Direction == projectGRPCContract.SortingDirection_DESCENDING {
					direction = common.Descending
				}

				return common.SortingOptionPair{
					Name:      sortingOption.Name,
					Direction: direction,
				}
			}).([]common.SortingOptionPair)
	}

	pagination := common.Pagination{}

	if castedRequest.Pagination.HasAfter {
		pagination.After = &castedRequest.Pagination.After
	}

	if castedRequest.Pagination.HasFirst {
		first := int(castedRequest.Pagination.First)
		pagination.First = &first
	}

	if castedRequest.Pagination.HasBefore {
		pagination.Before = &castedRequest.Pagination.Before
	}

	if castedRequest.Pagination.HasLast {
		last := int(castedRequest.Pagination.Last)
		pagination.Last = &last
	}

	return &business.ListProjectsRequest{
		Pagination:     pagination,
		ProjectIDs:     castedRequest.ProjectIDs,
		SortingOptions: sortingOptions,
	}, nil
}

// encodeListProjectsResponse encodes ListProjects response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeListProjectsResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.ListProjectsResponse)
	if castedResponse.Err == nil {
		return &projectGRPCContract.ListProjectsResponse{
			Error:           projectGRPCContract.Error_NO_ERROR,
			HasPreviousPage: castedResponse.HasPreviousPage,
			HasNextPage:     castedResponse.HasNextPage,
			TotalCount:      castedResponse.TotalCount,
			Projects: funk.Map(castedResponse.Projects, func(project models.ProjectWithCursor) *projectGRPCContract.ProjectWithCursor {
				return &projectGRPCContract.ProjectWithCursor{
					ProjectID: project.ProjectID,
					Project: &projectGRPCContract.Project{
						Name: project.Project.Name,
					},
					Cursor: project.Cursor,
				}
			}).([]*projectGRPCContract.ProjectWithCursor),
		}, nil
	}

	return &projectGRPCContract.ListProjectsResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

func mapError(err error) projectGRPCContract.Error {
	if commonErrors.IsUnknownError(err) {
		return projectGRPCContract.Error_UNKNOWN
	}

	if commonErrors.IsAlreadyExistsError(err) {
		return projectGRPCContract.Error_PROJECT_ALREADY_EXISTS
	}

	if commonErrors.IsNotFoundError(err) {
		return projectGRPCContract.Error_PROJECT_NOT_FOUND
	}

	if commonErrors.IsArgumentNilError(err) || commonErrors.IsArgumentError(err) {
		return projectGRPCContract.Error_BAD_REQUEST
	}

	return projectGRPCContract.Error_UNKNOWN
}
