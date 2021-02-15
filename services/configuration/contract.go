// Package configuration implements configuration service required by the project service
package configuration

// ConfigurationContract declares the service that provides configuration required by different Tenat modules
type ConfigurationContract interface {
	// GetGrpcHost retrieves the gRPC host name
	// Returns the gRPC host name or error if something goes wrong
	GetGrpcHost() (string, error)

	// GetGrpcPort retrieves the gRPC port number
	// Returns the gRPC port number or error if something goes wrong
	GetGrpcPort() (int, error)

	// GetHttpHost retrieves the HTTP host name
	// Returns the HTTP host name or error if something goes wrong
	GetHttpHost() (string, error)

	// GetHttpPort retrieves the HTTP port number
	// Returns the HTTP port number or error if something goes wrong
	GetHttpPort() (int, error)

	// GetDatabaseConnectionString retrieves the database connection string
	// Returns the database connection string or error if something goes wrong
	GetDatabaseConnectionString() (string, error)

	// GetDatabaseName retrieves the database name
	// Returns the database name or error if something goes wrong
	GetDatabaseName() (string, error)

	// GetDatabaseCollectionName retrieves the database collection name
	// Returns the database collection name or error if something goes wrong
	GetDatabaseCollectionName() (string, error)

	// GetJwksURL retrieves the JWKS URL
	// Returns the JWKS URL or error if something goes wrong
	GetJwksURL() (string, error)
}
