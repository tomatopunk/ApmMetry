package span

type SpanDetailViewModel struct {
	TraceId         string
	Duration        int64
	StartTimestamp  int64
	FinishTimestamp int64
}

type SpanViewModel struct {
	TraceId         string
	SpanId          string
	Duration        int64
	StartTimestamp  int64
	FinishTimestamp int64
	Children        []SpanViewModel
	ServiceName     string
}
