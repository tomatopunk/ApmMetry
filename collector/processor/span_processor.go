package processor


func NewSpanProcessor() SpanProcess{
	//sp := newSpanProcessor()
	//return sp
}

func newSpanProcessor() *grpcProcessor{
	sp := grpcProcessor{
	}
	return &sp
}