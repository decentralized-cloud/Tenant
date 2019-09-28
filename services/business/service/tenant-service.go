// Package service implements the different tenant business services
package service

import (
	"context"

	"github.com/decentralized-cloud/tenant/services/business/contract"
	repositoryContract "github.com/decentralized-cloud/tenant/services/repository/contract"
	commonErrors "github.com/micro-business/go-core/system/errors"
)

// TenantService implements the service that create new tenant, read, update and delete existing tenants.
type TenantService struct {
	repositoryService repositoryContract.TenantRepositoryServiceContract
}

// NewTenantService creates new instance of the TenantService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewTenantService(
	repositoryService repositoryContract.TenantRepositoryServiceContract) (contract.TenantServiceContract, error) {
	if repositoryService == nil {
		return nil, commonErrors.NewArgumentError("repositoryService", "repositoryService is required")
	}

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
	request *contract.CreateTenantRequest) (*contract.CreateTenantResponse, error) {
	if ctx == nil {
		return &contract.CreateTenantResponse{
			Err: commonErrors.NewArgumentError("ctx", "ctx is required"),
		}, nil
	}

	if request == nil {
		return &contract.CreateTenantResponse{
			Err: commonErrors.NewArgumentError("request", "request is required"),
		}, nil
	}

	if err := request.Validate(); err != nil {
		return &contract.CreateTenantResponse{
			Err: commonErrors.NewArgumentError("request", err.Error()),
		}, nil
	}

	response, err := service.repositoryService.CreateTenant(ctx, &repositoryContract.CreateTenantRequest{
		Tenant: request.Tenant,
	})

	if err != nil {
		if _, ok := err.(repositoryContract.TenantAlreadyExistsError); ok {
			return &contract.CreateTenantResponse{
				Err: contract.NewTenantAlreadyExistsError(),
			}, nil
		}

		return &contract.CreateTenantResponse{
			Err: contract.NewUnknownError(err.Error()),
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
func (service *TenantService) ReadTenant(
	ctx context.Context,
	request *contract.ReadTenantRequest) (*contract.ReadTenantResponse, error) {
	if ctx == nil {
		return &contract.ReadTenantResponse{
			Err: commonErrors.NewArgumentError("ctx", "ctx is required"),
		}, nil
	}

	if request == nil {
		return &contract.ReadTenantResponse{
			Err: commonErrors.NewArgumentError("request", "request is required"),
		}, nil
	}

	if err := request.Validate(); err != nil {
		return &contract.ReadTenantResponse{
			Err: commonErrors.NewArgumentError("request", err.Error()),
		}, nil
	}

	response, err := service.repositoryService.ReadTenant(ctx, &repositoryContract.ReadTenantRequest{
		TenantID: request.TenantID,
	})

	if err != nil {
		if _, ok := err.(repositoryContract.TenantNotFoundError); ok {
			return &contract.ReadTenantResponse{
				Err: contract.NewTenantNotFoundError(request.TenantID),
			}, nil
		}

		return &contract.ReadTenantResponse{
			Err: contract.NewUnknownError(err.Error()),
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
func (service *TenantService) UpdateTenant(
	ctx context.Context,
	request *contract.UpdateTenantRequest) (*contract.UpdateTenantResponse, error) {
	if ctx == nil {
		return &contract.UpdateTenantResponse{
			Err: commonErrors.NewArgumentError("ctx", "ctx is required"),
		}, nil
	}

	if request == nil {
		return &contract.UpdateTenantResponse{
			Err: commonErrors.NewArgumentError("request", "request is required"),
		}, nil
	}

	if err := request.Validate(); err != nil {
		return &contract.UpdateTenantResponse{
			Err: commonErrors.NewArgumentError("request", err.Error()),
		}, nil
	}

	_, err := service.repositoryService.UpdateTenant(ctx, &repositoryContract.UpdateTenantRequest{
		TenantID: request.TenantID,
		Tenant:   request.Tenant,
	})

	if err != nil {
		if _, ok := err.(repositoryContract.TenantNotFoundError); ok {
			return &contract.UpdateTenantResponse{
				Err: contract.NewTenantNotFoundError(request.TenantID),
			}, nil
		}

		return &contract.UpdateTenantResponse{
			Err: contract.NewUnknownError(err.Error()),
		}, nil
	}

	return &contract.UpdateTenantResponse{}, nil
}

// DeleteTenant delete an existing tenant
// context: Mandatory The reference to the context
// request: Mandatory. The request to delete an existing tenant
// Returns either the result of deleting an existing tenant or error if something goes wrong.
func (service *TenantService) DeleteTenant(
	ctx context.Context,
	request *contract.DeleteTenantRequest) (*contract.DeleteTenantResponse, error) {
	if ctx == nil {
		return &contract.DeleteTenantResponse{
			Err: commonErrors.NewArgumentError("ctx", "ctx is required"),
		}, nil
	}

	if request == nil {
		return &contract.DeleteTenantResponse{
			Err: commonErrors.NewArgumentError("request", "request is required"),
		}, nil
	}

	if err := request.Validate(); err != nil {
		return &contract.DeleteTenantResponse{
			Err: commonErrors.NewArgumentError("request", err.Error()),
		}, nil
	}

	_, err := service.repositoryService.DeleteTenant(ctx, &repositoryContract.DeleteTenantRequest{
		TenantID: request.TenantID,
	})

	if err != nil {
		if _, ok := err.(repositoryContract.TenantNotFoundError); ok {
			return &contract.DeleteTenantResponse{
				Err: contract.NewTenantNotFoundError(request.TenantID),
			}, nil
		}

		return &contract.DeleteTenantResponse{
			Err: contract.NewUnknownError(err.Error()),
		}, nil
	}

	return &contract.DeleteTenantResponse{}, nil
}
