// Package util implements different utilities required by the tenant service
package util

import (
	"log"
	"os"
	"os/signal"

	"github.com/decentralized-cloud/tenant/services/business"
	"github.com/decentralized-cloud/tenant/services/configuration"
	"github.com/decentralized-cloud/tenant/services/endpoint"
	"github.com/decentralized-cloud/tenant/services/repository/memory"
	"github.com/decentralized-cloud/tenant/services/transport/grpc"
	"github.com/decentralized-cloud/tenant/services/transport/https"
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

	defer func() {
		_ = logger.Sync()
	}()

	if err = setupDependencies(); err != nil {
		logger.Fatal("Failed to setup dependecies", zap.Error(err))
	}

	grpcTransportService, err := grpc.NewTransportService(
		logger,
		configurationService,
		endpointCreatorService)
	if err != nil {
		logger.Fatal("Failed to create gRPC transport service", zap.Error(err))
	}

	httpsTansportService, err := https.NewTransportService(
		logger,
		configurationService)
	if err != nil {
		logger.Fatal("Failed to create HTTPS transport service", zap.Error(err))
	}

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		if serviceErr := grpcTransportService.Start(); serviceErr != nil {
			logger.Fatal("Failed to start gRPC transport service", zap.Error(serviceErr))
		}
	}()

	go func() {
		if serviceErr := httpsTansportService.Start(); serviceErr != nil {
			logger.Fatal("Failed to start HTTPS transport service", zap.Error(serviceErr))
		}
	}()

	go func() {
		<-signalChan
		logger.Info("Received an interrupt, stopping services...")

		if err := grpcTransportService.Stop(); err != nil {
			logger.Error("Failed to stop gRPC transport service", zap.Error(err))
		}

		if err := httpsTansportService.Stop(); err != nil {
			logger.Error("Failed to stop HTTPS transport service", zap.Error(err))
		}

		close(cleanupDone)
	}()
	<-cleanupDone
}

func setupDependencies() (err error) {
	if configurationService, err = configuration.NewEnvConfigurationService(); err != nil {
		return
	}

	repositoryService, err := memory.NewRepositoryService()
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
