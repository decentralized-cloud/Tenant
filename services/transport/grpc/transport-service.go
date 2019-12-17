// Package grpc implements functions to expose tenant service endpoint using GRPC protocol.
package grpc

import (
	"context"
	"fmt"
	"net"

	tenantGRPCContract "github.com/decentralized-cloud/tenant/contract/grpc/go"
	"github.com/decentralized-cloud/tenant/services/configuration"
	"github.com/decentralized-cloud/tenant/services/endpoint"
	"github.com/decentralized-cloud/tenant/services/transport"
	gokitEndpoint "github.com/go-kit/kit/endpoint"
	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	commonErrors "github.com/micro-business/go-core/system/errors"
	"github.com/micro-business/gokit-core/middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type transportService struct {
	logger                    *zap.Logger
	configurationService      configuration.ConfigurationContract
	endpointCreatorService    endpoint.EndpointCreatorContract
	middlewareProviderService middleware.MiddlewareProviderContract
	createTenantHandler       gokitgrpc.Handler
	readTenantHandler         gokitgrpc.Handler
	updateTenantHandler       gokitgrpc.Handler
	deleteTenantHandler       gokitgrpc.Handler
	searchHandler             gokitgrpc.Handler
}

var Live bool
var Ready bool

func init() {
	Live = false
	Ready = false
}

// NewTransportService creates new instance of the transportService, setting up all dependencies and returns the instance
// logger: Mandatory. Reference to the logger service
// configurationService: Mandatory. Reference to the service that provides required configurations
// endpointCreatorService: Mandatory. Reference to the service that creates go-kit compatible endpoints
// middlewareProviderService: Mandatory. Reference to the service that provides different go-kit middlewares
// Returns the new service or error if something goes wrong
func NewTransportService(
	logger *zap.Logger,
	configurationService configuration.ConfigurationContract,
	endpointCreatorService endpoint.EndpointCreatorContract,
	middlewareProviderService middleware.MiddlewareProviderContract) (transport.TransportContract, error) {
	if logger == nil {
		return nil, commonErrors.NewArgumentNilError("logger", "logger is required")
	}

	if configurationService == nil {
		return nil, commonErrors.NewArgumentNilError("configurationService", "configurationService is required")
	}

	if endpointCreatorService == nil {
		return nil, commonErrors.NewArgumentNilError("endpointCreatorService", "endpointCreatorService is required")
	}

	if middlewareProviderService == nil {
		return nil, commonErrors.NewArgumentNilError("middlewareProviderService", "middlewareProviderService is required")
	}

	return &transportService{
		logger:                    logger,
		configurationService:      configurationService,
		endpointCreatorService:    endpointCreatorService,
		middlewareProviderService: middlewareProviderService,
	}, nil
}

// Start starts the GRPC transport service
// Returns error if something goes wrong
func (service *transportService) Start() error {
	service.setupHandlers()

	host, err := service.configurationService.GetGrpcHost()
	if err != nil {
		return err
	}

	port, err := service.configurationService.GetGrpcPort()
	if err != nil {
		return err
	}

	address := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	gRPCServer := grpc.NewServer()
	tenantGRPCContract.RegisterTenantServiceServer(gRPCServer, service)
	service.logger.Info("gRPC service started", zap.String("address", address))

	Live = true
	Ready = true

	err = gRPCServer.Serve(listener)

	Live = false
	Ready = false

	return err
}

// Stop stops the GRPC transport service
// Returns error if something goes wrong
func (service *transportService) Stop() error {
	return nil
}

func (service *transportService) setupHandlers() {
	var createTenantEndpoint gokitEndpoint.Endpoint
	{
		createTenantEndpoint = service.endpointCreatorService.CreateTenantEndpoint()
		createTenantEndpoint = service.middlewareProviderService.CreateLoggingMiddleware("CreateTenant")(createTenantEndpoint)
		service.createTenantHandler = gokitgrpc.NewServer(
			createTenantEndpoint,
			decodeSearchRequest,
			encodeSearchResponse,
		)
	}

	var readTenantEndpoint gokitEndpoint.Endpoint
	{
		readTenantEndpoint = service.endpointCreatorService.ReadTenantEndpoint()
		readTenantEndpoint = service.middlewareProviderService.CreateLoggingMiddleware("ReadTenant")(readTenantEndpoint)
		service.readTenantHandler = gokitgrpc.NewServer(
			readTenantEndpoint,
			decodeSearchRequest,
			encodeSearchResponse,
		)
	}

	var updateTenantEndpoint gokitEndpoint.Endpoint
	{
		updateTenantEndpoint = service.endpointCreatorService.UpdateTenantEndpoint()
		updateTenantEndpoint = service.middlewareProviderService.CreateLoggingMiddleware("UpdateTenant")(updateTenantEndpoint)
		service.updateTenantHandler = gokitgrpc.NewServer(
			updateTenantEndpoint,
			decodeSearchRequest,
			encodeSearchResponse,
		)
	}

	var deleteTenantEndpoint gokitEndpoint.Endpoint
	{
		deleteTenantEndpoint = service.endpointCreatorService.DeleteTenantEndpoint()
		deleteTenantEndpoint = service.middlewareProviderService.CreateLoggingMiddleware("DeleteTenant")(deleteTenantEndpoint)
		service.deleteTenantHandler = gokitgrpc.NewServer(
			deleteTenantEndpoint,
			decodeSearchRequest,
			encodeSearchResponse,
		)
	}

	var searchEndpoint gokitEndpoint.Endpoint
	{
		searchEndpoint = service.endpointCreatorService.SearchEndpoint()
		searchEndpoint = service.middlewareProviderService.CreateLoggingMiddleware("Search")(searchEndpoint)
		service.searchHandler = gokitgrpc.NewServer(
			searchEndpoint,
			decodeSearchRequest,
			encodeSearchResponse,
		)
	}
}

// CreateTenant creates a new tenant
// context: Mandatory. The reference to the context
// request: mandatory. The request to create a new tenant
// Returns the result of creating new tenant
func (service *transportService) CreateTenant(
	ctx context.Context,
	request *tenantGRPCContract.CreateTenantRequest) (*tenantGRPCContract.CreateTenantResponse, error) {
	_, response, err := service.createTenantHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*tenantGRPCContract.CreateTenantResponse), nil
}

// ReadTenant read an existing tenant
// context: Mandatory. The reference to the context
// request: Mandatory. The request to read an existing tenant
// Returns the result of reading an exiting tenant
func (service *transportService) ReadTenant(
	ctx context.Context,
	request *tenantGRPCContract.ReadTenantRequest) (*tenantGRPCContract.ReadTenantResponse, error) {
	_, response, err := service.readTenantHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*tenantGRPCContract.ReadTenantResponse), nil

}

// UpdateTenant update an existing tenant
// context: Mandatory. The reference to the context
// request: Mandatory. The request to update an existing tenant
// Returns the result of updateing an exiting tenant
func (service *transportService) UpdateTenant(
	ctx context.Context,
	request *tenantGRPCContract.UpdateTenantRequest) (*tenantGRPCContract.UpdateTenantResponse, error) {
	_, response, err := service.updateTenantHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*tenantGRPCContract.UpdateTenantResponse), nil

}

// DeleteTenant delete an existing tenant
// context: Mandatory. The reference to the context
// request: Mandatory. The request to delete an existing tenant
// Returns the result of deleting an exiting tenant
func (service *transportService) DeleteTenant(
	ctx context.Context,
	request *tenantGRPCContract.DeleteTenantRequest) (*tenantGRPCContract.DeleteTenantResponse, error) {
	_, response, err := service.deleteTenantHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*tenantGRPCContract.DeleteTenantResponse), nil

}

// Search returns the list  of tenant that matched the provided criteria
// context: Mandatory. The reference to the context
// request: Mandatory. The request contains the filter criteria to look for existing tenant
// Returns the list of tenant that matched the provided criteria
func (service *transportService) Search(
	ctx context.Context,
	request *tenantGRPCContract.SearchRequest) (*tenantGRPCContract.SearchResponse, error) {
	_, response, err := service.searchHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*tenantGRPCContract.SearchResponse), nil
}
