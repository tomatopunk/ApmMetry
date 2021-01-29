package processor

import "github.com/jaegertracing/jaeger/model"

type grpcProcessor struct {
}

func (sp *grpcProcessor) Close() error {
	return nil
}

func (sp *grpcProcessor) ProcessSpans(mSpans []*model.Span, options SpansOptions) ([]bool, error) {
	return nil, nil
}
