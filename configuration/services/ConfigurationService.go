// Package services defines configuration service implementation
package services

import (
	"os"
	"strconv"

	"github.com/decentralized-cloud/Tenant/configuration/contracts"
)

// ConfigurationService implements the service that provides configuration required by different Tenat modules
type ConfigurationService struct {
}

// NewConfigurationService creates new instance of the ConfigurationService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewConfigurationService() (contracts.ConfigurationServiceContract, error) {
	return &ConfigurationService{}, nil
}

// GetPort retrieves port number from environment variable
// Returns the port number or error if something goes wrong
func (service *ConfigurationService) GetPort() (int, error) {
	portNumberString := os.Getenv("PORT")
	portNumber, err := strconv.Atoi(portNumberString)

	if err != nil {
		return 0, contracts.NewUnknownError(err.Error())
	}

	return portNumber, nil
}

// GetHost retrieves host from environment variable
// Returns the host or error if something goes wrong
func (service *ConfigurationService) GetHost() (string, error) {
	hostName := os.Getenv("HOST")

	return hostName, nil
}
