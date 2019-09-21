// Package service implements the service that provides endpoints to be used by the transport layer
package service

import (
	"context"

	businessContract "github.com/decentralized-cloud/tenant/services/business/contract"
	"github.com/decentralized-cloud/tenant/services/endpoint/contract"
	"github.com/go-kit/kit/endpoint"
)

// EndpointCreatorService implements the service that creates endpoints to create new tenant,
// read, update and delete existing tenant.
type EndpointCreatorService struct {
	businessService businessContract.TenantServiceContract
}

// NewEndpointCreatorService creates new instance of the EndpointCreatorService, setting up all dependencies and returns the instance
// businessService: Mandatory. Reference to the instance of the Tenant  service
// Returns the new service or error if something goes wrong
func NewEndpointCreatorService(
	businessService businessContract.TenantServiceContract) (contract.EndpointCreatorContract, error) {
	return &EndpointCreatorService{
		businessService: businessService,
	}, nil
}

// CreateTenantEndpoint creates Create Tenant endpoint
// Returns the Create Tenant endpoint
func (service *EndpointCreatorService) CreateTenantEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return service.businessService.CreateTenant(ctx, request.(*businessContract.CreateTenantRequest))
	}
}

// ReadTenantEndpoint creates Read Tenant endpoint
// Returns the Read Tenant endpoint
func (service *EndpointCreatorService) ReadTenantEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return service.businessService.ReadTenant(ctx, request.(*businessContract.ReadTenantRequest))
	}
}

// UpdateTenantEndpoint creates Update Tenant endpoint
// Returns the Update Tenant endpoint
func (service *EndpointCreatorService) UpdateTenantEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return service.businessService.UpdateTenant(ctx, request.(*businessContract.UpdateTenantRequest))
	}
}

// DeleteTenantEndpoint creates Delete Tenant endpoint
// Returns the Delete Tenant endpoint
func (service *EndpointCreatorService) DeleteTenantEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return service.businessService.DeleteTenant(ctx, request.(*businessContract.DeleteTenantRequest))
	}
}
