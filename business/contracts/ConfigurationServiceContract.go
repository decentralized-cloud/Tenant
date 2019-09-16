// Package contracts defines the configuration business contracts
package contracts

// ConfigurationServiceContract declares the service that provides environment variables
type ConfigurationServiceContract interface {
	// GetPort retrieves port number from environment variable
	GetPort() (int, error)
	// GetHost retrieves host name from environment variable
	GetHost() (string, error)
}
