// Package contracts defines the contracts that provides endpoint to be used by the transport layer
package contracts

import "github.com/go-kit/kit/endpoint"

// EndpointCreatorContract declares the contract that creates endpoints to create new tenant,
// read, update and delete existing tenants.
type EndpointCreatorContract interface {
	// CreateTenantEndpoint creates Create Tenant endpoint
	// Returns the Create Tenant endpoint
	CreateTenantEndpoint() endpoint.Endpoint

	// ReadTenantEndpoint creates Read Tenant endpoint
	// Returns the Read Tenant endpoint
	ReadTenantEndpoint() endpoint.Endpoint

	// UpdateTenantEndpoint creates Update Tenant endpoint
	// Returns the Update Tenant endpoint
	UpdateTenantEndpoint() endpoint.Endpoint

	// DeleteTenantEndpoint creates Delete Tenant endpoint
	// Returns the Delete Tenant endpoint
	DeleteTenantEndpoint() endpoint.Endpoint
}
