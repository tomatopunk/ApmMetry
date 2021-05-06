package service

import (
	"context"
	"github.com/jaegertracing/jaeger/model"
	"query/service/es"
	"time"
)

type QueryService interface {
	GetTrace(ctx context.Context, traceID model.TraceID) (*model.Trace, error)
	ArchiveTrace(ctx context.Context, traceID model.TraceID) error
	FindTraces(ctx context.Context, queryParams TraceQueryParameters) ([]model.Trace, error)
	GetServices(ctx context.Context)
	GetOperations(ctx context.Context)
	GetDependencies(ctx context.Context, endTs time.Time, lookback time.Duration)
}

func NewQueryService() QueryService {
	return &es.ElasticsearchQueryService{}
}
