package es

import "github.com/jaegertracing/jaeger/model"

func (dbSpan Span) SpanToDomain() (*model.Span, error) {
	tags, err := dbSpan.convertKeyValues(dbSpan.Tags)
	if err != nil {
		return nil, err
	}
	logs, err := dbSpan.convertLogs(dbSpan.Logs)
	if err != nil {
		return nil, err
	}
	refs, err := dbSpan.convertRefs(dbSpan.References)
	if err != nil {
		return nil, err
	}
	process, err := dbSpan.convertProcess(dbSpan.Process)
	if err != nil {
		return nil, err
	}
	traceID, err := model.TraceIDFromString(string(dbSpan.TraceID))
	if err != nil {
		return nil, err
	}

	spanIDInt, err := model.SpanIDFromString(string(dbSpan.SpanID))
	if err != nil {
		return nil, err
	}

	if dbSpan.ParentSpanID != "" {
		parentSpanID, err := model.SpanIDFromString(string(dbSpan.ParentSpanID))
		if err != nil {
			return nil, err
		}
		refs = model.MaybeAddParentSpanID(traceID, parentSpanID, refs)
	}

	fieldTags, err := dbSpan.convertTagFields(dbSpan.Tag)
	if err != nil {
		return nil, err
	}
	tags = append(tags, fieldTags...)

	span := &model.Span{
		TraceID:       traceID,
		SpanID:        model.NewSpanID(uint64(spanIDInt)),
		OperationName: dbSpan.OperationName,
		References:    refs,
		Flags:         model.Flags(uint32(dbSpan.Flags)),
		StartTime:     model.EpochMicrosecondsAsTime(dbSpan.StartTime),
		Duration:      model.MicrosecondsAsDuration(dbSpan.Duration),
		Tags:          tags,
		Logs:          logs,
		Process:       process,
	}
	return span, nil
}

func (dbSpan span) convertKeyValues() {

}
