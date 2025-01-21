package mongodb

import (
	"context"
	"errors"
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

// All implements database.Query.
func (m *MongoModel[T]) All() ([]*T, error) {
	var res []*T
	result, err := m.client.Find(m.ctx, m.filter, options.Find().
		SetLimit(m.limit).SetSkip(m.offset).SetSort(m.order))
	if err != nil {
		return nil, err
	}

	defer result.Close(m.ctx)

	for result.Next(m.ctx) {
		var single bson.D
		if err := result.Decode(&single); err != nil {
			return nil, err
		}
		var singleRes T
		if err := convertFromBson(&singleRes, single); err != nil {
			return nil, err
		}
		res = append(res, &singleRes)
	}
	m.reset()
	return res, nil
}

// Count implements database.Query.
func (m *MongoModel[T]) Count() (int64, error) {
	count, err := m.client.CountDocuments(m.ctx, m.filter, options.Count().SetLimit(m.limit).
		SetSkip(m.offset))
	if err != nil {
		return 0, err
	}
	m.reset()
	return count, nil
}

// Delete implements database.Query.
func (m *MongoModel[T]) Delete() error {
	_, err := m.client.DeleteOne(m.ctx, m.filter)
	if err != nil {
		return err
	}
	m.reset()
	return nil
}

// DeleteMany implements database.Query.
func (m *MongoModel[T]) DeleteMany() error {
	_, err := m.client.DeleteMany(m.ctx, m.filter)
	if err != nil {
		return err
	}
	m.reset()
	return nil
}

// First implements database.Query.
func (m *MongoModel[T]) First() (*T, error) {
	var res T
	result := m.client.FindOne(m.ctx, m.filter, options.FindOne().
		SetSort(m.order))

	var single bson.D
	if err := result.Decode(&single); err != nil {
		return nil, err
	}
	if err := convertFromBson(&res, single); err != nil {
		return nil, err
	}
	m.reset()
	return &res, nil
}

// Update implements database.Query.
func (m *MongoModel[T]) Update(doc T) error {
	d, err := convertToBson(doc)
	if err != nil {
		return err
	}
	result, err := m.client.UpdateOne(m.ctx, m.filter, bson.D{{Key: "$set", Value: d}})
	if err != nil {
		return err
	}
	if result.MatchedCount < 0 {
		return errors.New("error updating document")
	}
	m.reset()
	return nil
}

// UpdateMany implements database.Query.
func (m *MongoModel[T]) UpdateMany(doc T) error {
	d, err := convertToBson(doc)
	if err != nil {
		return err
	}
	result, err := m.client.UpdateMany(m.ctx, m.filter, bson.D{{Key: "$set", Value: d}})
	if err != nil {
		return err
	}
	if result.MatchedCount < 0 {
		return errors.New("error updating documents")
	}
	m.reset()
	return nil
}

// ExecRaw implements database.Store.
func (m *MongoModel[T]) ExecRaw() error {
	panic("unimplemented")
}

// Query implements database.Store.
func (m *MongoModel[T]) Query(query_params ...database.Params) database.Query[T] {
	var q_params []database.QueryStruct

	for _, param := range query_params {
		q_params = append(q_params, param())
	}

	for _, qq := range q_params {
		switch qq.Key() {
		case database.QueryFilter:
			val, ok := qq.Value().(database.FilterStruct)
			if !ok {
				panic(errors.New("unsupported"))
			}
			m.filter = append(m.filter, bson.E{Key: val.Key(), Value: val.Value()})
		case database.QuerySort:
			var val int
			order_val, ok := qq.Value().(database.OrderStruct)
			if !ok {
				panic(errors.New("unsupported"))
			}
			switch order_val.Value() {
			case database.ASC:
				val = -1
			case database.DESC:
				val = 1
			default:
				panic(errors.New("unsupported"))
			}
			m.order = append(m.order, bson.E{Key: order_val.Key(), Value: val})
		case database.QueryLimit:
			val, ok := qq.Value().(int64)
			if !ok {
			}
			m.limit = val
		case database.QueryOffset:
			val, ok := qq.Value().(int64)
			if !ok {
			}
			m.offset = val
		default:
			panic(errors.New("unsupported"))
		}
	}
	return m
}

// Save implements database.Store.
func (m *MongoModel[T]) Save(doc ...T) error {
	for _, d := range doc {
		v, err := convertToBson(d)
		if err != nil {
			return err
		}
		_, err = m.client.InsertOne(m.ctx, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// WithContext implements database.Store.
func (m *MongoModel[T]) WithContext(ctx context.Context) database.Model[T] {
	m.ctx = ctx
	return m
}

func (m *MongoModel[T]) reset() {
	m.filter = bson.D{}
	m.limit = 0
	m.offset = 0
	m.order = bson.D{}
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
		order:  bson.D{},
		filter: bson.D{},
		limit:  0,
		offset: 0,
	}, nil
}
