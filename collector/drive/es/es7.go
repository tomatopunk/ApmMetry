package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	es7 "github.com/elastic/go-elasticsearch/v7"
	"io"
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

func (e es7client) Search(ctx context.Context, queries []SearchBody, size int, indices ...string) (*SearchResponse, error) {
	panic("implement me")
}
