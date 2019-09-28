//Package contract defines configuration service contracts
package contract

// ConfigurationServiceContract declares the service that provides configuration required by different Tenat modules
type ConfigurationServiceContract interface {
	// GetGRPCPort retrieves gRPC port number from environment variable
	// Returns the gRPC port number or error if something goes wrong
	GetGRPCPort() (int, error)

	// GetGRPCHostName retrieves gRPC host name from environment variable
	// Returns the gRPC host name or error if something goes wrong
	GetGRPCHost() (string, error)
}
