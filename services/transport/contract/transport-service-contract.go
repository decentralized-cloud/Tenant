// Package contract defines the different transport contracts
package contract

// TransportServiceContract declares the methods to be implemented by the transport service
type TransportServiceContract interface {
	// Start the transport service.
	// Returns error if something goes wrong.
	Start() error

	// Stop the transport service.
	// Returns error if something goes wrong.
	Stop() error
}
