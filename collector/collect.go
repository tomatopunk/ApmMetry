package main

import (
	"collector/handle"
	"collector/processor"
	"collector/server"
	"fmt"
	"google.golang.org/grpc"
)

type Collector struct {
	gRPCServer *grpc.Server
}

func New() *Collector {
	return &Collector{}
}

func (c *Collector) StartServer(options *CollectOptions) error {
	process := processor.NewSpanProcessor()

	gRPCServer, err := server.StartGRPCServer(&server.GRPCServerParams{
		HostPort: options.CollectorGRPCHostPort,
		Handler:  handle.NewGRPCHandler(process),
	})
	c.gRPCServer = gRPCServer
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (c *Collector) Close() error {
	if c.gRPCServer != nil {
		c.gRPCServer.GracefulStop()
	}
	return nil
}
