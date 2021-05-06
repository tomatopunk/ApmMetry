package model

type Span struct {
	TraceID       string
	SpanID        string
	OperationName string
	References    []Reference
	StartTime     uint64
	Duration      uint64
	Flags         uint64
	Tags          []KeyValue
	Tag           map[string]interface{}
	Logs          []Log
	Process       Process
}

type ReferenceType string

type ValueType string

const (
	// ChildOf means a span is the child of another span
	ChildOf ReferenceType = "CHILD_OF"
	// FollowsFrom means a span follows from another span
	FollowsFrom ReferenceType = "FOLLOWS_FROM"

	// StringType indicates a string value stored in KeyValue
	StringType ValueType = "string"
	// BoolType indicates a Boolean value stored in KeyValue
	BoolType ValueType = "bool"
	// Int64Type indicates a 64bit signed integer value stored in KeyValue
	Int64Type ValueType = "int64"
	// Float64Type indicates a 64bit float value stored in KeyValue
	Float64Type ValueType = "float64"
	// BinaryType indicates an arbitrary byte array stored in KeyValue
	BinaryType ValueType = "binary"
)

type Reference struct {
	RefType ReferenceType
	TraceID string
	SpanID  string
}

type KeyValue struct {
	Key   string
	Type  ValueType
	Value interface{}
}

type Log struct {
	Timestamp uint64
	Fields    []KeyValue
}

type Process struct {
	ServiceName string
	Tags        []KeyValue
	Tag         map[string]interface{}
}
