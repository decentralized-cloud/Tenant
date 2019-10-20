// Package configuration implements configuration service required by the tenant service
package configuration

import (
	"os"
	"strconv"
	"strings"
)

type envConfigurationService struct {
}

// NewEnvConfigurationService creates new instance of the EnvConfigurationService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewEnvConfigurationService() (ConfigurationContract, error) {
	return &envConfigurationService{}, nil
}

// GetHost retrieves host name
// Returns the host name or error if something goes wrong
func (service *envConfigurationService) GetHost() (string, error) {
	return os.Getenv("HOST"), nil
}

// GetPort retrieves port number
// Returns the port number or error if something goes wrong
func (service *envConfigurationService) GetPort() (int, error) {
	portNumberString := os.Getenv("PORT")
	if strings.Trim(portNumberString, " ") == "" {
		return 0, NewUnknownError("PORT is required")
	}

	portNumber, err := strconv.Atoi(portNumberString)
	if err != nil {
		return 0, NewUnknownErrorWithError("Failed to convert PORT to integer", err)
	}

	return portNumber, nil
}
