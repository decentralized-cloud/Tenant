// Package endpoint implements different endpoint services required by the tenant service
package endpoint

import (
	"context"

	"github.com/decentralized-cloud/tenant/services/business"
	"github.com/go-kit/kit/endpoint"
	commonErrors "github.com/micro-business/go-core/system/errors"
)

type endpointCreatorService struct {
	businessService business.BusinessContract
}

// NewEndpointCreatorService creates new instance of the EndpointCreatorService, setting up all dependencies and returns the instance
// businessService: Mandatory. Reference to the instance of the Tenant  service
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

// CreateTenantEndpoint creates Create Tenant endpoint
// Returns the Create Tenant endpoint
func (service *endpointCreatorService) CreateTenantEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if ctx == nil {
			return &business.CreateTenantResponse{
				Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
			}, nil
		}

		if request == nil {
			return &business.CreateTenantResponse{
				Err: commonErrors.NewArgumentNilError("request", "request is required"),
			}, nil
		}

		castedRequest := request.(*business.CreateTenantRequest)
		if err := castedRequest.Validate(); err != nil {
			return &business.CreateTenantResponse{
				Err: commonErrors.NewArgumentErrorWithError("request", "", err),
			}, nil
		}

		return service.businessService.CreateTenant(ctx, castedRequest)
	}
}

// ReadTenantEndpoint creates Read Tenant endpoint
// Returns the Read Tenant endpoint
func (service *endpointCreatorService) ReadTenantEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if ctx == nil {
			return &business.ReadTenantResponse{
				Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
			}, nil
		}

		if request == nil {
			return &business.ReadTenantResponse{
				Err: commonErrors.NewArgumentNilError("request", "request is required"),
			}, nil
		}

		castedRequest := request.(*business.ReadTenantRequest)
		if err := castedRequest.Validate(); err != nil {
			return &business.ReadTenantResponse{
				Err: commonErrors.NewArgumentErrorWithError("request", "", err),
			}, nil
		}

		return service.businessService.ReadTenant(ctx, castedRequest)
	}
}

// UpdateTenantEndpoint creates Update Tenant endpoint
// Returns the Update Tenant endpoint
func (service *endpointCreatorService) UpdateTenantEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if ctx == nil {
			return &business.UpdateTenantResponse{
				Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
			}, nil
		}

		if request == nil {
			return &business.UpdateTenantResponse{
				Err: commonErrors.NewArgumentNilError("request", "request is required"),
			}, nil
		}

		castedRequest := request.(*business.UpdateTenantRequest)
		if err := castedRequest.Validate(); err != nil {
			return &business.UpdateTenantResponse{
				Err: commonErrors.NewArgumentErrorWithError("request", "", err),
			}, nil
		}

		return service.businessService.UpdateTenant(ctx, castedRequest)
	}
}

// DeleteTenantEndpoint creates Delete Tenant endpoint
// Returns the Delete Tenant endpoint
func (service *endpointCreatorService) DeleteTenantEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if ctx == nil {
			return &business.DeleteTenantResponse{
				Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
			}, nil
		}

		if request == nil {
			return &business.DeleteTenantResponse{
				Err: commonErrors.NewArgumentNilError("request", "request is required"),
			}, nil
		}

		castedRequest := request.(*business.DeleteTenantRequest)
		if err := castedRequest.Validate(); err != nil {
			return &business.DeleteTenantResponse{
				Err: commonErrors.NewArgumentErrorWithError("request", "", err),
			}, nil
		}

		return service.businessService.DeleteTenant(ctx, castedRequest)
	}
}
