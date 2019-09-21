// Package contract defines the different tenant repository contracts
package contract

import "context"

// TenantRepositoryServiceContract declares the repository service that can create new tenant, read, update
// and delete existing tenants.
type TenantRepositoryServiceContract interface {
	// CreateTenant creates a new tenant.
	// context: Mandatory The reference to the context
	// request: Mandatory. The request to create a new tenant
	// Returns either the result of creating new tenant or error if something goes wrong.
	CreateTenant(
		ctx context.Context,
		request *CreateTenantRequest) (*CreateTenantResponse, error)

	// ReadTenant read an existing tenant
	// context: Mandatory The reference to the context
	// request: Mandatory. The request to read an esiting tenant
	// Returns either the result of reading an exiting tenant or error if something goes wrong.
	ReadTenant(
		ctx context.Context,
		request *ReadTenantRequest) (*ReadTenantResponse, error)

	// UpdateTenant update an existing tenant
	// context: Mandatory The reference to the context
	// request: Mandatory. The request to update an esiting tenant
	// Returns either the result of updateing an exiting tenant or error if something goes wrong.
	UpdateTenant(
		ctx context.Context,
		request *UpdateTenantRequest) (*UpdateTenantResponse, error)

	// DeleteTenant delete an existing tenant
	// context: Mandatory The reference to the context
	// request: Mandatory. The request to delete an esiting tenant
	// Returns either the result of deleting an exiting tenant or error if something goes wrong.
	DeleteTenant(
		ctx context.Context,
		request *DeleteTenantRequest) (*DeleteTenantResponse, error)
}
