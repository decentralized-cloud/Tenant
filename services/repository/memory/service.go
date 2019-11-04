// Package memory implements im-memory repository services
package memory

import (
	"context"
	"sort"

	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/repository"
	"github.com/lucsky/cuid"
	"github.com/micro-business/go-core/common"
	"github.com/thoas/go-funk"
)

type repositoryService struct {
	tenants map[string]models.Tenant
}

// NewRepositoryService creates new instance of the RepositoryService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewRepositoryService() (repository.RepositoryContract, error) {
	return &repositoryService{
		tenants: make(map[string]models.Tenant),
	}, nil
}

// CreateTenant creates a new tenant.
// ctx: Optional The reference to the context
// request: Mandatory. The request to create a new tenant
// Returns either the result of creating new tenant or error if something goes wrong.
func (service *repositoryService) CreateTenant(
	ctx context.Context,
	request *repository.CreateTenantRequest) (*repository.CreateTenantResponse, error) {
	tenantID := cuid.New()
	service.tenants[tenantID] = request.Tenant

	return &repository.CreateTenantResponse{
		TenantID: tenantID,
		Tenant:   request.Tenant,
	}, nil
}

// ReadTenant read an existing tenant
// ctx: Optional The reference to the context
// request: Mandatory. The request to read an existing tenant
// Returns either the result of reading an exiting tenant or error if something goes wrong.
func (service *repositoryService) ReadTenant(
	ctx context.Context,
	request *repository.ReadTenantRequest) (*repository.ReadTenantResponse, error) {
	tenant, ok := service.tenants[request.TenantID]
	if !ok {
		return nil, repository.NewTenantNotFoundError(request.TenantID)
	}

	return &repository.ReadTenantResponse{Tenant: tenant}, nil
}

// UpdateTenant update an existing tenant
// ctx: Optional The reference to the context
// request: Mandatory. The request to update an existing tenant
// Returns either the result of updateing an exiting tenant or error if something goes wr<F2>ong.
func (service *repositoryService) UpdateTenant(
	ctx context.Context,
	request *repository.UpdateTenantRequest) (*repository.UpdateTenantResponse, error) {
	_, ok := service.tenants[request.TenantID]
	if !ok {
		return nil, repository.NewTenantNotFoundError(request.TenantID)
	}

	service.tenants[request.TenantID] = request.Tenant

	return &repository.UpdateTenantResponse{
		Tenant: request.Tenant,
	}, nil
}

// DeleteTenant delete an existing tenant
// ctx: Optional The reference to the context
// request: Mandatory. The request to delete an existing tenant
// Returns either the result of deleting an exiting tenant or error if something goes wrong.
func (service *repositoryService) DeleteTenant(
	ctx context.Context,
	request *repository.DeleteTenantRequest) (*repository.DeleteTenantResponse, error) {
	_, ok := service.tenants[request.TenantID]
	if !ok {
		return nil, repository.NewTenantNotFoundError(request.TenantID)
	}

	delete(service.tenants, request.TenantID)

	return &repository.DeleteTenantResponse{}, nil
}

// Search returns the list of tenants that matched the criteria
// ctx: Mandatory The reference to the context
// request: Mandatory. The request contains the search criteria
// Returns the list of tenants that matched the criteria
func (service *repositoryService) Search(
	ctx context.Context,
	request *repository.SearchRequest) (*repository.SearchResponse, error) {
	response := &repository.SearchResponse{
		HasPreviousPage: false,
		HasNextPage:     false,
	}

	tenantsWithCursor := funk.Map(service.tenants, func(tenantID string, tenant models.Tenant) models.TenantWithCursor {
		return models.TenantWithCursor{
			TenantID: tenantID,
			Tenant:   tenant,
			Cursor:   "Not implemented",
		}
	})

	if len(request.TenantIDs) > 0 {
		tenantsWithCursor = funk.Filter(tenantsWithCursor, func(tenantWithCursor models.TenantWithCursor) bool {
			return funk.Contains(request.TenantIDs, tenantWithCursor.TenantID)
		})
	}

	response.Tenants = tenantsWithCursor.([]models.TenantWithCursor)

	// Default sorting is acsending if not provided, aslo as we only have one field currenrly stored againsst a tenant, we are ignroing the provided field name to sort on
	sortingDirection := common.Ascending
	if len(request.SortingOptions) > 0 {
		sortingDirection = request.SortingOptions[0].Direction
	}

	sort.Slice(response.Tenants, func(i, j int) bool {
		if sortingDirection == common.Ascending {
			return response.Tenants[i].Tenant.Name < response.Tenants[j].Tenant.Name
		}

		return response.Tenants[i].Tenant.Name > response.Tenants[j].Tenant.Name
	})

	return response, nil
}
