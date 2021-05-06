package spanStorage

import (
	"bytes"
	"collector/drive/es"
	"collector/storage/es/model"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type esSpanWriter struct {
	client    es.ElasticsearchClient
	indexName string
}

func (w esSpanWriter) writeSpans(ctx context.Context, span []model.Span) (int, error) {

	buf := &bytes.Buffer{}
	var errs []error

	for _, span := range span {
		data, err := json.Marshal(span)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		w.client.AddDataToBulkBuffer(buf, data, w.indexName, "span")
	}

	res, err := w.client.Bulk(ctx, buf)
	if err != nil {
		errs = append(errs, err)
		return 0, CombineErrors(errs)
	}
	errLen, err := handleBulkResponse(res)
	if err != nil {
		errs = append(errs, err)

		println("Batch is failed,len:{%s},reason:{%s}", errLen, err.Error())
	}
	return len(span), CombineErrors(errs)
}

func handleBulkResponse(res *es.BulkResponse) (int, error) {
	var errs []error
	for _, item := range res.Items {
		if item.Index.Status > 201 {
			errs = append(errs, fmt.Errorf("Failded reason:{%v},result:{%v}", item.Index.Error.Reason, item.Index.Result))
		}
	}
	return len(errs), CombineErrors(errs)
}

func CombineErrors(errs []error) error {
	num := len(errs)
	if num == 0 {
		return nil
	}

	if num == 1 {
		return errs[0]
	}

	errMsgs := make([]string, 0, num)

	for i, err := range errs {
		errMsgs[i] = err.Error()
	}
	err := fmt.Errorf("[%s]", strings.Join(errMsgs, "; "))
	return err
}
