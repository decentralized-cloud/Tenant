// Package grpctransport implements functions to expose Tenant service endpoint using GRPC protocol.
package grpctransport

import (
	"context"
	"fmt"
	"log"
	"net"

	business "github.com/decentralized-cloud/tenant/business/services"
	configuration "github.com/decentralized-cloud/tenant/configuration/services"
	endpoint "github.com/decentralized-cloud/tenant/endpoint/services"
	repository "github.com/decentralized-cloud/tenant/repository/services"
	"google.golang.org/grpc"

	tenantGRPCContract "github.com/decentralized-cloud/tenant-contract"
	configurationServiceContracts "github.com/decentralized-cloud/tenant/configuration/contracts"
	endpointContracts "github.com/decentralized-cloud/tenant/endpoint/contracts"
	gokitgrpc "github.com/go-kit/kit/transport/grpc"
)

// Server provides function to start GRPC server to serve tenant requests
type Server struct {
	endpointCreatorService endpointContracts.EndpointCreatorContract
	configurationService   configurationServiceContracts.ConfigurationServiceContract

	createTenantHandler gokitgrpc.Handler
	readTenantHandler   gokitgrpc.Handler
	updateTenantHandler gokitgrpc.Handler
	deleteTenantHandler gokitgrpc.Handler
}

// ListenAndServe creates a new GRPC server instance, listens on a port and start serving GRPC requests
func (server *Server) ListenAndServe() {
	if err := server.setupDependencies(); err != nil {
		log.Fatal(err)
	}

	server.setupHandlers()

	errors := make(chan error)

	go func() {

		portNumber, err := server.configurationService.GetPort()
		if err != nil {
			errors <- err
			return
		}

		host, err := server.configurationService.GetHost()
		if err != nil {
			errors <- err
			return
		}

		address := fmt.Sprintf("%s:%d", host, portNumber)
		listener, err := net.Listen("tcp", address)

		if err != nil {
			errors <- err

			return
		}

		gRPCServer := grpc.NewServer()

		tenantGRPCContract.RegisterTenantServiceServer(gRPCServer, server)

		fmt.Printf("gRPC listen on %s\n", address)
		errors <- gRPCServer.Serve(listener)
	}()

	fmt.Println(<-errors)
}

func (server *Server) setupDependencies() error {
	repositoryService, err := repository.NewTenantRepositoryService()

	if err != nil {
		return err
	}

	businessServer, err := business.NewTenantService(repositoryService)

	if err != nil {
		return err
	}

	if server.endpointCreatorService, err = endpoint.NewEndpointCreatorService(businessServer); err != nil {
		return err
	}

	if server.configurationService, err = configuration.NewConfigurationService(); err != nil {
		return err
	}

	return nil
}

// newServer creates a new GRPC server that can serve tenant GRPC requests and process them
func (server *Server) setupHandlers() {
	server.createTenantHandler = gokitgrpc.NewServer(
		server.endpointCreatorService.CreateTenantEndpoint(),
		decodeCreateTenantRequest,
		encodeCreateTenantResponse,
	)

	server.readTenantHandler = gokitgrpc.NewServer(
		server.endpointCreatorService.ReadTenantEndpoint(),
		decodeReadTenantRequest,
		encodeReadTenantResponse,
	)

	server.updateTenantHandler = gokitgrpc.NewServer(
		server.endpointCreatorService.UpdateTenantEndpoint(),
		decodeUpdateTenantRequest,
		encodeUpdateTenantResponse,
	)

	server.deleteTenantHandler = gokitgrpc.NewServer(
		server.endpointCreatorService.DeleteTenantEndpoint(),
		decodeDeleteTenantRequest,
		encodeDeleteTenantResponse,
	)
}

// CreateTenant creates a new tenant
// context: Mandatory. The reference to the context
// request: mandatory. The request to create a new tenant
// Returns the result of creating new tenant
func (server *Server) CreateTenant(
	ctx context.Context,
	request *tenantGRPCContract.CreateTenantRequest) (*tenantGRPCContract.CreateTenantResponse, error) {
	_, response, err := server.createTenantHandler.ServeGRPC(ctx, request)

	if err != nil {
		return nil, err
	}

	return response.(*tenantGRPCContract.CreateTenantResponse), nil
}

// ReadTenant read an existing tenant
// context: Mandatory. The reference to the context
// request: Mandatory. The request to read an esiting tenant
// Returns the result of reading an exiting tenant
func (server *Server) ReadTenant(
	ctx context.Context,
	request *tenantGRPCContract.ReadTenantRequest) (*tenantGRPCContract.ReadTenantResponse, error) {
	_, response, err := server.readTenantHandler.ServeGRPC(ctx, request)

	if err != nil {
		return nil, err
	}

	return response.(*tenantGRPCContract.ReadTenantResponse), nil

}

// UpdateTenant update an existing tenant
// context: Mandatory. The reference to the context
// request: Mandatory. The request to update an esiting tenant
// Returns the result of updateing an exiting tenant
func (server *Server) UpdateTenant(
	ctx context.Context,
	request *tenantGRPCContract.UpdateTenantRequest) (*tenantGRPCContract.UpdateTenantResponse, error) {
	_, response, err := server.updateTenantHandler.ServeGRPC(ctx, request)

	if err != nil {
		return nil, err
	}

	return response.(*tenantGRPCContract.UpdateTenantResponse), nil

}

// DeleteTenant delete an existing tenant
// context: Mandatory. The reference to the context
// request: Mandatory. The request to delete an esiting tenant
// Returns the result of deleting an exiting tenant
func (server *Server) DeleteTenant(
	ctx context.Context,
	request *tenantGRPCContract.DeleteTenantRequest) (*tenantGRPCContract.DeleteTenantResponse, error) {
	_, response, err := server.deleteTenantHandler.ServeGRPC(ctx, request)

	if err != nil {
		return nil, err
	}

	return response.(*tenantGRPCContract.DeleteTenantResponse), nil

}
