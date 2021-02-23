package spanStorage

import (
	"collector/drive/es"
	"collector/storage/es/model"
	"context"
	jModel "github.com/jaegertracing/jaeger/model"
	"net/http"
)

type EsStorage struct {
	client es.ElasticsearchClient
	config es.ClientConfig
	w      esSpanWriter
}

func NewEsStorage(esConfig es.ClientConfig) (EsStorage, error) {
	var storage = EsStorage{}
	httpTransport := &http.Transport{}
	client, err := es.NewClient(esConfig, httpTransport)

	if err != nil {
		return storage, err
	}
	storage.client = client
	storage.config = esConfig
	storage.w = esSpanWriter{
		client:    client,
		indexName: "span",
	}
	return storage, nil
}

func (s EsStorage) WriteSpan(ctx context.Context, m *jModel.Span) error {
	dbSpan := ConvertDbSpan(m)
	_, err := s.w.writeSpans(ctx, []model.Span{dbSpan})
	return err
}

func (s EsStorage) WriteSpans(ctx context.Context, m []*jModel.Span) error {
	dbSpans := ConvertDbSpans(m)
	_, err := s.w.writeSpans(ctx, dbSpans)
	return err
}
