package es

import (
	"bytes"
	"collector/storage/es/model"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

const mapping = "{\"index_patterns\":\"*span-*\",\"settings\":{\"index.number_of_shards\":3,\"index.number_of_replicas\":1,\"index.mapping.nested_fields.limit\":50,\"index.requests.cache.enable\":true},\"mappings\":{\"dynamic_templates\":[{\"span_tags_map\":{\"mapping\":{\"type\":\"keyword\",\"ignore_above\":256},\"path_match\":\"tag.*\"}},{\"process_tags_map\":{\"mapping\":{\"type\":\"keyword\",\"ignore_above\":256},\"path_match\":\"process.tag.*\"}}],\"properties\":{\"traceID\":{\"type\":\"keyword\",\"ignore_above\":256},\"parentSpanID\":{\"type\":\"keyword\",\"ignore_above\":256},\"spanID\":{\"type\":\"keyword\",\"ignore_above\":256},\"operationName\":{\"type\":\"keyword\",\"ignore_above\":256},\"startTime\":{\"type\":\"long\"},\"startTimeMillis\":{\"type\":\"date\",\"format\":\"epoch_millis\"},\"duration\":{\"type\":\"long\"},\"flags\":{\"type\":\"integer\"},\"logs\":{\"type\":\"nested\",\"dynamic\":false,\"properties\":{\"timestamp\":{\"type\":\"long\"},\"fields\":{\"type\":\"nested\",\"dynamic\":false,\"properties\":{\"key\":{\"type\":\"keyword\",\"ignore_above\":256},\"value\":{\"type\":\"keyword\",\"ignore_above\":256},\"tagType\":{\"type\":\"keyword\",\"ignore_above\":256}}}}},\"process\":{\"properties\":{\"serviceName\":{\"type\":\"keyword\",\"ignore_above\":256},\"tag\":{\"type\":\"object\"},\"tags\":{\"type\":\"nested\",\"dynamic\":false,\"properties\":{\"key\":{\"type\":\"keyword\",\"ignore_above\":256},\"value\":{\"type\":\"keyword\",\"ignore_above\":256},\"tagType\":{\"type\":\"keyword\",\"ignore_above\":256}}}}},\"references\":{\"type\":\"nested\",\"dynamic\":false,\"properties\":{\"refType\":{\"type\":\"keyword\",\"ignore_above\":256},\"traceID\":{\"type\":\"keyword\",\"ignore_above\":256},\"spanID\":{\"type\":\"keyword\",\"ignore_above\":256}}},\"tag\":{\"type\":\"object\"},\"tags\":{\"type\":\"nested\",\"dynamic\":false,\"properties\":{\"key\":{\"type\":\"keyword\",\"ignore_above\":256},\"value\":{\"type\":\"keyword\",\"ignore_above\":256},\"tagType\":{\"type\":\"keyword\",\"ignore_above\":256}}}}}}"

func TestCreateIndex(t *testing.T) {
	client, err := createClient()
	if err != nil {
		t.Error(err)
	}
	err = client.PutTemplate(context.TODO(), "span", strings.NewReader(mapping))
	if err != nil {
		t.Error(err)
	}
}

func TestSearch(t *testing.T) {
	client, err := createClient()
	if err != nil {
		t.Error(err)
	}

	terms := Terms{}
	terms["traceID"] = TermQuery{
		Value: "hahahah",
	}

	body := SearchBody{
		Query: &Query{
			Term: &terms,
		},
		Size: 100,
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(body)
	test := string(buf.Bytes())

	t.Log(test)

	res, err := client.Search(context.TODO(), body, 100, "span")

	if err != nil {
		t.Error(err)
	}

	t.Log(res)
}

func TestBulk(t *testing.T) {
	client, err := createClient()
	if err != nil {
		t.Error(err)
	}
	dates := make([]model.Span, 0, 100)
	for i := 0; i <= 100; i++ {
		dates = append(dates, model.Span{
			TraceID: "hahahah",
		})
	}
	buf := get_buf(client, dates)

	test := string(buf.Bytes())

	t.Log(test)

	res, err := client.Bulk(context.TODO(), buf)
	if err != nil {
		t.Error(err)
	}

	if res.Errors {
		t.Error(res)
	}
	t.Log(res)
}

func get_buf(client *es7client, dates []model.Span) *bytes.Buffer {
	buf := &bytes.Buffer{}
	var errs []error

	for _, span := range dates {
		data, err := json.Marshal(span)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		client.AddDataToBulkBuffer(buf, data, "span", "")
	}
	return buf
}

func createClient() (*es7client, error) {
	httpTransport := &http.Transport{}
	config := ClientConfig{
		Addresses: []string{"http://localhost:9200/"},
	}
	client, err := NewElasticsearch7Client(config, httpTransport)
	return client, err
}
