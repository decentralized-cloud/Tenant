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
	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/micro-business/go-core/gokit/middleware"
	commonErrors "github.com/micro-business/go-core/system/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type transportService struct {
	logger                    *zap.Logger
	configurationService      configuration.ConfigurationContract
	endpointCreatorService    endpoint.EndpointCreatorContract
	middlewareProviderService middleware.MiddlewareProviderContract
	jwksURL                   string
	createProjectHandler      gokitgrpc.Handler
	readProjectHandler        gokitgrpc.Handler
	updateProjectHandler      gokitgrpc.Handler
	deleteProjectHandler      gokitgrpc.Handler
	ListProjectsHandler       gokitgrpc.Handler
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

	jwksURL, err := configurationService.GetJwksURL()
	if err != nil {
		return nil, err
	}

	return &transportService{
		logger:                    logger,
		configurationService:      configurationService,
		endpointCreatorService:    endpointCreatorService,
		middlewareProviderService: middlewareProviderService,
		jwksURL:                   jwksURL,
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
	projectGRPCContract.RegisterServiceServer(gRPCServer, service)
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
	endpoint := service.endpointCreatorService.CreateProjectEndpoint()
	endpoint = service.middlewareProviderService.CreateLoggingMiddleware("CreateProject")(endpoint)
	endpoint = service.createAuthMiddleware()(endpoint)
	service.createProjectHandler = gokitgrpc.NewServer(
		endpoint,
		decodeCreateProjectRequest,
		encodeCreateProjectResponse,
	)

	endpoint = service.endpointCreatorService.ReadProjectEndpoint()
	endpoint = service.middlewareProviderService.CreateLoggingMiddleware("ReadProject")(endpoint)
	endpoint = service.createAuthMiddleware()(endpoint)
	service.readProjectHandler = gokitgrpc.NewServer(
		endpoint,
		decodeReadProjectRequest,
		encodeReadProjectResponse,
	)

	endpoint = service.endpointCreatorService.UpdateProjectEndpoint()
	endpoint = service.middlewareProviderService.CreateLoggingMiddleware("UpdateProject")(endpoint)
	endpoint = service.createAuthMiddleware()(endpoint)
	service.updateProjectHandler = gokitgrpc.NewServer(
		endpoint,
		decodeUpdateProjectRequest,
		encodeUpdateProjectResponse,
	)

	endpoint = service.endpointCreatorService.DeleteProjectEndpoint()
	endpoint = service.middlewareProviderService.CreateLoggingMiddleware("DeleteProject")(endpoint)
	endpoint = service.createAuthMiddleware()(endpoint)
	service.deleteProjectHandler = gokitgrpc.NewServer(
		endpoint,
		decodeDeleteProjectRequest,
		encodeDeleteProjectResponse,
	)

	endpoint = service.endpointCreatorService.ListProjectsEndpoint()
	endpoint = service.middlewareProviderService.CreateLoggingMiddleware("ListProjects")(endpoint)
	endpoint = service.createAuthMiddleware()(endpoint)
	service.ListProjectsHandler = gokitgrpc.NewServer(
		endpoint,
		decodeListProjectsRequest,
		encodeListProjectsResponse,
	)
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
// Returns the result of reading an existing project
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
// Returns the result of updateing an existing project
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
// Returns the result of deleting an existing project
func (service *transportService) DeleteProject(
	ctx context.Context,
	request *projectGRPCContract.DeleteProjectRequest) (*projectGRPCContract.DeleteProjectResponse, error) {
	_, response, err := service.deleteProjectHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*projectGRPCContract.DeleteProjectResponse), nil

}

// ListProjects returns the list  of project that matched the provided criteria
// context: Mandatory. The reference to the context
// request: Mandatory. The request contains the filter criteria to look for existing project
// Returns the list of project that matched the provided criteria
func (service *transportService) ListProjects(
	ctx context.Context,
	request *projectGRPCContract.ListProjectsRequest) (*projectGRPCContract.ListProjectsResponse, error) {
	_, response, err := service.ListProjectsHandler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.(*projectGRPCContract.ListProjectsResponse), nil
}
