package es

type Span struct {
	TraceID         TraceID                `json:"traceID"`
	SpanID          SpanID                 `json:"spanID"`
	ParentSpanID    SpanID                 `json:"parentSpanID,omitempty"` // deprecated
	Flags           uint32                 `json:"flags,omitempty"`
	OperationName   string                 `json:"operationName"`
	References      []Reference            `json:"references"`
	StartTime       uint64                 `json:"startTime"` // microseconds since Unix epoch
	StartTimeMillis uint64                 `json:"startTimeMillis"`
	Duration        uint64                 `json:"duration"` // microseconds
	Tags            []KeyValue             `json:"tags"`
	Tag             map[string]interface{} `json:"tag,omitempty"`
	Logs            []Log                  `json:"logs"`
	Process         Process                `json:"process,omitempty"`
}

const (
	// ChildOf means a span is the child of another span
	ChildOf ReferenceType = "CHILD_OF"
	// FollowsFrom means a span follows from another span
	FollowsFrom ReferenceType = "FOLLOWS_FROM"
)

type Process struct {
	ServiceName string                 `json:"serviceName"`
	Tags        []KeyValue             `json:"tags"`
	Tag         map[string]interface{} `json:"tag,omitempty"`
}
type Log struct {
	Timestamp uint64     `json:"timestamp"`
	Fields    []KeyValue `json:"fields"`
}

type KeyValue struct {
	Key   string      `json:"key"`
	Type  ValueType   `json:"type"`
	Value interface{} `json:"value"`
}

type Reference struct {
	RefType ReferenceType `json:"refType"`
	TraceID TraceID       `json:"traceID"`
	SpanID  SpanID        `json:"spanID"`
}

type ReferenceType string

type TraceID string

type SpanID string

type ValueType string
