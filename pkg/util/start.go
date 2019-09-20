// Package util implements different utilities required by the tenant service
package util

import grpctransport "github.com/decentralized-cloud/Tenant/transport/grpc"

// StartService setups all dependecies required to start the tenant service and
// start the service
func StartService() {
	server := &grpctransport.Server{}

	server.ListenAndServe()
}
