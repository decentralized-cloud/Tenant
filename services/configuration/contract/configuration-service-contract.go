//Package contract defines configuration service contracts
package contract

// ConfigurationServiceContract declares the service that provides configuration required by different Tenat modules
type ConfigurationServiceContract interface {
	// GetPort retrieves port number from environment variable
	// Returns the port number or error if something goes wrong
	GetPort() (int, error)

	// GetHostName retrieves host name from environment variable
	// Returns the host name or error if something goes wrong
	GetHost() (string, error)
}
