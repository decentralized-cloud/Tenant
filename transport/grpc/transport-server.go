// Package grpctransport implements functions to expose Tenant service endpoint using GRPC protocol.
package grpctransport

import (
	"context"
	"fmt"
	"net"

	tenantGRPCContract "github.com/decentralized-cloud/tenant-contract/grpc"
	configurationServiceContract "github.com/decentralized-cloud/tenant/services/configuration/contract"
	endpointContract "github.com/decentralized-cloud/tenant/services/endpoint/contract"
	transportContract "github.com/decentralized-cloud/tenant/transport/contract"
	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	commonErrors "github.com/micro-business/go-core/system/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type transportService struct {
	logger                 *zap.Logger
	endpointCreatorService endpointContract.EndpointCreatorContract
	configurationService   configurationServiceContract.ConfigurationServiceContract
	createTenantHandler    gokitgrpc.Handler
	readTenantHandler      gokitgrpc.Handler
	updateTenantHandler    gokitgrpc.Handler
	deleteTenantHandler    gokitgrpc.Handler
}

// NewTransportService creates new instance of the GRPCService, setting up all dependencies and returns the instance
// logger: Mandatory. Reference to the logger service
// configurationService: Mandatory. Reference to the service that provides required configurations
// endpointCreatorService: Mandatory. Reference to the service that creates go-kit compatible endpoints
// Returns the new service or error if something goes wrong
func NewTransportService(
	logger *zap.Logger,
	configurationService configurationServiceContract.ConfigurationServiceContract,
	endpointCreatorService endpointContract.EndpointCreatorContract) (transportContract.TransportServiceContract, error) {
	if logger == nil {
		return nil, commonErrors.NewArgumentNilError("logger", "logger is required")
	}

	if configurationService == nil {
		return nil, commonErrors.NewArgumentNilError("configurationService", "configurationService is required")
	}

	if endpointCreatorService == nil {
		return nil, commonErrors.NewArgumentNilError("endpointCreatorService", "endpointCreatorService is required")
	}

	return &transportService{
		logger:                 logger,
		configurationService:   configurationService,
		endpointCreatorService: endpointCreatorService,
	}, nil
}

// Start starts the GRPC transport service
// Returns error if something goes wrong
func (service *transportService) Start() error {
	service.setupHandlers()

	portNumber, err := service.configurationService.GetGRPCPort()
	if err != nil {
		return err
	}

	host, err := service.configurationService.GetGRPCHost()
	if err != nil {
		return err
	}

	address := fmt.Sprintf("%s:%d", host, portNumber)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	gRPCServer := grpc.NewServer()
	tenantGRPCContract.RegisterTenantServiceServer(gRPCServer, service)
	service.logger.Info("gRPC server started", zap.String("address", address))

	return gRPCServer.Serve(listener)
}

// Stop stops the GRPC transport service
// Returns error if something goes wrong
func (service *transportService) Stop() error {
	return nil
}

// newServer creates a new GRPC server that can serve tenant GRPC requests and process them
func (service *transportService) setupHandlers() {
	service.createTenantHandler = gokitgrpc.NewServer(
		service.endpointCreatorService.CreateTenantEndpoint(),
		decodeCreateTenantRequest,
		encodeCreateTenantResponse,
	)

	service.readTenantHandler = gokitgrpc.NewServer(
		service.endpointCreatorService.ReadTenantEndpoint(),
		decodeReadTenantRequest,
		encodeReadTenantResponse,
	)

	service.updateTenantHandler = gokitgrpc.NewServer(
		service.endpointCreatorService.UpdateTenantEndpoint(),
		decodeUpdateTenantRequest,
		encodeUpdateTenantResponse,
	)

	service.deleteTenantHandler = gokitgrpc.NewServer(
		service.endpointCreatorService.DeleteTenantEndpoint(),
		decodeDeleteTenantRequest,
		encodeDeleteTenantResponse,
	)
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
