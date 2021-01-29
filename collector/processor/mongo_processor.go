package processor

import (
	"github.com/jaegertracing/jaeger/model"
)

type MongoProcessor struct {
}

func (sp MongoProcessor) ProcessSpans(mSpans []*model.Span, options SpansOptions) ([]bool, error) {
	panic("implement me")
}

func (sp MongoProcessor) Close() error {
	panic("implement me")
}
