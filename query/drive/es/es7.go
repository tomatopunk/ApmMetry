package es

import (
	"context"
	"encoding/json"
	es7 "github.com/elastic/go-elasticsearch/v7"
	"io/ioutil"
	"net/http"
)

type es7client struct {
	client *es7.Client
}

const (
	bulkES7MetaFormat = `{"index":{"_index":"%s"}}` + "\n"
)

func NewElasticsearch7Client(config ClientConfig, roundTripper http.RoundTripper) (*es7client, error) {
	client, err := es7.NewClient(es7.Config{
		Addresses: config.Addresses,
		Username:  config.Username,
		Password:  config.Password,
		Transport: roundTripper,
	})
	if err != nil {
		return nil, err
	}
	return &es7client{
		client: client,
	}, nil
}

func (e es7client) Search(ctx context.Context, queries SearchBody, size int, indices ...string) (*SearchResponse, error) {
	body, err := queries.encodeSearchBody()
	if err != nil {
		return nil, err
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex(indices...),
		e.client.Search.WithBody(body),
		e.client.Search.WithIgnoreUnavailable(true),
		e.client.Search.WithSize(size))

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	r := &es7searchResponse{}
	if err = json.Unmarshal(data, r); err != nil {
		return nil, err
	}
	return &SearchResponse{
		Aggs: r.Aggs,
		Hits: Hits{
			Total: r.Hits.Total.Value,
			Hits:  r.Hits.Hits,
		},
	}, nil
}
