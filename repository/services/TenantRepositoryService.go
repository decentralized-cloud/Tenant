// Package services implements the different tenant repository services
package services

import (
	"context"

	"github.com/decentralized-cloud/Tenant/repository/contracts"
	"github.com/lucsky/cuid"
)

// TenantRepositoryService implements the repository service that create new tenant, read, update and delete existing tenants.
type TenantRepositoryService struct {
}

// NewTenantRepositoryService creates new instance of the TenantRepositoryService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewTenantRepositoryService() (contracts.TenantRepositoryServiceContract, error) {
	return &TenantRepositoryService{}, nil
}

// CreateTenant creates a new tenant.
// context: Optional The reference to the context
// request: Mandatory. The request to create a new tenant
// Returns either the result of creating new tenant or error if something goes wrong.
func (service *TenantRepositoryService) CreateTenant(
	ctx context.Context,
	request *contracts.CreateTenantRequest) (*contracts.CreateTenantResponse, error) {
	return &contracts.CreateTenantResponse{
		TenantID: cuid.New(),
	}, nil
}

// ReadTenant read an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to read an esiting tenant
// Returns either the result of reading an exiting tenant or error if something goes wrong.
func (service *TenantRepositoryService) ReadTenant(
	ctx context.Context,
	request *contracts.ReadTenantRequest) (*contracts.ReadTenantResponse, error) {
	return &contracts.ReadTenantResponse{}, nil
}

// UpdateTenant update an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to update an esiting tenant
// Returns either the result of updateing an exiting tenant or error if something goes wrong.
func (service *TenantRepositoryService) UpdateTenant(
	ctx context.Context,
	request *contracts.UpdateTenantRequest) (*contracts.UpdateTenantResponse, error) {
	return &contracts.UpdateTenantResponse{}, nil
}

// DeleteTenant delete an existing tenant
// context: Optional The reference to the context
// request: Mandatory. The request to delete an esiting tenant
// Returns either the result of deleting an exiting tenant or error if something goes wrong.
func (service *TenantRepositoryService) DeleteTenant(
	ctx context.Context,
	request *contracts.DeleteTenantRequest) (*contracts.DeleteTenantResponse, error) {
	return &contracts.DeleteTenantResponse{}, nil
}
