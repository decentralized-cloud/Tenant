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
		return nil, commonErrors.NewArgumentError("repositoryService", "repositoryService is required")
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
		return nil, commonErrors.NewArgumentError("ctx", "ctx is required")
	}

	if request == nil {
		return nil, commonErrors.NewArgumentError("request", "request is required")
	}

	if err := request.Validate(); err != nil {
		return nil, commonErrors.NewArgumentError("request", err.Error())
	}

	response, err := service.repositoryService.CreateTenant(ctx, &repositoryContract.CreateTenantRequest{
		Tenant: request.Tenant,
	})

	if err != nil {
		if _, ok := err.(repositoryContract.TenantAlreadyExistsError); ok {
			return nil, contract.NewTenantAlreadyExistsError()
		}

		return nil, contract.NewUnknownError(err.Error())
	}

	return &contract.CreateTenantResponse{
		TenantID: response.TenantID,
	}, nil
}

// ReadTenant read an existing tenant
// context: Mandatory The reference to the context
// request: Mandatory. The request to read an esiting tenant
// Returns either the result of reading an exiting tenant or error if something goes wrong.
func (service *tenantService) ReadTenant(
	ctx context.Context,
	request *contract.ReadTenantRequest) (*contract.ReadTenantResponse, error) {
	if ctx == nil {
		return nil, commonErrors.NewArgumentError("ctx", "ctx is required")
	}

	if request == nil {
		return nil, commonErrors.NewArgumentError("request", "request is required")
	}

	if err := request.Validate(); err != nil {
		return nil, commonErrors.NewArgumentError("request", err.Error())
	}

	response, err := service.repositoryService.ReadTenant(ctx, &repositoryContract.ReadTenantRequest{
		TenantID: request.TenantID,
	})

	if err != nil {
		if _, ok := err.(repositoryContract.TenantNotFoundError); ok {
			return nil, contract.NewTenantNotFoundError(request.TenantID)
		}

		return nil, contract.NewUnknownError(err.Error())
	}

	return &contract.ReadTenantResponse{
		Tenant: response.Tenant,
	}, nil
}

// UpdateTenant update an existing tenant
// context: Mandatory The reference to the context
// request: Mandatory. The request to update an esiting tenant
// Returns either the result of updateing an exiting tenant or error if something goes wrong.
func (service *tenantService) UpdateTenant(
	ctx context.Context,
	request *contract.UpdateTenantRequest) (*contract.UpdateTenantResponse, error) {
	if ctx == nil {
		return nil, commonErrors.NewArgumentError("ctx", "ctx is required")
	}

	if request == nil {
		return nil, commonErrors.NewArgumentError("request", "request is required")
	}

	if err := request.Validate(); err != nil {
		return nil, commonErrors.NewArgumentError("request", err.Error())
	}

	_, err := service.repositoryService.UpdateTenant(ctx, &repositoryContract.UpdateTenantRequest{
		TenantID: request.TenantID,
		Tenant:   request.Tenant,
	})

	if err != nil {
		if _, ok := err.(repositoryContract.TenantNotFoundError); ok {
			return nil, contract.NewTenantNotFoundError(request.TenantID)
		}

		return nil, contract.NewUnknownError(err.Error())
	}

	return &contract.UpdateTenantResponse{}, nil
}

// DeleteTenant delete an existing tenant
// context: Mandatory The reference to the context
// request: Mandatory. The request to delete an esiting tenant
// Returns either the result of deleting an exiting tenant or error if something goes wrong.
func (service *tenantService) DeleteTenant(
	ctx context.Context,
	request *contract.DeleteTenantRequest) (*contract.DeleteTenantResponse, error) {
	if ctx == nil {
		return nil, commonErrors.NewArgumentError("ctx", "ctx is required")
	}

	if request == nil {
		return nil, commonErrors.NewArgumentError("request", "request is required")
	}

	if err := request.Validate(); err != nil {
		return nil, commonErrors.NewArgumentError("request", err.Error())
	}

	_, err := service.repositoryService.DeleteTenant(ctx, &repositoryContract.DeleteTenantRequest{
		TenantID: request.TenantID,
	})

	if err != nil {
		if _, ok := err.(repositoryContract.TenantNotFoundError); ok {
			return nil, contract.NewTenantNotFoundError(request.TenantID)
		}

		return nil, contract.NewUnknownError(err.Error())
	}

	return &contract.DeleteTenantResponse{}, nil
}
