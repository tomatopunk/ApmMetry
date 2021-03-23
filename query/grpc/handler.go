package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jaegertracing/jaeger/model"
	"query/drive/es"
	"query/grpc/protoc/api_v2"
	"time"
)

type GRPCHandler struct {
	client es.ElasticsearchClient
	index  string
}

func (g GRPCHandler) GetTrace(request *api_v2.GetTraceRequest, server api_v2.QueryService_GetTraceServer) error {
	currentTime := time.Now()
	trace, err := g.traceIDsSearch(server.Context(), []model.TraceID{request.TraceID}, currentTime)

	panic("implement me")
}

// todo implement me
func (g GRPCHandler) ArchiveTrace(ctx context.Context, request *api_v2.ArchiveTraceRequest) (*api_v2.ArchiveTraceResponse, error) {
	panic("implement me")

	//request.TraceID
}

func (g GRPCHandler) FindTraces(request *api_v2.FindTracesRequest, server api_v2.QueryService_FindTracesServer) error {
	panic("implement me")
}

const getServicesAggregation = `{
    "serviceName": {
      "terms": {
        "field": "serviceName",
        "size": %d
      }
    }
  }
`

func (g GRPCHandler) GetServices(ctx context.Context, request *api_v2.GetServicesRequest) (*api_v2.GetServicesResponse, error) {
	aggs := fmt.Sprintf(getServicesAggregation, 1000)
	searchBody := es.SearchBody{
		Aggregations: json.RawMessage(aggs),
	}
	res, err := g.client.Search(ctx, searchBody, 0, g.index)

	if err != nil {
		return nil, err
	}

	if res.Error != nil {
		return nil, fmt.Errorf("%s", res.Error)
	}

	var serviceNames []string

	for _, k := range res.Aggs["serviceName"].Buckets {
		serviceNames = append(serviceNames, k.Key)
	}
	return &api_v2.GetServicesResponse{Services: serviceNames}, nil
}

const getOperationsAggregation = `{
    "operationName": {
      "terms": {
        "field": "operationName",
        "size": %d
      }
    }
  }
`

func (g GRPCHandler) GetOperations(ctx context.Context, request *api_v2.GetOperationsRequest) (*api_v2.GetOperationsResponse, error) {

	aggs := fmt.Sprintf(getOperationsAggregation, "1000")
	searchBody := es.SearchBody{
		Aggregations: json.RawMessage(aggs),
		Query: &es.Query{
			Term: &es.Terms{
				"serviceName": es.TermQuery{Value: request.Service},
			},
		},
	}

	res, err := g.client.Search(ctx, searchBody, 0, g.index)

	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, fmt.Errorf("%s", res.Error)
	}

	var operations []*api_v2.Operation

	for _, k := range res.Aggs["operationName"].Buckets {
		operations = append(operations, &api_v2.Operation{
			Name: k.Key,
		})
	}
	return &api_v2.GetOperationsResponse{
		OperationNames: uniqueOperationName(operations),
		Operations:     operations,
	}, nil
}

//todo implement me
func (g GRPCHandler) GetDependencies(ctx context.Context, request *api_v2.GetDependenciesRequest) (*api_v2.GetDependenciesResponse, error) {
	panic("implement me")
}

func NewGRPCHandler() *GRPCHandler {
	return &GRPCHandler{}
}

func uniqueOperationName(operations []*api_v2.Operation) []string {
	set := make(map[string]struct{})
	for _, operation := range operations {
		set[operation.Name] = struct{}{}
	}
	var operationNames []string
	for operation := range set {
		operationNames = append(operationNames, operation)
	}
	return operationNames
}

func (h GRPCHandler) traceIDsSearch(ctx context.Context, traceIDs []model.TraceID, endTime time.Time) ([]*model.Trace, error) {

}
