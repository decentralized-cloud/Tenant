// Package services implements the different tenant business services
package services

import (
	"context"

	"github.com/decentralized-cloud/Tenant/business/contracts"
	commonErrors "github.com/decentralized-cloud/Tenant/common/errors"
	repositoryContracts "github.com/decentralized-cloud/Tenant/repository/contracts"
)

// TenantService implements the service that create new tenant, read, update and delete existing tenants.
type TenantService struct {
	repositoryService repositoryContracts.TenantRepositoryServiceContract
}

// NewTenantService creates new instance of the TenantService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewTenantService(
	repositoryService repositoryContracts.TenantRepositoryServiceContract) (contracts.TenantServiceContract, error) {
	return &TenantService{
		repositoryService: repositoryService,
	}, nil
}

// CreateTenant creates a new tenant.
// context: Mandatory The reference to the context
// request: Mandatory. The request to create a new tenant
// Returns either the result of creating new tenant or error if something goes wrong.
func (service *TenantService) CreateTenant(
	ctx context.Context,
	request *contracts.CreateTenantRequest) (*contracts.CreateTenantResponse, error) {
	if ctx == nil {
		return nil, commonErrors.NewArgumentError("ctx", "ctx is required")
	}

	if request == nil {
		return nil, commonErrors.NewArgumentError("request", "request is required")
	}

	if err := request.Validate(); err != nil {
		return nil, commonErrors.NewArgumentError("request", err.Error())
	}

	response, err := service.repositoryService.CreateTenant(ctx, &repositoryContracts.CreateTenantRequest{
		Tenant: request.Tenant,
	})

	if err != nil {
		if _, ok := err.(repositoryContracts.TenantAlreadyExistsError); ok {
			return nil, contracts.NewTenantAlreadyExistsError()
		}

		return nil, contracts.NewUnknownError(err.Error())
	}

	return &contracts.CreateTenantResponse{
		TenantID: response.TenantID,
	}, nil
}

// ReadTenant read an existing tenant
// context: Mandatory The reference to the context
// request: Mandatory. The request to read an esiting tenant
// Returns either the result of reading an exiting tenant or error if something goes wrong.
func (service *TenantService) ReadTenant(
	ctx context.Context,
	request *contracts.ReadTenantRequest) (*contracts.ReadTenantResponse, error) {
	if ctx == nil {
		return nil, commonErrors.NewArgumentError("ctx", "ctx is required")
	}

	if request == nil {
		return nil, commonErrors.NewArgumentError("request", "request is required")
	}

	if err := request.Validate(); err != nil {
		return nil, commonErrors.NewArgumentError("request", err.Error())
	}

	response, err := service.repositoryService.ReadTenant(ctx, &repositoryContracts.ReadTenantRequest{
		TenantID: request.TenantID,
	})

	if err != nil {
		if _, ok := err.(repositoryContracts.TenantNotFoundError); ok {
			return nil, contracts.NewTenantNotFoundError(request.TenantID)
		}

		return nil, contracts.NewUnknownError(err.Error())
	}

	return &contracts.ReadTenantResponse{
		Tenant: response.Tenant,
	}, nil
}

// UpdateTenant update an existing tenant
// context: Mandatory The reference to the context
// request: Mandatory. The request to update an esiting tenant
// Returns either the result of updateing an exiting tenant or error if something goes wrong.
func (service *TenantService) UpdateTenant(
	ctx context.Context,
	request *contracts.UpdateTenantRequest) (*contracts.UpdateTenantResponse, error) {
	if ctx == nil {
		return nil, commonErrors.NewArgumentError("ctx", "ctx is required")
	}

	if request == nil {
		return nil, commonErrors.NewArgumentError("request", "request is required")
	}

	if err := request.Validate(); err != nil {
		return nil, commonErrors.NewArgumentError("request", err.Error())
	}

	_, err := service.repositoryService.UpdateTenant(ctx, &repositoryContracts.UpdateTenantRequest{
		TenantID: request.TenantID,
		Tenant:   request.Tenant,
	})

	if err != nil {
		if _, ok := err.(repositoryContracts.TenantNotFoundError); ok {
			return nil, contracts.NewTenantNotFoundError(request.TenantID)
		}

		return nil, contracts.NewUnknownError(err.Error())
	}

	return &contracts.UpdateTenantResponse{}, nil
}

// DeleteTenant delete an existing tenant
// context: Mandatory The reference to the context
// request: Mandatory. The request to delete an esiting tenant
// Returns either the result of deleting an exiting tenant or error if something goes wrong.
func (service *TenantService) DeleteTenant(
	ctx context.Context,
	request *contracts.DeleteTenantRequest) (*contracts.DeleteTenantResponse, error) {
	if ctx == nil {
		return nil, commonErrors.NewArgumentError("ctx", "ctx is required")
	}

	if request == nil {
		return nil, commonErrors.NewArgumentError("request", "request is required")
	}

	if err := request.Validate(); err != nil {
		return nil, commonErrors.NewArgumentError("request", err.Error())
	}

	_, err := service.repositoryService.DeleteTenant(ctx, &repositoryContracts.DeleteTenantRequest{
		TenantID: request.TenantID,
	})

	if err != nil {
		if _, ok := err.(repositoryContracts.TenantNotFoundError); ok {
			return nil, contracts.NewTenantNotFoundError(request.TenantID)
		}

		return nil, contracts.NewUnknownError(err.Error())
	}

	return &contracts.DeleteTenantResponse{}, nil
}
