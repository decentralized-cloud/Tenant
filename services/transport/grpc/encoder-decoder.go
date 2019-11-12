// Package grpc implements functions to expose tenant service endpoint using GRPC protocol.
package grpc

import (
	"context"

	tenantGRPCContract "github.com/decentralized-cloud/tenant/contract/grpc/go"
	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/business"
	"github.com/micro-business/go-core/common"
	commonErrors "github.com/micro-business/go-core/system/errors"
	"github.com/thoas/go-funk"
)

// decodeCreateTenantRequest decodes CreateTenant request message from GRPC object to business object
// context: Mandatory The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeCreateTenantRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*tenantGRPCContract.CreateTenantRequest)

	return &business.CreateTenantRequest{
		Tenant: models.Tenant{
			Name: castedRequest.Tenant.Name,
		}}, nil
}

// encodeCreateTenantResponse encodes CreateTenant response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeCreateTenantResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.CreateTenantResponse)

	if castedResponse.Err == nil {
		return &tenantGRPCContract.CreateTenantResponse{
			Error:    tenantGRPCContract.Error_NO_ERROR,
			TenantID: castedResponse.TenantID,
			Tenant: &tenantGRPCContract.Tenant{
				Name: castedResponse.Tenant.Name,
			},
		}, nil
	}

	return &tenantGRPCContract.CreateTenantResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

// decodeReadTenantRequest decodes ReadTenant request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeReadTenantRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*tenantGRPCContract.ReadTenantRequest)

	return &business.ReadTenantRequest{
		TenantID: castedRequest.TenantID,
	}, nil
}

// encodeReadTenantResponse encodes ReadTenant response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeReadTenantResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.ReadTenantResponse)

	if castedResponse.Err == nil {
		return &tenantGRPCContract.ReadTenantResponse{
			Error: tenantGRPCContract.Error_NO_ERROR,
			Tenant: &tenantGRPCContract.Tenant{
				Name: castedResponse.Tenant.Name,
			},
		}, nil
	}

	return &tenantGRPCContract.ReadTenantResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

// decodeUpdateTenantRequest decodes UpdateTenant request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeUpdateTenantRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*tenantGRPCContract.UpdateTenantRequest)

	return &business.UpdateTenantRequest{
		TenantID: castedRequest.TenantID,
		Tenant: models.Tenant{
			Name: castedRequest.Tenant.Name,
		}}, nil
}

// encodeUpdateTenantResponse encodes UpdateTenant response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeUpdateTenantResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.UpdateTenantResponse)

	if castedResponse.Err == nil {
		return &tenantGRPCContract.UpdateTenantResponse{
			Error: tenantGRPCContract.Error_NO_ERROR,
			Tenant: &tenantGRPCContract.Tenant{
				Name: castedResponse.Tenant.Name,
			},
		}, nil
	}

	return &tenantGRPCContract.UpdateTenantResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

// decodeDeleteTenantRequest decodes DeleteTenant request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeDeleteTenantRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*tenantGRPCContract.DeleteTenantRequest)

	return &business.DeleteTenantRequest{
		TenantID: castedRequest.TenantID,
	}, nil
}

// encodeDeleteTenantResponse encodes DeleteTenant response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeDeleteTenantResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.DeleteTenantResponse)
	if castedResponse.Err == nil {
		return &tenantGRPCContract.DeleteTenantResponse{
			Error: tenantGRPCContract.Error_NO_ERROR,
		}, nil
	}

	return &tenantGRPCContract.DeleteTenantResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

// decodeSearchRequest decodes Search request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrongw
func decodeSearchRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*tenantGRPCContract.SearchRequest)
	sortingOptions := []common.SortingOptionPair{}

	if len(castedRequest.SortingOptions) > 0 {
		sortingOptions = funk.Map(
			castedRequest.SortingOptions,
			func(sortingOption *tenantGRPCContract.SortingOptionPair) common.SortingOptionPair {
				direction := common.Ascending

				if sortingOption.Direction == tenantGRPCContract.SortingDirection_DESCENDING {
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
		*pagination.After = castedRequest.Pagination.After
	}

	if castedRequest.Pagination.HasFirst {
		*pagination.First = int(castedRequest.Pagination.First)
	}

	if castedRequest.Pagination.HasBefore {
		*pagination.Before = castedRequest.Pagination.Before
	}

	if castedRequest.Pagination.HasLast {
		*pagination.Last = int(castedRequest.Pagination.Last)
	}

	return &business.SearchRequest{
		Pagination:     pagination,
		TenantIDs:      castedRequest.TenantIDs,
		SortingOptions: sortingOptions,
	}, nil
}

// encodeSearchResponse encodes Search response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeSearchResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*business.SearchResponse)
	if castedResponse.Err == nil {
		return &tenantGRPCContract.SearchResponse{
			Error:           tenantGRPCContract.Error_NO_ERROR,
			HasPreviousPage: castedResponse.HasPreviousPage,
			HasNextPage:     castedResponse.HasNextPage,
			TotalCount:      castedResponse.TotalCount,
			Tenants: funk.Map(castedResponse.Tenants, func(tenant models.TenantWithCursor) *tenantGRPCContract.TenantWithCursor {
				return &tenantGRPCContract.TenantWithCursor{
					TenantID: tenant.TenantID,
					Tenant: &tenantGRPCContract.Tenant{
						Name: tenant.Tenant.Name,
					},
					Cursor: tenant.Cursor,
				}
			}).([]*tenantGRPCContract.TenantWithCursor),
		}, nil
	}

	return &tenantGRPCContract.SearchResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

func mapError(err error) tenantGRPCContract.Error {
	if business.IsUnknownError(err) {
		return tenantGRPCContract.Error_UNKNOWN
	}

	if business.IsTenantAlreadyExistsError(err) {
		return tenantGRPCContract.Error_TENANT_ALREADY_EXISTS
	}

	if business.IsTenantNotFoundError(err) {
		return tenantGRPCContract.Error_TENANT_NOT_FOUND
	}

	if commonErrors.IsArgumentNilError(err) || commonErrors.IsArgumentError(err) {
		return tenantGRPCContract.Error_BAD_REQUEST
	}

	panic("Error type undefined.")
}
