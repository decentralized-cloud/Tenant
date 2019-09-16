package main

import grpctransport "github.com/decentralized-cloud/Tenant/transport/grpc"

func main() {
	server := &grpctransport.Server{}

	server.ListenAndServe()
}
