// Package business implements different business services required by the tenant service
package business

import (
	"github.com/decentralized-cloud/tenant/models"
	"github.com/micro-business/go-core/common"
)

// CreateTenantRequest contains the request to create a new tenant
type CreateTenantRequest struct {
	Tenant models.Tenant
}

// CreateTenantResponse contains the result of creating a new tenant
type CreateTenantResponse struct {
	TenantID string
	Err      error
	Tenant   models.Tenant
}

// ReadTenantRequest contains the request to read an existing tenant
type ReadTenantRequest struct {
	TenantID string
}

// ReadTenantResponse contains the result of reading an existing tenant
type ReadTenantResponse struct {
	Tenant models.Tenant
	Err    error
}

// UpdateTenantRequest contains the request to update an existing tenant
type UpdateTenantRequest struct {
	TenantID string
	Tenant   models.Tenant
}

// UpdateTenantResponse contains the result of updating an existing tenant
type UpdateTenantResponse struct {
	Err    error
	Tenant models.Tenant
}

// DeleteTenantRequest contains the request to delete an existing tenant
type DeleteTenantRequest struct {
	TenantID string
}

// DeleteTenantResponse contains the result of deleting an existing tenant
type DeleteTenantResponse struct {
	Err error
}

// SearchRequest contains the filter criteria to look for existing tenants
type SearchRequest struct {
	Pagination     common.Pagination
	SortingOptions []common.SortingOptionPair
	TenantIDs      []string
}

// SearchResponse contains the list of the tenants that matched the result
type SearchResponse struct {
	Err             error
	HasPreviousPage bool
	HasNextPage     bool
	Tenants         []models.TenantWithCursor
}
