// Package models defines the different object models used in Tenant
package models

// Tenant defines the tenant object
type Tenant struct {
	Name string
}

// TenantWithCursor implements the pair of the tenant with a cursor that determines the
// location of the tennat in the repository.
type TenantWithCursor struct {
	TenantID string
	Tenant   Tenant
	Cursor   string
}
