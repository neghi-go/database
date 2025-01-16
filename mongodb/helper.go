package mongodb

import (
	"errors"

	"github.com/google/uuid"
	"github.com/neghi-go/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func getIndexes[T any](model T) ([]mongo.IndexModel, error) {
	var res []mongo.IndexModel
	parsed, err := database.EncodeModel(model)
	if err != nil {
		return nil, err
	}
	for _, e := range parsed {
		if e.Index && e.MongoID {
			return nil, errors.New("setting mongoid already sets index")
		}
		if e.Index {
			res = append(res, mongo.IndexModel{
				Keys:    bson.D{{Key: e.Key, Value: 1}},
				Options: options.Index().SetUnique(e.Unique),
			})
		}
	}
	return res, nil
}

func convertToBson[T any](data T) (bson.D, error) {
	var res = bson.D{}
	parsed, err := database.EncodeModel(data)
	if err != nil {
		return nil, err
	}
	for _, e := range parsed {
		if e.Required && e.Value == nil {
			return nil, errors.New("field is required but not provided")
		}
		if e.MongoID && e.Value != "" {
			id, err := bson.ObjectIDFromHex(e.Value.(string))
			if err != nil {
				return nil, err
			}
			e.Value = id
		}

		res = append(res, bson.E{Key: e.Key, Value: e.Value})
	}
	return res, nil
}

func convertFromBson[T any](obj T, doc bson.D) error {
	var parserModel = database.M{}
	for _, d := range doc {
		var val interface{}
		switch d.Value.(type) {
		case bson.ObjectID:
			val = d.Value.(bson.ObjectID).Hex()
		case bson.DateTime:
			t := d.Value.(bson.DateTime).Time()
			val = t
		case bson.Binary:
			t, ok := d.Value.(bson.Binary)
			if !ok {
				return errors.New("invalid doc provided")
			}
			if t.Subtype == 4 {
				val, _ = uuid.FromBytes(t.Data)
			}
		default:
			val = d.Value
		}
		parserModel = append(parserModel, database.P{Key: d.Key, Value: val})
	}
	return database.DecodeModel(obj, parserModel)
}
