// Package endpoint implements different endpoint services required by the project service
package endpoint

import "github.com/go-kit/kit/endpoint"

// EndpointCreatorContract declares the contract that creates endpoints to create new project,
// read, update and delete existing projects.
type EndpointCreatorContract interface {
	// CreateProjectEndpoint creates Create Project endpoint
	// Returns the Create Project endpoint
	CreateProjectEndpoint() endpoint.Endpoint

	// ReadProjectEndpoint creates Read Project endpoint
	// Returns the Read Project endpoint
	ReadProjectEndpoint() endpoint.Endpoint

	// UpdateProjectEndpoint creates Update Project endpoint
	// Returns the Update Project endpoint
	UpdateProjectEndpoint() endpoint.Endpoint

	// DeleteProjectEndpoint creates Delete Project endpoint
	// Returns the Delete Project endpoint
	DeleteProjectEndpoint() endpoint.Endpoint

	// SearchEndpoint creates Search Project endpoint
	// Returns the Search Project endpoint
	SearchEndpoint() endpoint.Endpoint
}
