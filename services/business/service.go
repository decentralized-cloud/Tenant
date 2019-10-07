// Package business implements different business services required by the tenant service
package business

import (
	"context"

	"github.com/decentralized-cloud/tenant/services/repository"
	commonErrors "github.com/micro-business/go-core/system/errors"
)

type businessService struct {
	repositoryService repository.RepositoryContract
}

// NewBusinessService creates new instance of the BusinessService, setting up all dependencies and returns the instance
// repositoryService: Mandatory. Reference to the repository service that can persist the tenant related data
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

// CreateTenant creates a new tenant.
// context: Mandatory The reference to the context
// request: Mandatory. The request to create a new tenant
// Returns either the result of creating new tenant or error if something goes wrong.
func (service *businessService) CreateTenant(
	ctx context.Context,
	request *CreateTenantRequest) (*CreateTenantResponse, error) {
	response, err := service.repositoryService.CreateTenant(ctx, &repository.CreateTenantRequest{
		Tenant: request.Tenant,
	})

	if err != nil {
		return &CreateTenantResponse{
			Err: mapRepositoryError(err, ""),
		}, nil
	}

	return &CreateTenantResponse{
		TenantID: response.TenantID,
	}, nil
}

// ReadTenant read an existing tenant
// context: Mandatory The reference to the context
// request: Mandatory. The request to read an existing tenant
// Returns either the result of reading an existing tenant or error if something goes wrong.
func (service *businessService) ReadTenant(
	ctx context.Context,
	request *ReadTenantRequest) (*ReadTenantResponse, error) {
	response, err := service.repositoryService.ReadTenant(ctx, &repository.ReadTenantRequest{
		TenantID: request.TenantID,
	})

	if err != nil {
		return &ReadTenantResponse{
			Err: mapRepositoryError(err, request.TenantID),
		}, nil
	}

	return &ReadTenantResponse{
		Tenant: response.Tenant,
	}, nil
}

// UpdateTenant update an existing tenant
// context: Mandatory The reference to the context
// request: Mandatory. The request to update an existing tenant
// Returns either the result of updateing an existing tenant or error if something goes wrong.
func (service *businessService) UpdateTenant(
	ctx context.Context,
	request *UpdateTenantRequest) (*UpdateTenantResponse, error) {
	_, err := service.repositoryService.UpdateTenant(ctx, &repository.UpdateTenantRequest{
		TenantID: request.TenantID,
		Tenant:   request.Tenant,
	})

	if err != nil {
		return &UpdateTenantResponse{
			Err: mapRepositoryError(err, request.TenantID),
		}, nil
	}

	return &UpdateTenantResponse{}, nil
}

// DeleteTenant delete an existing tenant
// context: Mandatory The reference to the context
// request: Mandatory. The request to delete an existing tenant
// Returns either the result of deleting an existing tenant or error if something goes wrong.
func (service *businessService) DeleteTenant(
	ctx context.Context,
	request *DeleteTenantRequest) (*DeleteTenantResponse, error) {
	_, err := service.repositoryService.DeleteTenant(ctx, &repository.DeleteTenantRequest{
		TenantID: request.TenantID,
	})

	if err != nil {
		return &DeleteTenantResponse{
			Err: mapRepositoryError(err, request.TenantID),
		}, nil
	}

	return &DeleteTenantResponse{}, nil
}

func mapRepositoryError(err error, tenantID string) error {
	if repository.IsTenantAlreadyExistsError(err) {
		return NewTenantAlreadyExistsErrorWithError(err)
	}

	if repository.IsTenantNotFoundError(err) {
		return NewTenantNotFoundErrorWithError(tenantID, err)
	}

	return NewUnknownErrorWithError("", err)
}
