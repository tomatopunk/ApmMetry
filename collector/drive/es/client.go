package es

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type ClientConfig struct {
	Addresses []string
	Username  string
	Password  string
}

type ElasticsearchClient interface {
	PutTemplate(ctx context.Context, name string, template io.Reader) error

	Bulk(ctx context.Context, bulkBody io.Reader) (*BulkResponse, error)

	Search(ctx context.Context, query SearchBody, size int, indices ...string) (*SearchResponse, error)

	AddDataToBulkBuffer(bulkBody *bytes.Buffer, data []byte, index, typ string)
}

type BulkResponse struct {
	Errors bool               `json:"errors"`
	Items  []BulkResponseItem `json:"items"`
}

type BulkResponseItem struct {
	Index BulkIndexResponse `json:"index"`
}

type BulkIndexResponse struct {
	ID     string `json:"_id"`
	Result string `json:"result"`
	Status int    `json:"status"`
	Error  struct {
		Type   string `json:"type"`
		Reason string `json:"reason"`
		Cause  struct {
			Type   string `json:"type"`
			Reason string `json:"reason"`
		} `json:"caused_by"`
	} `json:"error"`
}

type es7searchResponse struct {
	Hits es7its                         `json:"hits"`
	Aggs map[string]AggregationResponse `json:"aggregations,omitempty"`
}

type es7its struct {
	Total struct {
		Value int `json:"value"`
	} `json:"total"`
	Hits []Hit `json:"hits"`
}

func NewClient(config ClientConfig, roundTripper http.RoundTripper) (ElasticsearchClient, error) {
	client, err := NewElasticsearch7Client(config, roundTripper)
	if err != nil {
		return nil, err
	}
	return client, nil
}
