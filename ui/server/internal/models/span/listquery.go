package span

type ListQuery struct {
	TraceId        string
	ServiceName    string
	StartTimestamp int64
	EndTimestamp   int64
	Tags           string
	Limit          int32
}
