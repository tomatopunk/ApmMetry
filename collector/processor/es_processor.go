package processor

import "github.com/jaegertracing/jaeger/model"

type esProcessor struct {
}

func (e esProcessor) ProcessSpans(mSpans []*model.Span, options SpansOptions) ([]bool, error) {
	panic("implement me")
}

func (e esProcessor) Close() error {
	panic("implement me")
}
