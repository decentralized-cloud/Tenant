// Package service defines configuration service implementation
package service

import (
	"os"
	"strconv"

	"github.com/decentralized-cloud/tenant/services/configuration/contract"
)

type configurationService struct {
}

// NewConfigurationService creates new instance of the ConfigurationService, setting up all dependencies and returns the instance
// Returns the new service or error if something goes wrong
func NewConfigurationService() (contract.ConfigurationServiceContract, error) {
	return &configurationService{}, nil
}

// GetGRPCPort retrieves gRPC port number from environment variable
// Returns the gRPC port number or error if something goes wrong
func (service *configurationService) GetGRPCPort() (int, error) {
	portNumberString := os.Getenv("PORT")
	portNumber, err := strconv.Atoi(portNumberString)

	if err != nil {
		return 0, contract.NewUnknownError(err.Error())
	}

	return portNumber, nil
}

// GetGRPCHostName retrieves gRPC host name from environment variable
// Returns the gRPC host name or error if something goes wrong
func (service *configurationService) GetGRPCHost() (string, error) {
	hostName := os.Getenv("HOST")

	return hostName, nil
}
