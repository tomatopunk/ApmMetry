package grpc

import (
	"context"
	"query/grpc/protoc/api_v2"
)

type GRPCHandler struct {
}

func (G GRPCHandler) GetTrace(request *api_v2.GetTraceRequest, server api_v2.QueryService_GetTraceServer) error {
	panic("implement me")
}

func (G GRPCHandler) ArchiveTrace(ctx context.Context, request *api_v2.ArchiveTraceRequest) (*api_v2.ArchiveTraceResponse, error) {
	panic("implement me")

	//request.TraceID
}

func (G GRPCHandler) FindTraces(request *api_v2.FindTracesRequest, server api_v2.QueryService_FindTracesServer) error {
	panic("implement me")
}

func (G GRPCHandler) GetServices(ctx context.Context, request *api_v2.GetServicesRequest) (*api_v2.GetServicesResponse, error) {
	panic("implement me")
}

func (G GRPCHandler) GetOperations(ctx context.Context, request *api_v2.GetOperationsRequest) (*api_v2.GetOperationsResponse, error) {
	panic("implement me")
}

func (G GRPCHandler) GetDependencies(ctx context.Context, request *api_v2.GetDependenciesRequest) (*api_v2.GetDependenciesResponse, error) {
	panic("implement me")
}

func NewGRPCHandler() *GRPCHandler {
	return &GRPCHandler{}
}
