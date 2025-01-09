package database

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"strings"
)

const (
	databaseTag = "db"
)

const (
	propertyRequired = "required"
	propertyIndex    = "index"
	propertyUnique   = "unique"
	propertyMongoID  = "mongoid"
)

type P struct {
	Key      string
	Value    interface{}
	Required bool
	Unique   bool
	Index    bool
	MongoID  bool
}

type M []P

func EncodeModel(obj interface{}) (M, error) {
	var res M = M{}
	parsed, err := parse(databaseTag, obj)
	if err != nil {
		return nil, err
	}

	for _, p := range parsed {
		var val interface{}
		switch p.fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val = p.fieldValue.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val = p.fieldValue.Uint()
		case reflect.Float32, reflect.Float64:
			val = p.fieldValue.Float()
		case reflect.Bool:
			val = p.fieldValue.Bool()
		case reflect.String:
			val = p.fieldValue.String()
		default:
			return nil, errors.New("case not provided")
		}

		var key string
		switch checkTag(p.fieldTag, propertyMongoID) {
		case true:
			key = "_id"
		default:
			key = getFieldname(p.fieldTag)
		}

		res = append(res, P{Key: key, Value: val,
			Required: checkTag(p.fieldTag, propertyRequired),
			Index:    checkTag(p.fieldTag, propertyIndex),
			Unique:   checkTag(p.fieldTag, propertyUnique),
			MongoID:  checkTag(p.fieldTag, propertyMongoID),
		})

	}
	return res, nil
}

func DecodeModel(obj interface{}, data M) error {
	v := reflect.ValueOf(obj)

	if v.Kind() != reflect.Pointer {
		return errors.New("expect a pointer to a struct")
	}

	p := v.Elem()
	for i := 0; i < p.NumField(); i++ {
		field := p.Field(i)
		if field.CanSet() && field.IsValid() {
			for _, d := range data {
				tags := strings.Split(p.Type().Field(i).Tag.Get(databaseTag), ",")
				tag := getFieldname(tags)
				if tag == d.Key {
					field.Set(reflect.ValueOf(d.Value))
				}
			}
		}
	}

	return nil
}

func checkTag(fieldTags []string, tag string) bool {
	return slices.Contains(fieldTags, tag)
}

func getFieldname(fieldTags []string) string {
	for _, tags := range fieldTags {
		switch tags {
		case propertyIndex, propertyRequired, propertyUnique:
			continue
		case propertyMongoID:
			return "_id"
		default:
			return tags
		}
	}
	return ""
}

type Model[T any] interface {
	// MISC
	WithContext(ctx context.Context) Model[T]

	//SAVE
	Save(data ...T) error

	//QUERY
	Filter(filter D) Model[T]
	Limit(limit int64) Model[T]
	Offset(offset int64) Model[T]
	Order(order D) Model[T]
	Count() (count int64, err error)

	First() (*T, error)
	Find() ([]*T, error)

	//UPDATE
	UpdateOne(data T) error
	UpdateMany(data T) error

	//DELETE
	Delete() error
}
