package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	es7 "github.com/elastic/go-elasticsearch/v7"
	"io"
	"io/ioutil"
	"net/http"
)

type es7client struct {
	client *es7.Client
}

const (
	bulkES7MetaFormat = `{"index":{"_index":"%s"}}` + "\n"
)

func NewElasticsearch7Client(config clientConfig, roundTripper http.RoundTripper) (*es7client, error) {
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

func (e es7client) PutTemplate(ctx context.Context, name string, template io.Reader) error {
	response, err := e.client.Indices.PutTemplate(name, template, e.client.Indices.PutTemplate.WithContext(ctx))
	if err != nil {
		return err
	}
	response.Body.Close()
	return nil
}

func (e es7client) AddDataToBulkBuffer(bulkBody *bytes.Buffer, data []byte, index, typ string) {
	meta := []byte(fmt.Sprintf(bulkES7MetaFormat, index))
	bulkBody.Grow(len(meta) + len(data) + len("\n"))
	bulkBody.Write(meta)
	bulkBody.Write(data)
	bulkBody.Write([]byte("\n"))
}

func (e es7client) Bulk(ctx context.Context, bulkBody io.Reader) (*BulkResponse, error) {
	response, err := e.client.Bulk(bulkBody, e.client.Bulk.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("bulk failed with code %d", response.StatusCode)
	}
	var blk BulkResponse
	err = json.NewDecoder(response.Body).Decode(&blk)
	if err != nil {
		return nil, err
	}
	return &blk, nil
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
