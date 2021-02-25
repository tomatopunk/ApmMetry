package server

import (
	"collector/drive/es"
	"collector/handle"
	"collector/processor"
	"collector/protoc/api_v2"
	"context"
	"github.com/jaegertracing/jaeger/model"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"net"
	"testing"
)

type mockSpanProcessor struct {
}

func (p *mockSpanProcessor) Close() error {
	return nil
}

func (p *mockSpanProcessor) ProcessSpans(mSpans []*model.Span, options processor.SpansOptions) ([]bool, error) {
	return []bool{}, nil
}

func TestSpanCollector(t *testing.T) {
	process, _ := processor.NewSpanProcessor(es.ClientConfig{
		Addresses: []string{"http://127.0.0.1:9200/"},
	})

	params := &GRPCServerParams{
		Handler: handle.NewGRPCHandler(process),
	}

	server := grpc.NewServer()
	defer server.Stop()

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer listener.Close()

	go func() {
		rpcServer(listener, server, params)
	}()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	c := api_v2.NewCollectorServiceClient(conn)
	response, err := c.PostSpans(context.Background(), &api_v2.PostSpansRequest{})
	require.NoError(t, err)
	require.NotNil(t, response)
}
