// Package contracts defines the different tenant business contracts
package contracts

import "github.com/decentralized-cloud/Tenant/models"

// Request to create a new tenant
type CreateTenantRequest struct {
	Tenant models.Tenant
}

// Response contains the result of creating a new tenant
type CreateTenantResponse struct {
	TenantID string
}

// Request to read an existing tenant
type ReadTenantRequest struct {
	TenantID string
}

// Response contains the result of reading an existing tenant
type ReadTenantResponse struct {
	Tenant models.Tenant
}

// Request to update an existing tenant
type UpdateTenantRequest struct {
	TenantID string
	Tenant   models.Tenant
}

// Response contains the result of updating an existing tenant
type UpdateTenantResponse struct {
}

// Request to delete an existing tenant
type DeleteTenantRequest struct {
	TenantID string
}

// Response contains the result of deleting an existing tenant
type DeleteTenantResponse struct {
}
