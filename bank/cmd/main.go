package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"software-webservices-go/bank"
)

func main() {
	// this is going to get a bit weird, but we can always refactor later
	socket, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen on port %v\n", err)
	}
	var opts []grpc.ServerOption
	// we need to create a new gRPC server
	grpcServer := grpc.NewServer(opts...)
	bank.RegisterBankServer(grpcServer, bank.NewBankService())
	// Timeout should be 120 seconds by default for the HTTP connection establishment
	// TODO: add UnaryInterceptor timeout for the gRPC handlers
	// https://pkg.go.dev/google.golang.org/grpc#UnaryInterceptor
	// https://github.com/grpc/grpc-go/issues/5059
	err = grpcServer.Serve(socket)
	if err != nil {
		log.Fatalf("Failed to serve gRPC server over port %v\n", err)
	}
}
