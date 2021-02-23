package processor

import (
	"collector/drive/es"
	"collector/storage/es/spanStorage"
	"time"
)

func NewSpanProcessor(config es.ClientConfig) (SpanProcess, error) {
	//todo Should by in assembler
	sp, err := newEsSpanProcessor(config)
	if err != nil {
		return nil, err
	}
	return sp, nil
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

func newEsSpanProcessor(esConfig es.ClientConfig) (*esProcessor, error) {
	//esConfig := es.ClientConfig{}
	storage, err := spanStorage.NewEsStorage(esConfig)
	if err != nil {
		return nil, err
	}
	sp := esProcessor{
		processTimeout: time.Second * 10,
		storage:        storage,
	}
	return &sp, nil
}
