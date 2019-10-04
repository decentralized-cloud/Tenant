// Package service implements the different tenant business services
package service

import (
	"context"

	"github.com/decentralized-cloud/tenant/services/business/contract"
	repositoryContract "github.com/decentralized-cloud/tenant/services/repository/contract"
	commonErrors "github.com/micro-business/go-core/system/errors"
)

type tenantService struct {
	repositoryService repositoryContract.TenantRepositoryServiceContract
}

// NewTenantService creates new instance of the TenantService, setting up all dependencies and returns the instance
// repositoryService: Mandatory. Reference to the repository service that can persist the tenant related data
// Returns the new service or error if something goes wrong
func NewTenantService(
	repositoryService repositoryContract.TenantRepositoryServiceContract) (contract.TenantServiceContract, error) {
	if repositoryService == nil {
		return nil, commonErrors.NewArgumentNilError("repositoryService", "repositoryService is required")
	}

	return &tenantService{
		repositoryService: repositoryService,
	}, nil
}

// CreateTenant creates a new tenant.
// context: Mandatory The reference to the context
// request: Mandatory. The request to create a new tenant
// Returns either the result of creating new tenant or error if something goes wrong.
func (service *tenantService) CreateTenant(
	ctx context.Context,
	request *contract.CreateTenantRequest) (*contract.CreateTenantResponse, error) {
	if ctx == nil {
		return &contract.CreateTenantResponse{
			Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
		}, nil
	}

	if request == nil {
		return &contract.CreateTenantResponse{
			Err: commonErrors.NewArgumentNilError("request", "request is required"),
		}, nil
	}

	if err := request.Validate(); err != nil {
		return &contract.CreateTenantResponse{
			Err: commonErrors.NewArgumentErrorWithError("request", "", err),
		}, nil
	}

	response, err := service.repositoryService.CreateTenant(ctx, &repositoryContract.CreateTenantRequest{
		Tenant: request.Tenant,
	})

	if err != nil {
		return &contract.CreateTenantResponse{
			Err: mapRepositoryError(err, ""),
		}, nil
	}

	return &contract.CreateTenantResponse{
		TenantID: response.TenantID,
	}, nil
}

// ReadTenant read an existing tenant
// context: Mandatory The reference to the context
// request: Mandatory. The request to read an existing tenant
// Returns either the result of reading an existing tenant or error if something goes wrong.
func (service *tenantService) ReadTenant(
	ctx context.Context,
	request *contract.ReadTenantRequest) (*contract.ReadTenantResponse, error) {
	if ctx == nil {
		return &contract.ReadTenantResponse{
			Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
		}, nil
	}

	if request == nil {
		return &contract.ReadTenantResponse{
			Err: commonErrors.NewArgumentNilError("request", "request is required"),
		}, nil
	}

	if err := request.Validate(); err != nil {
		return &contract.ReadTenantResponse{
			Err: commonErrors.NewArgumentErrorWithError("request", "", err),
		}, nil
	}

	response, err := service.repositoryService.ReadTenant(ctx, &repositoryContract.ReadTenantRequest{
		TenantID: request.TenantID,
	})

	if err != nil {
		return &contract.ReadTenantResponse{
			Err: mapRepositoryError(err, request.TenantID),
		}, nil
	}

	return &contract.ReadTenantResponse{
		Tenant: response.Tenant,
	}, nil
}

// UpdateTenant update an existing tenant
// context: Mandatory The reference to the context
// request: Mandatory. The request to update an existing tenant
// Returns either the result of updateing an existing tenant or error if something goes wrong.
func (service *tenantService) UpdateTenant(
	ctx context.Context,
	request *contract.UpdateTenantRequest) (*contract.UpdateTenantResponse, error) {
	if ctx == nil {
		return &contract.UpdateTenantResponse{
			Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
		}, nil
	}

	if request == nil {
		return &contract.UpdateTenantResponse{
			Err: commonErrors.NewArgumentNilError("request", "request is required"),
		}, nil
	}

	if err := request.Validate(); err != nil {
		return &contract.UpdateTenantResponse{
			Err: commonErrors.NewArgumentErrorWithError("request", "", err),
		}, nil
	}

	_, err := service.repositoryService.UpdateTenant(ctx, &repositoryContract.UpdateTenantRequest{
		TenantID: request.TenantID,
		Tenant:   request.Tenant,
	})

	if err != nil {
		return &contract.UpdateTenantResponse{
			Err: mapRepositoryError(err, request.TenantID),
		}, nil
	}

	return &contract.UpdateTenantResponse{}, nil
}

// DeleteTenant delete an existing tenant
// context: Mandatory The reference to the context
// request: Mandatory. The request to delete an existing tenant
// Returns either the result of deleting an existing tenant or error if something goes wrong.
func (service *tenantService) DeleteTenant(
	ctx context.Context,
	request *contract.DeleteTenantRequest) (*contract.DeleteTenantResponse, error) {
	if ctx == nil {
		return &contract.DeleteTenantResponse{
			Err: commonErrors.NewArgumentNilError("ctx", "ctx is required"),
		}, nil
	}

	if request == nil {
		return &contract.DeleteTenantResponse{
			Err: commonErrors.NewArgumentNilError("request", "request is required"),
		}, nil
	}

	if err := request.Validate(); err != nil {
		return &contract.DeleteTenantResponse{
			Err: commonErrors.NewArgumentErrorWithError("request", "", err),
		}, nil
	}

	_, err := service.repositoryService.DeleteTenant(ctx, &repositoryContract.DeleteTenantRequest{
		TenantID: request.TenantID,
	})

	if err != nil {
		return &contract.DeleteTenantResponse{
			Err: mapRepositoryError(err, request.TenantID),
		}, nil
	}

	return &contract.DeleteTenantResponse{}, nil
}

func mapRepositoryError(err error, tenantID string) error {
	if repositoryContract.IsTenantAlreadyExistsError(err) {
		return contract.NewTenantAlreadyExistsErrorWithError(err)
	}

	if repositoryContract.IsTenantNotFoundError(err) {
		return contract.NewTenantNotFoundErrorWithError(tenantID, err)
	}

	return contract.NewUnknownErrorWithError("", err)
}
