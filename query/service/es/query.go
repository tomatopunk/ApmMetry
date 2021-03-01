package es

import (
	"context"
	"github.com/jaegertracing/jaeger/model"
	"query/service"
	"time"
)

type ElasticsearchQueryService struct {
}

func (e ElasticsearchQueryService) GetTrace(ctx context.Context, traceID model.TraceID) (*model.Trace, error) {
	panic("implement me")
}

func (e ElasticsearchQueryService) ArchiveTrace(ctx context.Context, traceID model.TraceID) error {
	panic("implement me")
}

func (e ElasticsearchQueryService) FindTraces(ctx context.Context, queryParams service.TraceQueryParameters) ([]model.Trace, error) {
	panic("implement me")
}

func (e ElasticsearchQueryService) GetServices(ctx context.Context) {
	panic("implement me")
}

func (e ElasticsearchQueryService) GetOperations(ctx context.Context) {
	panic("implement me")
}

func (e ElasticsearchQueryService) GetDependencies(ctx context.Context, endTs time.Time, lookback time.Duration) {
	panic("implement me")
}
