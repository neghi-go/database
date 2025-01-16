package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func initClient(url string) (*mongo.Client, error) {
	registry := mongoRegistry
	registry.RegisterTypeEncoder(tUUID, bson.ValueEncoderFunc(uuidEncodeValue))
	registry.RegisterTypeDecoder(tUUID, bson.ValueDecoderFunc(uuidDecodeValue))

	return mongo.Connect(options.Client().ApplyURI(url).SetRegistry(registry))
}

type mongoDatabase struct {
	db *mongo.Database
}

func New(url, db string) (*mongoDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	client, err := initClient(url)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return &mongoDatabase{
		db: client.Database(db),
	}, nil
}

func (m *mongoDatabase) Disconnect(ctx context.Context) error {
	return m.db.Client().Disconnect(ctx)
}
