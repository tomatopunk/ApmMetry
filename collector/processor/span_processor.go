package processor

func NewSpanProcessor() SpanProcess {
	sp := newEsSpanProcessor()
	return sp
}

func newSpanProcessor() *grpcProcessor {
	sp := grpcProcessor{}
	return &sp
}

func newRedisSpanProcessor() *redisProcessor {
	sp := redisProcessor{}
	return &sp
}

func newMongoSpanProcessor() *MongoProcessor {
	sp := MongoProcessor{}
	return &sp
}

func newEsSpanProcessor() *esProcessor {
	sp := esProcessor{}
	return &sp
}
