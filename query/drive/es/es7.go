package es

import (
	"bytes"
	"context"
	"encoding/json"
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

func (e es7client) MultiSearch(ctx context.Context, queries []SearchBody) (*MultiSearchResponse, error) {
	var indices []string
	var es7Queries []es7SearchBody
	for _, q := range queries {
		es7Queries = append(es7Queries, es7SearchBody{
			SearchBody:     q,
			TrackTotalHits: true,
		})
		indices = append(indices, q.Indices...)
	}
	body, err := es7QueryBodies(es7Queries)
	if err != nil {
		return nil, err
	}

	response, err := e.client.Msearch(body,
		e.client.Msearch.WithContext(ctx),
		e.client.Msearch.WithIndex(indices...),
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	r := &es7multiSearchResponse{}
	if err = json.Unmarshal(data, r); err != nil {
		return nil, err
	}
	return convertMultiSearchResponse(r), nil
}

type es7SearchBody struct {
	SearchBody
	TrackTotalHits bool `json:"track_total_hits"`
}

func es7QueryBodies(searchBodies []es7SearchBody) (io.Reader, error) {
	buf := &bytes.Buffer{}
	for _, sb := range searchBodies {
		data, err := json.Marshal(sb)
		if err != nil {
			return nil, err
		}
		addDataToMSearchBuffer(buf, data)
	}
	return buf, nil
}

func addDataToMSearchBuffer(buffer *bytes.Buffer, data []byte) {
	meta := []byte(multiSearchHeaderFormat)
	buffer.Grow(len(data) + len(meta) + len("\n"))
	buffer.Write(meta)
	buffer.Write(data)
	buffer.Write([]byte("\n"))
}

type es7multiSearchResponse struct {
	Responses []es7searchResponse `json:"responses"`
}

const multiSearchHeaderFormat = `{"ignore_unavailable": "true"}` + "\n"

func convertMultiSearchResponse(response *es7multiSearchResponse) *MultiSearchResponse {
	mResponse := &MultiSearchResponse{}
	for _, r := range response.Responses {
		mResponse.Responses = append(mResponse.Responses, convertSearchResponse(r))
	}
	return mResponse
}

func convertSearchResponse(response es7searchResponse) SearchResponse {
	return SearchResponse{
		Aggs: response.Aggs,
		Hits: Hits{
			Total: response.Hits.Total.Value,
			Hits:  response.Hits.Hits,
		},
	}
}
