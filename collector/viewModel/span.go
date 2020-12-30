package viewModel

type Span struct {
	traceId      string `json:"traceId""`
	spanName     string `json:"spanName"`
	spanId       string `json:"spanId"`
	spanParentId string `json:"spanParentId"`
}
