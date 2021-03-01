package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"query/grpc/protoc/api_v2"
)

func CreateServer() error {
	conn, err := net.Listen("", "8080")
	if err != nil {
		return fmt.Errorf("Query listen is faild {%d}", err)
	}
	server, err := createGRPCServer()

	if err != nil {
		return fmt.Errorf("Create grpc server is faild :{%s}", err)
	}

	err = server.Serve(conn)
	if err != nil {
		return fmt.Errorf("Could not start GRPC server :{%s}", err)
	}
	return nil
}

func createGRPCServer() (*grpc.Server, error) {
	var grpcOpts []grpc.ServerOption

	server := grpc.NewServer(grpcOpts...)
	handler := NewGRPCHandler()
	api_v2.RegisterQueryServiceServer(server, handler)
	return server, nil
}
