package produce

type Produce interface {
	SendMessage() error
}

type Span struct {
	TraceId      string `json:"traceId"`
	SpanName     string `json:"spanName"`
	SpanId       string `json:"spanId"`
	SpanParentId string `json:"spanParentId"`
}
