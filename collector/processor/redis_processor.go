package processor

import (
	"collector/drive"
	"context"
	"github.com/jaegertracing/jaeger/model"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type redisProcessor struct {
}

func (sp redisProcessor) ProcessSpans(mSpans []*model.Span, options SpansOptions) ([]bool, error) {
	mongoOption := drive.MongoDriveOption{
		URI: "",
	}
	client, err := mongoOption.NewClient()
	if err != nil {
		return nil, err
	}
	collection := client.Database("trace").Collection("trace")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection.InsertOne(ctx, bson.D.Map(mSpans))
	return nil, nil
}
func (sp redisProcessor) Close() error {

	//todo
	return nil
}
