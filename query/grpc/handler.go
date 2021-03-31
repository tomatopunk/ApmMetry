package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jaegertracing/jaeger/model"
	"query/drive/es"
	"query/grpc/protoc/api_v2"
	"strings"
	"time"
)

type GRPCHandler struct {
	client es.ElasticsearchClient
	index  string
}

func (g GRPCHandler) GetTrace(request *api_v2.GetTraceRequest, server api_v2.QueryService_GetTraceServer) error {
	currentTime := time.Now()

	trace, err := g.traceIDsSearch(server.Context(), []model.TraceID{request.TraceID}, currentTime)

	if err != nil {
		return err
	}
	return g.sendSpanOnChunks(trace[0].Spans, server.Send)
}

func (g GRPCHandler) sendSpanOnChunks(spans []*model.Span, sendFn func(*api_v2.SpansResponseChunk) error) error {
	chunk := make([]model.Span, 0, len(spans))
	for i := 0; i < len(spans); i += 10 {
		chunk = chunk[:0]
		for j := i; j < len(spans) && j < i+10; j++ {
			chunk = append(chunk, *spans[j])
		}
		if err := sendFn(&api_v2.SpansResponseChunk{Spans: chunk}); err != nil {
			return err
		}
	}
	return nil
}

const findTraceIDsAggregation = `{
    "traceID": {
      "aggs": {
        "startTime": {
          "max": {
            "field": "startTime"
          }
        }
      },
      "terms": {
        "field": "traceID",
        "size": %d,
        "order": {
          "startTime": "desc"
        }
      }
    }
  }
`

func (g GRPCHandler) FindTraces(request *api_v2.FindTracesRequest, server api_v2.QueryService_FindTracesServer) error {

	var queryParam = request.Query
	numTraces := int(queryParam.SearchDepth)

	if numTraces <= 0 {
		numTraces = 20
	}

	searchQuery := es.Query{}
	searchQuery.BoolQuery = map[es.BoolQueryType][]es.BoolQuery{}

	minTimeMicros := model.TimeAsEpochMicroseconds(queryParam.StartTimeMin)
	maxTimeMicros := model.TimeAsEpochMicroseconds(queryParam.StartTimeMax)
	searchQuery.BoolQuery[es.Must] = append(
		searchQuery.BoolQuery[es.Must],
		es.BoolQuery{
			RangeQueries: map[string]es.RangeQuery{
				"startTime": {
					GTE: minTimeMicros,
					LTE: maxTimeMicros,
				},
			},
		},
	)

	if queryParam.DurationMax != 0 || queryParam.DurationMin != 0 {
		minDurationMicros := model.DurationAsMicroseconds(queryParam.DurationMin)
		maxDurationMicros := model.DurationAsMicroseconds(time.Hour * 24)

		if queryParam.DurationMax != 0 {
			maxDurationMicros = model.DurationAsMicroseconds(queryParam.DurationMax)
		}

		searchQuery.BoolQuery[es.Must] = append(searchQuery.BoolQuery[es.Must],
			es.BoolQuery{
				RangeQueries: map[string]es.RangeQuery{
					"duration": {
						GTE: minDurationMicros,
						LTE: maxDurationMicros,
					},
				},
			},
		)
	}

	if queryParam.ServiceName != "" {
		searchQuery.BoolQuery[es.Must] = append(searchQuery.BoolQuery[es.Must],
			es.BoolQuery{
				Term: map[string]string{
					"process.serviceName": queryParam.ServiceName,
				},
			},
		)
	}

	if queryParam.OperationName != "" {
		searchQuery.BoolQuery[es.Must] = append(searchQuery.BoolQuery[es.Must],
			es.BoolQuery{Term: map[string]string{"operationName": queryParam.OperationName}})
	}

	if len(queryParam.Tags) != 0 {
		tagQueries := es.BoolQuery{BoolQuery: map[es.BoolQueryType][]es.BoolQuery{}}

		addTagQuery("tag", queryParam.Tags, tagQueries)
		addTagQuery("process.tag", queryParam.Tags, tagQueries)

		addNestedQuery("tags", queryParam.Tags, tagQueries)

		addNestedQuery("process.tags", queryParam.Tags, tagQueries)

		addNestedQuery("logs.fields", queryParam.Tags, tagQueries)

	}

	aggs := fmt.Sprintf(findTraceIDsAggregation, numTraces)

	searchBody := es.SearchBody{
		Aggregations: json.RawMessage(aggs),
		Query:        &searchQuery,
	}

	response, err := g.client.Search(server.Context(), searchBody, 0, g.index)

	if err != nil {
		return err
	}

	if response.Error != nil {
		return fmt.Errorf("%s", response.Error)
	}

	var traceIds []model.TraceID
	for _, id := range response.Aggs["traceID"].Buckets {
		traceID, err := model.TraceIDFromString(id.Key)
		if err != nil {
			return err
		}
		traceIds = append(traceIds, traceID)
	}
	traces, err := g.traceIDsSearch(server.Context(), traceIds, queryParam.StartTimeMax)
	for _, trace := range traces {
		g.sendSpanOnChunks(trace.Spans, server.Send)
	}
	return nil
}

func addNestedQuery(field string, tags map[string]string, queries es.BoolQuery) {
	k := fmt.Sprintf("%s.%s", field, "key")
	v := fmt.Sprintf("%s.%s", field, "value")

	for tagK, tagV := range tags {

		nestedQuery := &es.NestedQuery{
			Path: field,
			Query: es.Query{
				BoolQuery: map[es.BoolQueryType][]es.BoolQuery{
					es.Must: {
						{MatchQueries: map[string]es.MatchQuery{
							k: {Query: tagK},
						}},
						{Regexp: map[string]es.TermQuery{
							v: {Value: tagV},
						}},
					},
				},
			},
		}
		queries.BoolQuery[es.Should] = append(queries.BoolQuery[es.Should],
			es.BoolQuery{
				Nested: nestedQuery,
			},
		)
	}
}

func addTagQuery(field string, tags map[string]string, queries es.BoolQuery) {
	for k, v := range tags {
		kk := Replace(k)

		key := fmt.Sprintf("%s.%s", field, kk)

		queries.BoolQuery[es.Should] = append(queries.BoolQuery[es.Should],
			es.BoolQuery{
				BoolQuery: map[es.BoolQueryType][]es.BoolQuery{
					es.Must: {
						{Regexp: map[string]es.TermQuery{key: {Value: v}}},
					},
				},
			},
		)

	}
}

func Replace(k string) string {
	return strings.Replace(k, ".", "@", -1)
}

// todo implement me
func (g GRPCHandler) ArchiveTrace(ctx context.Context, request *api_v2.ArchiveTraceRequest) (*api_v2.ArchiveTraceResponse, error) {
	panic("implement me")

	//request.TraceID
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

func (g GRPCHandler) traceIDsSearch(ctx context.Context, traceIDs []model.TraceID, endTime time.Time) ([]*model.Trace, error) {
	if len(traceIDs) <= 0 {
		return []*model.Trace{}, nil
	}

	queries := make([]es.SearchBody, len(traceIDs))
	for i, traceID := range traceIDs {
		searchBody := es.SearchBody{
			Indices: []string{g.index},
			Query: &es.Query{
				Term: &es.Terms{
					"traceID": es.TermQuery{
						Value: traceID.String(),
					},
				},
			},
		}
		queries[i] = searchBody
	}

	//res,err := g.client.MultiSearch(ctx,queries)
	panic("imp")
}
