// Package memory implements im-memory repository services
package memory

import (
	"context"

	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/repository"
	"github.com/lucsky/cuid"
)

var tenants map[string]models.Tenant

type repositoryService struct {
}

func init() {
	tenants = make(map[string]models.Tenant)
}

// NewRepositoryService creates new instance of the RepositoryService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewRepositoryService() (repository.RepositoryContract, error) {
	return &repositoryService{}, nil
}

// CreateTenant creates a new tenant.
// context: Optional The reference to the context
// request: Mandatory. The request to create a new tenant
// Returns either the result of creating new tenant or error if something goes wrong.
func (service *repositoryService) CreateTenant(
	ctx context.Context,
	request *repository.CreateTenantRequest) (*repository.CreateTenantResponse, error) {

	tenantID := cuid.New()
	tenants[tenantID] = request.Tenant

	return &repository.CreateTenantResponse{
		TenantID: tenantID,
	}, nil
}

// ReadTenant read an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to read an existing tenant
// Returns either the result of reading an exiting tenant or error if something goes wrong.
func (service *repositoryService) ReadTenant(
	ctx context.Context,
	request *repository.ReadTenantRequest) (*repository.ReadTenantResponse, error) {

	tenant, ok := tenants[request.TenantID]
	if !ok {
		return nil, repository.NewTenantNotFoundError(request.TenantID)
	}

	return &repository.ReadTenantResponse{Tenant: tenant}, nil
}

// UpdateTenant update an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to update an existing tenant
// Returns either the result of updateing an exiting tenant or error if something goes wrong.
func (service *repositoryService) UpdateTenant(
	ctx context.Context,
	request *repository.UpdateTenantRequest) (*repository.UpdateTenantResponse, error) {

	_, ok := tenants[request.TenantID]
	if !ok {
		return nil, repository.NewTenantNotFoundError(request.TenantID)
	}

	tenants[request.TenantID] = request.Tenant

	return &repository.UpdateTenantResponse{}, nil
}

// DeleteTenant delete an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to delete an existing tenant
// Returns either the result of deleting an exiting tenant or error if something goes wrong.
func (service *repositoryService) DeleteTenant(
	ctx context.Context,
	request *repository.DeleteTenantRequest) (*repository.DeleteTenantResponse, error) {

	_, ok := tenants[request.TenantID]
	if !ok {
		return nil, repository.NewTenantNotFoundError(request.TenantID)
	}

	delete(tenants, request.TenantID)

	return &repository.DeleteTenantResponse{}, nil
}
