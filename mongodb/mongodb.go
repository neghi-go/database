package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/neghi-go/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoModel[T any] struct {
	ctx    context.Context
	filter bson.D
	order  bson.D
	limit  int64
	offset int64
	client *mongo.Collection
}

// Count implements database.Model.
func (m *MongoModel[T]) Count() (count int64, err error) {
	return m.client.CountDocuments(m.ctx, nil)

}

// Delete implements database.Model.
func (m *MongoModel[T]) Delete() error {
	_, err := m.client.DeleteOne(m.ctx, m.filter)
	if err != nil {
		return err
	}
	return nil
}

// Filter implements database.Model.
func (m *MongoModel[T]) Filter(filter database.D) database.Model[T] {
	for _, f := range filter {
		m.filter = append(m.filter, bson.E{Key: f.Key(), Value: f.Value()})
	}
	return m
}

// Find implements database.Model.
func (m *MongoModel[T]) Find() ([]*T, error) {
	var res []*T
	result, err := m.client.Find(m.ctx, m.filter, options.Find().
		SetSort(m.order).SetLimit(m.limit).SetSkip(m.offset))
	if err != nil {
		return nil, err
	}
	result.Close(m.ctx)
	for result.Next(m.ctx) {
		var singleDoc bson.D
		err := result.Decode(&singleDoc)
		if err != nil {
			return nil, err
		}

		var singleModel T
		err = convertFromBson(&singleModel, singleDoc)
		if err != nil {
			return nil, err
		}

		res = append(res, &singleModel)
	}

	return res, nil
}

// First implements database.Model.
func (m *MongoModel[T]) First() (*T, error) {
	var singleModel T
	result := m.client.FindOne(m.ctx, m.filter, options.FindOne().
		SetSort(m.order).SetSkip(m.offset))
	var singleDoc bson.D
	err := result.Decode(&singleDoc)
	if err != nil {
		return nil, err
	}

	err = convertFromBson(&singleModel, singleDoc)
	if err != nil {
		return nil, err
	}
	return &singleModel, nil
}

// Limit implements database.Model.
func (m *MongoModel[T]) Limit(limit int64) database.Model[T] {
	m.limit = limit
	return m
}

// Offset implements database.Model.
func (m *MongoModel[T]) Offset(offset int64) database.Model[T] {
	m.offset = offset
	return m
}

// Order implements database.Model.
func (m *MongoModel[T]) Order(order database.D) database.Model[T] {
	for _, o := range order {
		var val interface{}
		sortKey, ok := o.Value().(database.OrderType)
		if !ok {
			fmt.Println("error: Invalid Order value provided")
			break
		}
		switch sortKey {
		case database.ASC:
			val = -1
		case database.DESC:
			val = 1
		default:
			fmt.Println("unsupported Order Key type")
			return nil
		}
		m.order = append(m.order, bson.E{Key: o.Key(), Value: val})
	}
	return m
}

// Save implements database.Model.
func (m *MongoModel[T]) Save(data ...T) error {
	for _, d := range data {
		doc, err := convertToBson(d)
		if err != nil {
			return err
		}
		_, err = m.client.InsertOne(m.ctx, doc)
		if err != nil {
			return err
		}
	}
	return nil
}

// Update implements database.Model.
func (m *MongoModel[T]) UpdateOne(data T) error {
	doc, err := convertToBson(data)
	if err != nil {
		return err
	}
	_, err = m.client.UpdateOne(m.ctx, m.filter, bson.D{{Key: "$set", Value: doc}})
	if err != nil {
		return err
	}

	return nil
}

// Update implements database.Model.
func (m *MongoModel[T]) UpdateMany(data T) error {
	doc, err := convertToBson(data)
	if err != nil {
		return err
	}
	_, err = m.client.UpdateMany(m.ctx, m.filter, bson.D{{Key: "$set", Value: doc}})
	if err != nil {
		return err
	}

	return nil
}

// WithContext implements database.Model.
func (m *MongoModel[T]) WithContext(ctx context.Context) database.Model[T] {
	m.ctx = ctx
	return m
}

func RegisterModel[T any](conn *mongoDatabase, coll string, model T) (database.Model[T], error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	col := conn.db.Collection(coll)

	indexes, err := getIndexes(model)
	if err != nil {
		return nil, err
	}
	_, err = col.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return nil, err
	}

	return &MongoModel[T]{
		client: col,
		ctx:    context.Background(),
	}, nil
}
