package drive

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDriveOption struct {
	URI string
	//log
}

func (driveOption *MongoDriveOption) NewClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(driveOption.URI))
	if err != nil {
		return nil, err
	}
	return client, err
}
