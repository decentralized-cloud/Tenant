// Package util implements different utilities required by the tenant service
package util

import (
	"log"
	"os"
	"os/signal"

	grpctransport "github.com/decentralized-cloud/tenant/transport/grpc"
	"go.uber.org/zap"
)

// StartService setups all dependecies required to start the tenant service and
// start the service
func StartService() {
	server := &grpctransport.Server{}

	go server.StartListenAndServe()

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan

		logger, err := zap.NewProduction()
		if err != nil {
			log.Fatal(err)
		}

		defer logger.Sync()

		logger.Info("Received an interrupt, stopping services...")

		close(cleanupDone)
	}()
	<-cleanupDone
}
