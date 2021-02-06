// Package grpc implements functions to expose project service endpoint using GRPC protocol.
package grpc

import (
	"context"
	"fmt"
	"net"

	projectGRPCContract "github.com/decentralized-cloud/project/contract/grpc/go"
	"github.com/decentralized-cloud/project/services/configuration"
	"github.com/decentralized-cloud/project/services/endpoint"
	"github.com/decentralized-cloud/project/services/transport"
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
	createProjectHandler      gokitgrpc.Handler
	readProjectHandler        gokitgrpc.Handler
	updateProjectHandler      gokitgrpc.Handler
	deleteProjectHandler      gokitgrpc.Handler
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
	projectGRPCContract.RegisterProjectServiceServer(gRPCServer, service)
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
	var createProjectEndpoint gokitEndpoint.Endpoint
	{
		createProjectEndpoint = service.endpointCreatorService.CreateProjectEndpoint()
		createProjectEndpoint = service.middlewareProviderService.CreateLoggingMiddleware("CreateProject")(createProjectEndpoint)
		service.createProjectHandler = gokitgrpc.NewServer(
			createProjectEndpoint,
			decodeCreateProjectRequest,
			encodeCreateProjectResponse,
		)
	}

	var readProjectEndpoint gokitEndpoint.Endpoint
	{
		readProjectEndpoint = service.endpointCreatorService.ReadProjectEndpoint()
		readProjectEndpoint = service.middlewareProviderService.CreateLoggingMiddleware("ReadProject")(readProjectEndpoint)
		service.readProjectHandler = gokitgrpc.NewServer(
			readProjectEndpoint,
			decodeReadProjectRequest,
			encodeReadProjectResponse,
		)
	}

	var updateProjectEndpoint gokitEndpoint.Endpoint
	{
		updateProjectEndpoint = service.endpointCreatorService.UpdateProjectEndpoint()
		updateProjectEndpoint = service.middlewareProviderService.CreateLoggingMiddleware("UpdateProject")(updateProjectEndpoint)
		service.updateProjectHandler = gokitgrpc.NewServer(
			updateProjectEndpoint,
			decodeUpdateProjectRequest,
			encodeUpdateProjectResponse,
		)
	}

	var deleteProjectEndpoint gokitEndpoint.Endpoint
	{
		deleteProjectEndpoint = service.endpointCreatorService.DeleteProjectEndpoint()
		deleteProjectEndpoint = service.middlewareProviderService.CreateLoggingMiddleware("DeleteProject")(deleteProjectEndpoint)
		service.deleteProjectHandler = gokitgrpc.NewServer(
			deleteProjectEndpoint,
			decodeDeleteProjectRequest,
			encodeDeleteProjectResponse,
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

// CreateProject creates a new project
// context: Mandatory. The reference to the context
// request: mandatory. The request to create a new project
// Returns the result of creating new project
func (service *transportService) CreateProject(
	ctx context.Context,
	request *projectGRPCContract.CreateProjectRequest) (*projectGRPCContract.CreateProjectResponse, error) {
	_, response, err := service.createProjectHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*projectGRPCContract.CreateProjectResponse), nil
}

// ReadProject read an existing project
// context: Mandatory. The reference to the context
// request: Mandatory. The request to read an existing project
// Returns the result of reading an exiting project
func (service *transportService) ReadProject(
	ctx context.Context,
	request *projectGRPCContract.ReadProjectRequest) (*projectGRPCContract.ReadProjectResponse, error) {
	_, response, err := service.readProjectHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*projectGRPCContract.ReadProjectResponse), nil

}

// UpdateProject update an existing project
// context: Mandatory. The reference to the context
// request: Mandatory. The request to update an existing project
// Returns the result of updateing an exiting project
func (service *transportService) UpdateProject(
	ctx context.Context,
	request *projectGRPCContract.UpdateProjectRequest) (*projectGRPCContract.UpdateProjectResponse, error) {
	_, response, err := service.updateProjectHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*projectGRPCContract.UpdateProjectResponse), nil

}

// DeleteProject delete an existing project
// context: Mandatory. The reference to the context
// request: Mandatory. The request to delete an existing project
// Returns the result of deleting an exiting project
func (service *transportService) DeleteProject(
	ctx context.Context,
	request *projectGRPCContract.DeleteProjectRequest) (*projectGRPCContract.DeleteProjectResponse, error) {
	_, response, err := service.deleteProjectHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*projectGRPCContract.DeleteProjectResponse), nil

}

// Search returns the list  of project that matched the provided criteria
// context: Mandatory. The reference to the context
// request: Mandatory. The request contains the filter criteria to look for existing project
// Returns the list of project that matched the provided criteria
func (service *transportService) Search(
	ctx context.Context,
	request *projectGRPCContract.SearchRequest) (*projectGRPCContract.SearchResponse, error) {
	_, response, err := service.searchHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*projectGRPCContract.SearchResponse), nil
}
