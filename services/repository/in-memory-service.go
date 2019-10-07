// Package repository implements different repository services required by the tenant service
package repository

import (
	"context"

	"github.com/decentralized-cloud/tenant/models"
	"github.com/lucsky/cuid"
)

var tenants map[string]models.Tenant

type inMemoryRepositoryService struct {
}

func init() {
	tenants = make(map[string]models.Tenant)
}

// NewInMemoryRepositoryService creates new instance of the InMemoryRepositoryService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewInMemoryRepositoryService() (RepositoryContract, error) {
	return &inMemoryRepositoryService{}, nil
}

// CreateTenant creates a new tenant.
// context: Optional The reference to the context
// request: Mandatory. The request to create a new tenant
// Returns either the result of creating new tenant or error if something goes wrong.
func (service *inMemoryRepositoryService) CreateTenant(
	ctx context.Context,
	request *CreateTenantRequest) (*CreateTenantResponse, error) {

	tenantID := cuid.New()
	tenants[tenantID] = request.Tenant

	return &CreateTenantResponse{
		TenantID: tenantID,
	}, nil
}

// ReadTenant read an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to read an existing tenant
// Returns either the result of reading an exiting tenant or error if something goes wrong.
func (service *inMemoryRepositoryService) ReadTenant(
	ctx context.Context,
	request *ReadTenantRequest) (*ReadTenantResponse, error) {

	tenant, ok := tenants[request.TenantID]
	if !ok {
		return nil, NewTenantNotFoundError(request.TenantID)
	}

	return &ReadTenantResponse{Tenant: tenant}, nil
}

// UpdateTenant update an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to update an existing tenant
// Returns either the result of updateing an exiting tenant or error if something goes wrong.
func (service *inMemoryRepositoryService) UpdateTenant(
	ctx context.Context,
	request *UpdateTenantRequest) (*UpdateTenantResponse, error) {

	_, ok := tenants[request.TenantID]
	if !ok {
		return nil, NewTenantNotFoundError(request.TenantID)
	}

	tenants[request.TenantID] = request.Tenant

	return &UpdateTenantResponse{}, nil
}

// DeleteTenant delete an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to delete an existing tenant
// Returns either the result of deleting an exiting tenant or error if something goes wrong.
func (service *inMemoryRepositoryService) DeleteTenant(
	ctx context.Context,
	request *DeleteTenantRequest) (*DeleteTenantResponse, error) {

	_, ok := tenants[request.TenantID]
	if !ok {
		return nil, NewTenantNotFoundError(request.TenantID)
	}

	delete(tenants, request.TenantID)

	return &DeleteTenantResponse{}, nil
}
