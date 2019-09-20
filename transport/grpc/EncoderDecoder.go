// Package grpctransport implements functions to expose Tenant service endpoint using GRPC protocol.
package grpctransport

import (
	"context"

	tenantGRPCContract "github.com/decentralized-cloud/tenant-contract"
	businessContracts "github.com/decentralized-cloud/tenant/business/contracts"
	"github.com/decentralized-cloud/tenant/models"
)

// decodeCreateTenantRequest decodes CreateTenant request message from GRPC object to business object
// context: Mandatory The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeCreateTenantRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*tenantGRPCContract.CreateTenantRequest)

	return &businessContracts.CreateTenantRequest{
		Tenant: models.Tenant{
			Name: castedRequest.Tenant.Name,
		}}, nil
}

// encodeCreateTenantResponse encodes CreateTenant response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeCreateTenantResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*businessContracts.CreateTenantResponse)

	return &tenantGRPCContract.CreateTenantResponse{
		TenantID: castedResponse.TenantID,
	}, nil
}

// decodeReadTenantRequest decodes ReadTenant request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeReadTenantRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*tenantGRPCContract.ReadTenantRequest)

	return &businessContracts.ReadTenantRequest{
		TenantID: castedRequest.TenantID,
	}, nil
}

// encodeReadTenantResponse encodes ReadTenant response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeReadTenantResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	castedResponse := response.(*businessContracts.ReadTenantResponse)

	return &tenantGRPCContract.ReadTenantResponse{
		Tenant: &tenantGRPCContract.Tenant{
			Name: castedResponse.Tenant.Name,
		}}, nil
}

// decodeUpdateTenantRequest decodes UpdateTenant request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeUpdateTenantRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*tenantGRPCContract.UpdateTenantRequest)

	return &businessContracts.UpdateTenantRequest{
		TenantID: castedRequest.TenantID,
		Tenant: models.Tenant{
			Name: castedRequest.Tenant.Name,
		}}, nil
}

// encodeUpdateTenantResponse encodes UpdateTenant response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeUpdateTenantResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	return &tenantGRPCContract.UpdateTenantResponse{}, nil
}

// decodeDeleteTenantRequest decodes DeleteTenant request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeDeleteTenantRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*tenantGRPCContract.DeleteTenantRequest)

	return &businessContracts.DeleteTenantRequest{
		TenantID: castedRequest.TenantID,
	}, nil
}

// encodeDeleteTenantResponse encodes DeleteTenant response from business object to GRPC object
// context: Optional The reference to the context
// request: Mandatory. The reference to the business response
// Returns either the decoded response or error if something goes wrong
func encodeDeleteTenantResponse(
	ctx context.Context,
	response interface{}) (interface{}, error) {
	return &tenantGRPCContract.DeleteTenantResponse{}, nil
}
