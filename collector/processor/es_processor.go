package processor

import (
	"collector/storage/es/spanStorage"
	"context"
	"github.com/jaegertracing/jaeger/model"
	"time"
)

type esProcessor struct {
	processTimeout time.Duration
	storage        spanStorage.EsStorage
}

func (e esProcessor) ProcessSpans(mSpans []*model.Span, options SpansOptions) ([]bool, error) {
	res := make([]bool, len(mSpans))

	ctx, cancel := context.WithTimeout(context.Background(), e.processTimeout)

	defer cancel()

	err := e.storage.WriteSpans(ctx, mSpans)
	if err != nil {
		println(err)
	}
	return res, nil
}

func (e esProcessor) Close() error {
	panic("implement me")
}
