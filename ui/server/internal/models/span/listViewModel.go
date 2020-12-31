package span

type ListViewModel struct {
	TraceId        string
	Duration       int64
	StartTimestamp int64
	EndTimestamp   int64
	SpanName       string
}
