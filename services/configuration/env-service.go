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

// GetGrpcHost retrieves gRPC host name
// Returns the gRPC host name or error if something goes wrong
func (service *envConfigurationService) GetGrpcHost() (string, error) {
	return os.Getenv("GRPC_HOST"), nil
}

// GetGrpcPort retrieves gRPC port number
// Returns the gRPC port number or error if something goes wrong
func (service *envConfigurationService) GetGrpcPort() (int, error) {
	portNumberString := os.Getenv("GRPC_PORT")
	if strings.Trim(portNumberString, " ") == "" {
		return 0, NewUnknownError("GRPC_PORT is required")
	}

	portNumber, err := strconv.Atoi(portNumberString)
	if err != nil {
		return 0, NewUnknownErrorWithError("Failed to convert GRPC_PORT to integer", err)
	}

	return portNumber, nil
}

// GetHttpsHost retrieves HTTPS host name
// Returns the HTTPS host name or error if something goes wrong
func (service *envConfigurationService) GetHttpsHost() (string, error) {
	return os.Getenv("HTTPS_HOST"), nil
}

// GetHttpsPort retrieves HTTPS port number
// Returns the HTTPS port number or error if something goes wrong
func (service *envConfigurationService) GetHttpsPort() (int, error) {
	portNumberString := os.Getenv("HTTPS_PORT")
	if strings.Trim(portNumberString, " ") == "" {
		return 0, NewUnknownError("HTTPS_PORT is required")
	}

	portNumber, err := strconv.Atoi(portNumberString)
	if err != nil {
		return 0, NewUnknownErrorWithError("Failed to convert HTTPS_PORT to integer", err)
	}

	return portNumber, nil
}

// GetDbConnectionString retrieves database connection string
// Returns the database connection string or error if something goes wrong
func (service *envConfigurationService) GetDbConnectionString() (string, error) {
	connectionString := os.Getenv("DB_CONNECTION_STRING")

	if strings.Trim(connectionString, " ") == "" {
		return "", NewUnknownError("DB_CONNECTION_STRING is required")
	}

	return connectionString, nil
}
