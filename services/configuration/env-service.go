// Package configuration implements configuration service required by the project service
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

// GetGrpcHost retrieves the gRPC host name
// Returns the gRPC host name or error if something goes wrong
func (service *envConfigurationService) GetGrpcHost() (string, error) {
	return os.Getenv("GRPC_HOST"), nil
}

// GetGrpcPort retrieves the gRPC port number
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

// GetHttpHost retrieves the HTTP host name
// Returns the HTTP host name or error if something goes wrong
func (service *envConfigurationService) GetHttpHost() (string, error) {
	return os.Getenv("HTTP_HOST"), nil
}

// GetHttpPort retrieves the HTTP port number
// Returns the HTTP port number or error if something goes wrong
func (service *envConfigurationService) GetHttpPort() (int, error) {
	portNumberString := os.Getenv("HTTP_PORT")
	if strings.Trim(portNumberString, " ") == "" {
		return 0, NewUnknownError("HTTP_PORT is required")
	}

	portNumber, err := strconv.Atoi(portNumberString)
	if err != nil {
		return 0, NewUnknownErrorWithError("Failed to convert HTTP_PORT to integer", err)
	}

	return portNumber, nil
}

// GetDatabaseConnectionString retrieves the database connection string
// Returns the database connection string or error if something goes wrong
func (service *envConfigurationService) GetDatabaseConnectionString() (string, error) {
	connectionString := os.Getenv("DATABASE_CONNECTION_STRING")

	if strings.Trim(connectionString, " ") == "" {
		return "", NewUnknownError("DATABASE_CONNECTION_STRING is required")
	}

	return connectionString, nil
}

// GetDatabaseName retrieves the database name
// Returns the database name or error if something goes wrong
func (service *envConfigurationService) GetDatabaseName() (string, error) {
	databaseName := os.Getenv("PROJECT_DATABASE_NAME")

	if strings.Trim(databaseName, " ") == "" {
		return "", NewUnknownError("PROJECT_DATABASE_NAME is required")
	}

	return databaseName, nil

}

// GetDatabaseCollectionName retrieves the database collection name
// Returns the database collection name or error if something goes wrong
func (service *envConfigurationService) GetDatabaseCollectionName() (string, error) {
	databaseCollectionName := os.Getenv("USER_DATABASE_COLLECTION_NAME")

	if strings.Trim(databaseCollectionName, " ") == "" {
		return "", NewUnknownError("USER_DATABASE_COLLECTION_NAME is required")
	}

	return databaseCollectionName, nil
}

// GetJwksURL retrieves the JWKS URL
// Returns the JWKS URL or error if something goes wrong
func (service *envConfigurationService) GetJwksURL() (string, error) {
	jwksURL := os.Getenv("JWKS_URL")

	if strings.Trim(jwksURL, " ") == "" {
		return "", NewUnknownError("JWKS_URL is required")
	}

	return jwksURL, nil
}
