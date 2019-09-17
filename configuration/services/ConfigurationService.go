// Package services defines configuration service implementation
package services

import (
	"os"
	"strconv"

	"github.com/decentralized-cloud/Tenant/configuration/contracts"
)

// ConfigurationService implements the service that provides environment variables
type ConfigurationService struct {
}

// NewConfigurationService creates new instance of the ConfigurationService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewConfigurationService() (configurationService *ConfigurationService, err error) {
	return configurationService, nil
}

// GetPort retrieves port number from environment variable
// Returns the port number
func (service *ConfigurationService) GetPort() (*int, error) {
	//Todo : what is port number environment variable name
	portNumberString := os.Getenv("PORT")
	portNumber, err := strconv.Atoi(portNumberString)
	if err != nil {
		return nil, contracts.NewUnknownError(err.Error())
	}
	return &portNumber, nil
}

// GetHostName retrieves host name from environment variable
// Returns the port number
func (service *ConfigurationService) GetHostName() (*string, error) {
	hostName, err := os.GetHostName()
	if err != nil {
		return nil, contracts.NewUnknownError(err.Error())
	}
	return hostName, nil

}
