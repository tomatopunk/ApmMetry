package processor

import (
	"github.com/jaegertracing/jaeger/model"
	"io"
)

type SpansOptions struct {
	SpanFormat     SpanFormat
	TransportProto TransportProto
}

type TransportProto string

const (
	GRPCProto    TransportProto = "grpc"
	HttpProto    TransportProto = "http"
	UnknownProto TransportProto = "unknow"
)

type SpanFormat string

const (
	JaegerSpanFormat SpanFormat = "jaeger"
	ProtoSpanFormat  SpanFormat = "proto"
)

type SpanProcess interface {
	ProcessSpans(mSpans []*model.Span, options SpansOptions) ([]bool, error)
	io.Closer
}
