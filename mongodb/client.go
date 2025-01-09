package mongodb

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var mongo_sync sync.Once

var mongo_instance *mongo.Client

type mongoDatabase struct {
	db *mongo.Database
}

func New(url, db string) (*mongoDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mongo_sync.Do(func() {
		client, err := mongo.Connect(options.Client().ApplyURI(url))
		if err != nil {
			panic(err)
		}
		mongo_instance = client
	})
	if err := mongo_instance.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &mongoDatabase{
		db: mongo_instance.Database(db),
	}, nil
}
