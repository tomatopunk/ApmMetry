package server

import (
	"collector/handle"
	"collector/protoc/api_v2"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type GRPCServerParams struct {
	HostPort string
	Handler  *handle.GRPCHandler
}

func StartGRPCServer(params *GRPCServerParams) (*grpc.Server, error) {
	server := grpc.NewServer()

	listener, err := net.Listen("TCP", params.HostPort)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := rpcServer(listener, server, params); err != nil {

	}
	return server, nil
}

func rpcServer(listener net.Listener, server *grpc.Server, params *GRPCServerParams) error {
	api_v2.RegisterCollectorServiceServer(server, params.Handler)
	fmt.Println("Start gRPC server")

	go func() {
		if err := server.Serve(listener); err != nil {
			fmt.Errorf("%w", err)
		}
	}()
	return nil
}
