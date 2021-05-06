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

	listener, err := net.Listen("tcp", params.HostPort)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := rpcServer(listener, server, params); err != nil {
		return nil, err
	}
	return server, nil
}

func rpcServer(listener net.Listener, server *grpc.Server, params *GRPCServerParams) error {
	api_v2.RegisterCollectorServiceServer(server, params.Handler)
	fmt.Println("Start gRPC server")

	if err := server.Serve(listener); err != nil {
		fmt.Errorf("Could not launch gRPC service {%w}", err)
	}

	//go func() {
	//	if err := server.Serve(listener); err != nil {
	//		fmt.Errorf("Could not launch gRPC service {%w}", err)
	//	}
	//}()
	return nil
}
