package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jaegertracing/jaeger/model"
	"query/drive/es"
	"query/service"
	"time"
)

type ElasticsearchQueryService struct {
	client es.ElasticsearchClient
}

const (
	searchIndex string = "span"
)

func (e ElasticsearchQueryService) GetTrace(ctx context.Context, traceID model.TraceID) (*model.Trace, error) {

	terms := es.Terms{}
	terms["traceID"] = es.TermQuery{
		Value: traceID.String(),
	}

	searchBody := es.SearchBody{
		Indices:      []string{searchIndex},
		Aggregations: nil,
		Query: &es.Query{
			Term: &terms,
		},
		Sort:           nil,
		Size:           0,
		TerminateAfter: 0,
		SearchAfter:    nil,
	}

	response, err := e.client.Search(ctx, searchBody, 1000, searchIndex)
	if err != nil {
		return nil, err
	}

	if len(response.Hits.Hits) == 0 {
		return nil, fmt.Errorf("es result is failed")
	}

	spans, err := e.collectSpans(response.Hits.Hits)
	if err != nil {
		return nil, fmt.Errorf("convert spans failed :{%s}", err)
	}

}

func (e ElasticsearchQueryService) collectSpans(hits []es.Hit) ([]*Span, error) {
	spans := make([]*Span, len(hits))
	for i, hit := range hits {
		esSpanInByteArr := hit.Source

		var jsonSpan Span

		d := json.NewDecoder(bytes.NewReader(*esSpanInByteArr))
		d.UseNumber()
		err := d.Decode(&jsonSpan)

		if err != nil {
			return nil, err
		}
		spans[i] = &jsonSpan
	}
	return spans, nil
}

func (e ElasticsearchQueryService) ArchiveTrace(ctx context.Context, traceID model.TraceID) error {
	panic("implement me")
}

func (e ElasticsearchQueryService) FindTraces(ctx context.Context, queryParams service.TraceQueryParameters) ([]model.Trace, error) {
	panic("implement me")
}

func (e ElasticsearchQueryService) GetServices(ctx context.Context) {

}

func (e ElasticsearchQueryService) GetOperations(ctx context.Context) {
	panic("implement me")
}

func (e ElasticsearchQueryService) GetDependencies(ctx context.Context, endTs time.Time, lookback time.Duration) {
	panic("implement me")
}
