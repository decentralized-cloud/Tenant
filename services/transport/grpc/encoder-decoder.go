// Package grpc implements functions to expose tenant service endpoint using GRPC protocol.
package grpc

import (
	"context"

	tenantGRPCContract "github.com/decentralized-cloud/tenant/contract/grpc/go"
	"github.com/decentralized-cloud/tenant/models"
	"github.com/decentralized-cloud/tenant/services/business"
	commonErrors "github.com/micro-business/go-core/system/errors"
)

// decodeCreateTenantRequest decodes CreateTenant request message from GRPC object to business object
// context: Mandatory The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeCreateTenantRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*tenantGRPCContract.CreateTenantRequest)

	return &business.CreateTenantRequest{
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
	castedResponse := response.(*business.CreateTenantResponse)

	if castedResponse.Err == nil {
		return &tenantGRPCContract.CreateTenantResponse{
			TenantID: castedResponse.TenantID,
			Error:    tenantGRPCContract.Error_NO_ERROR,
		}, nil
	}

	return &tenantGRPCContract.CreateTenantResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
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

	return &business.ReadTenantRequest{
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
	castedResponse := response.(*business.ReadTenantResponse)

	if castedResponse.Err == nil {
		return &tenantGRPCContract.ReadTenantResponse{
			Tenant: &tenantGRPCContract.Tenant{
				Name: castedResponse.Tenant.Name,
			},
			Error: tenantGRPCContract.Error_NO_ERROR,
		}, nil
	}

	return &tenantGRPCContract.ReadTenantResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

// decodeUpdateTenantRequest decodes UpdateTenant request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeUpdateTenantRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*tenantGRPCContract.UpdateTenantRequest)

	return &business.UpdateTenantRequest{
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
	castedResponse := response.(*business.UpdateTenantResponse)

	if castedResponse.Err == nil {
		return &tenantGRPCContract.UpdateTenantResponse{
			Error: tenantGRPCContract.Error_NO_ERROR,
		}, nil
	}

	return &tenantGRPCContract.UpdateTenantResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

// decodeDeleteTenantRequest decodes DeleteTenant request message from GRPC object to business object
// context: Optional The reference to the context
// request: Mandatory. The reference to the GRPC request
// Returns either the decoded request or error if something goes wrong
func decodeDeleteTenantRequest(
	ctx context.Context,
	request interface{}) (interface{}, error) {
	castedRequest := request.(*tenantGRPCContract.DeleteTenantRequest)

	return &business.DeleteTenantRequest{
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
	castedResponse := response.(*business.DeleteTenantResponse)
	if castedResponse.Err == nil {
		return &tenantGRPCContract.DeleteTenantResponse{
			Error: tenantGRPCContract.Error_NO_ERROR,
		}, nil
	}

	return &tenantGRPCContract.DeleteTenantResponse{
		Error:        mapError(castedResponse.Err),
		ErrorMessage: castedResponse.Err.Error(),
	}, nil
}

func mapError(err error) tenantGRPCContract.Error {
	if business.IsUnknownError(err) {
		return tenantGRPCContract.Error_UNKNOWN
	}

	if business.IsTenantAlreadyExistsError(err) {
		return tenantGRPCContract.Error_TENANT_ALREADY_EXISTS
	}

	if business.IsTenantNotFoundError(err) {
		return tenantGRPCContract.Error_TENANT_NOT_FOUND
	}

	if commonErrors.IsArgumentNilError(err) || commonErrors.IsArgumentError(err) {
		return tenantGRPCContract.Error_BAD_REQUEST
	}

	panic("Error type undefined.")
}
