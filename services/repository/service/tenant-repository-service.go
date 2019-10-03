// Package service implements the different tenant repository services
package service

import (
	"context"

	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/repository/contract"

	"github.com/lucsky/cuid"
)

var tenants map[string]models.Tenant

type tenantRepositoryService struct {
}

func init() {
	tenants = make(map[string]models.Tenant)
}

// NewTenantRepositoryService creates new instance of the TenantRepositoryService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewTenantRepositoryService() (contract.TenantRepositoryServiceContract, error) {
	return &tenantRepositoryService{}, nil
}

// CreateTenant creates a new tenant.
// context: Optional The reference to the context
// request: Mandatory. The request to create a new tenant
// Returns either the result of creating new tenant or error if something goes wrong.
func (service *tenantRepositoryService) CreateTenant(
	ctx context.Context,
	request *contract.CreateTenantRequest) (*contract.CreateTenantResponse, error) {

	tenantID := cuid.New()
	tenants[tenantID] = request.Tenant

	return &contract.CreateTenantResponse{
		TenantID: tenantID,
	}, nil
}

// ReadTenant read an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to read an existing tenant
// Returns either the result of reading an exiting tenant or error if something goes wrong.
func (service *tenantRepositoryService) ReadTenant(
	ctx context.Context,
	request *contract.ReadTenantRequest) (*contract.ReadTenantResponse, error) {

	tenant, ok := tenants[request.TenantID]
	if !ok {
		return nil, contract.NewTenantNotFoundError(request.TenantID)
	}

	return &contract.ReadTenantResponse{Tenant: tenant}, nil
}

// UpdateTenant update an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to update an existing tenant
// Returns either the result of updateing an exiting tenant or error if something goes wrong.
func (service *tenantRepositoryService) UpdateTenant(
	ctx context.Context,
	request *contract.UpdateTenantRequest) (*contract.UpdateTenantResponse, error) {

	_, ok := tenants[request.TenantID]
	if !ok {
		return nil, contract.NewTenantNotFoundError(request.TenantID)
	}

	tenants[request.TenantID] = request.Tenant

	return &contract.UpdateTenantResponse{}, nil
}

// DeleteTenant delete an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to delete an existing tenant
// Returns either the result of deleting an exiting tenant or error if something goes wrong.
func (service *tenantRepositoryService) DeleteTenant(
	ctx context.Context,
	request *contract.DeleteTenantRequest) (*contract.DeleteTenantResponse, error) {

	_, ok := tenants[request.TenantID]
	if !ok {
		return nil, contract.NewTenantNotFoundError(request.TenantID)
	}

	delete(tenants, request.TenantID)

	return &contract.DeleteTenantResponse{}, nil
}
