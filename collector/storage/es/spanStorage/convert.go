package spanStorage

import (
	"collector/storage/es/model"
	jModel "github.com/jaegertracing/jaeger/model"
	"strings"
)

func ConvertDbSpans(mspan []*jModel.Span) []model.Span {
	spans := make([]model.Span, 0, len(mspan))
	for _, s := range mspan {
		span := ConvertDbSpan(s)
		spans = append(spans, span)
	}
	return spans
}
func ConvertDbSpan(mspan *jModel.Span) model.Span {
	s := convertSpan(mspan)
	s.Process = convertProcess(mspan.Process)
	s.References = convertReferences(mspan)
	return s
}

func convertReferences(mspan *jModel.Span) []model.Reference {
	out := make([]model.Reference, 0, len(mspan.References))
	for _, ref := range mspan.References {
		out = append(out, model.Reference{
			RefType: convertRefType(ref.RefType),
			TraceID: model.TraceID(ref.TraceID.String()),
			SpanID:  model.SpanID(ref.SpanID.String()),
		})
	}
	return out
}

func convertRefType(refType jModel.SpanRefType) model.ReferenceType {
	if refType == jModel.FollowsFrom {
		return model.FollowsFrom
	}
	return model.ChildOf
}

func convertProcess(process *jModel.Process) model.Process {
	tags, tagsMap := convertKeyValueString(process.Tags)
	return model.Process{
		ServiceName: process.ServiceName,
		Tags:        tags,
		Tag:         tagsMap,
	}
}

func convertSpan(mspan *jModel.Span) model.Span {
	tags, tagsMap := convertKeyValueString(mspan.Tags)
	return model.Span{
		TraceID:         model.TraceID(mspan.TraceID.String()),
		SpanID:          model.SpanID(mspan.SpanID.String()),
		Flags:           uint32(mspan.Flags),
		OperationName:   mspan.OperationName,
		StartTime:       model.TimeAsEpochMicroseconds(mspan.StartTime),
		StartTimeMillis: model.TimeAsEpochMicroseconds(mspan.StartTime) / 1000,
		Duration:        model.DurationAsMicroseconds(mspan.Duration),
		Tags:            tags,
		Tag:             tagsMap,
		Logs:            convertLogs(mspan.Logs),
	}

}
func convertLogs(logs []jModel.Log) []model.Log {
	out := make([]model.Log, len(logs))
	for i, log := range logs {
		var kvs []model.KeyValue
		for _, kv := range log.Fields {
			kvs = append(kvs, convertKeyValue(kv))
		}
		out[i] = model.Log{
			Timestamp: model.TimeAsEpochMicroseconds(log.Timestamp),
			Fields:    kvs,
		}
	}
	return out
}
func convertKeyValueString(keyValues jModel.KeyValues) ([]model.KeyValue, map[string]interface{}) {
	var tagsMap = map[string]interface{}{}
	var kvs []model.KeyValue

	for _, kv := range keyValues {
		if kv.GetVType() != jModel.BinaryType {
			tagsMap[kv.Key] = kv.Value()
		} else {
			kvs = append(kvs, convertKeyValue(kv))
		}
	}
	return kvs, tagsMap

}

func convertKeyValue(kv jModel.KeyValue) model.KeyValue {
	return model.KeyValue{
		Key:   kv.Key,
		Type:  model.ValueType(strings.ToLower(kv.VType.String())),
		Value: kv.AsString(),
	}
}
