// Package https implements functions to expose tenant service endpoint using HTTPS protocol.
package https

import (
	"fmt"
	"net/http"

	"github.com/decentralized-cloud/tenant/services/configuration"
	"github.com/decentralized-cloud/tenant/services/transport"
	"github.com/decentralized-cloud/tenant/services/transport/grpc"
	commonErrors "github.com/micro-business/go-core/system/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/savsgio/atreugo/v9"
	"go.uber.org/zap"
)

type transportService struct {
	logger               *zap.Logger
	configurationService configuration.ConfigurationContract
}

// NewTransportService creates new instance of the GRPCService, setting up all dependencies and returns the instance
// logger: Mandatory. Reference to the logger service
// configurationService: Mandatory. Reference to the service that provides required configurations
// Returns the new service or error if something goes wrong
func NewTransportService(
	logger *zap.Logger,
	configurationService configuration.ConfigurationContract) (transport.TransportContract, error) {
	if logger == nil {
		return nil, commonErrors.NewArgumentNilError("logger", "logger is required")
	}

	if configurationService == nil {
		return nil, commonErrors.NewArgumentNilError("configurationService", "configurationService is required")
	}

	return &transportService{
		logger:               logger,
		configurationService: configurationService,
	}, nil
}

// Start starts the GraphQL transport service
// Returns error if something goes wrong
func (service *transportService) Start() error {
	config := &atreugo.Config{GracefulShutdown: true}
	var err error

	host, err := service.configurationService.GetHttpsHost()
	if err != nil {
		return err
	}

	port, err := service.configurationService.GetHttpsPort()
	if err != nil {
		return err
	}

	config.Addr = fmt.Sprintf("%s:%d", host, port)
	server := atreugo.New(config)

	server.Path("GET", "/live", service.livenessCheckHandler)
	server.Path("GET", "/ready", service.readinessCheckHandler)
	server.NetHTTPPath("GET", "/metrics", promhttp.Handler())
	service.logger.Info("HTTPS server started", zap.String("address", config.Addr))

	return server.ListenAndServe()
}

// Stop stops the GraphQL transport service
// Returns error if something goes wrong
func (service *transportService) Stop() error {
	return nil
}

func (service *transportService) livenessCheckHandler(ctx *atreugo.RequestCtx) error {
	if grpc.Live {
		ctx.Response.SetStatusCode(http.StatusOK)
	} else {
		ctx.Response.SetStatusCode(http.StatusServiceUnavailable)
	}

	return nil
}

func (service *transportService) readinessCheckHandler(ctx *atreugo.RequestCtx) error {
	if grpc.Ready {
		ctx.Response.SetStatusCode(http.StatusOK)
	} else {
		ctx.Response.SetStatusCode(http.StatusServiceUnavailable)
	}

	return nil
}
