package handle

import (
	"collector/processor"
	"collector/protoc/api_v2"
	"context"
)

type GRPCHandler struct {
	spanProcess processor.SpanProcess
}

func NewGRPCHandler(spanProcess processor.SpanProcess) *GRPCHandler {
	return &GRPCHandler{
		spanProcess: spanProcess,
	}
}

func (g *GRPCHandler) PostSpans(ctx context.Context, r *api_v2.PostSpansRequest) (*api_v2.PostSpansResponse, error) {
	for _, span := range r.GetBatch().Spans {
		if span.GetProcess() == nil {
			span.Process = r.GetBatch().Process
		}
	}

	_, err := g.spanProcess.ProcessSpans(r.GetBatch().Spans, processor.SpansOptions{
		SpanFormat:     processor.ProtoSpanFormat,
		TransportProto: processor.GRPCProto,
	})

	if err != nil {
		return nil, err
	}
	return &api_v2.PostSpansResponse{}, nil
}
