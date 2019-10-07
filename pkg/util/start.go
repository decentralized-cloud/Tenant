// Package util implements different utilities required by the tenant service
package util

import (
	"log"
	"os"
	"os/signal"

	"github.com/decentralized-cloud/tenant/services/business"
	"github.com/decentralized-cloud/tenant/services/configuration"
	"github.com/decentralized-cloud/tenant/services/endpoint"
	"github.com/decentralized-cloud/tenant/services/repository"
	grpctransport "github.com/decentralized-cloud/tenant/services/transport/grpc"
	"go.uber.org/zap"
)

var configurationService configuration.ConfigurationContract
var endpointCreatorService endpoint.EndpointCreatorContract

// StartService setups all dependecies required to start the tenant service and
// start the service
func StartService() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	defer logger.Sync()

	if err = setupDependencies(); err != nil {
		logger.Fatal("Failed to setup dependecies", zap.Error(err))
	}

	grpcTransportService, err := grpctransport.NewTransportService(
		logger,
		configurationService,
		endpointCreatorService)
	if err != nil {
		logger.Fatal("Failed to create gRPC transport service", zap.Error(err))
	}

	go grpcTransportService.Start()

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan
		logger.Info("Received an interrupt, stopping services...")

		if err := grpcTransportService.Stop(); err != nil {
			logger.Error("Failed to stop GRPC transport service", zap.Error(err))
		}

		close(cleanupDone)
	}()
	<-cleanupDone
}

func setupDependencies() (err error) {
	if configurationService, err = configuration.NewEnvConfigurationService(); err != nil {
		return
	}

	repositoryService, err := repository.NewInMemoryRepositoryService()

	if err != nil {
		return
	}

	businessService, err := business.NewBusinessService(repositoryService)

	if err != nil {
		return err
	}

	if endpointCreatorService, err = endpoint.NewEndpointCreatorService(businessService); err != nil {
		return
	}

	return
}
