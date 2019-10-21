// Package configuration implements configuration service required by the tenant service
package configuration

// ConfigurationContract declares the service that provides configuration required by different Tenat modules
type ConfigurationContract interface {
	// GetGrpcHost retrieves gRPC host name
	// Returns the gRPC host name or error if something goes wrong
	GetGrpcHost() (string, error)

	// GetGrpcPort retrieves gRPC port number
	// Returns the gRPC port number or error if something goes wrong
	GetGrpcPort() (int, error)

	// GetHttpsHost retrieves HTTPS host name
	// Returns the HTTPS host name or error if something goes wrong
	GetHttpsHost() (string, error)

	// GetHttpsPort retrieves HTTPS port number
	// Returns the HTTPS port number or error if something goes wrong
	GetHttpsPort() (int, error)

	// GetDbConnectionString retrieves database connection string
	// Returns the database connection string or error if something goes wrong
	GetDbConnectionString() (string, error)
}
