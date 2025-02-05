package database

import (
	"errors"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	databaseTag = "db"
)

const (
	propertyRequired = "required"
	propertyIndex    = "index"
	propertyUnique   = "unique"
	propertyMongoID  = "mongoid"
	propertyUUID     = "uuid"
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
		case reflect.Int8:
			val = int8(p.fieldValue.Int())
		case reflect.Int:
			val = int(p.fieldValue.Int())
		case reflect.Int16:
			val = int16(p.fieldValue.Int())
		case reflect.Int32:
			val = int32(p.fieldValue.Int())
		case reflect.Int64:
			val = p.fieldValue.Int()
		case reflect.Uint8:
			val = uint8(p.fieldValue.Uint())
		case reflect.Uint:
			val = uint(p.fieldValue.Uint())
		case reflect.Uint16:
			val = uint16(p.fieldValue.Uint())
		case reflect.Uint32:
			val = uint32(p.fieldValue.Uint())
		case reflect.Uint64:
			val = p.fieldValue.Uint()
		case reflect.Float32, reflect.Float64:
			val = p.fieldValue.Float()
		case reflect.Bool:
			val = p.fieldValue.Bool()
		case reflect.String:
			val = p.fieldValue.String()
		default:
			val = p.fieldValue.Interface()
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
					switch field.Type() {
					case reflect.TypeOf(uuid.UUID{}):
						if str, ok := d.Value.(uuid.UUID); ok {
							field.Set(reflect.ValueOf(str))
						}
					case reflect.TypeOf(time.Time{}):
						if str, ok := d.Value.(time.Time); ok {
							field.Set(reflect.ValueOf(str))
						}
					default:
						switch field.Kind() {
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							val, err := handleIntTypes(d.Value)
							if err != nil {
								return err
							}
							field.SetInt(val)
						case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
							val, err := handleUintTypes(d.Value)
							if err != nil {
								return err
							}
							field.SetUint(val)
						case reflect.Float32, reflect.Float64:
							val, err := handleFloatTypes(d.Value)
							if err != nil {
								return err
							}
							field.SetFloat(val)
						case reflect.String:
							val, ok := d.Value.(string)
							if !ok {
								return errors.New("not ok")
							}
							field.SetString(val)
						default:
							field.Set(reflect.ValueOf(d.Value))
						}
					}

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
		case propertyMongoID:
			return "_id"
		default:
			return tags
		}
	}
	return ""
}

func handleIntTypes(value interface{}) (int64, error) {
	switch value.(type) {
	case int:
		val, _ := value.(int)
		return int64(val), nil
	case int8:
		val, _ := value.(int8)
		return int64(val), nil
	case int16:
		val, _ := value.(int16)
		return int64(val), nil
	case int32:
		val, _ := value.(int32)
		return int64(val), nil
	case int64:
		val, _ := value.(int64)
		return val, nil
	default:
		return 0, errors.New("invalid type provided!")
	}
}

func handleUintTypes(value interface{}) (uint64, error) {
	switch value.(type) {
	case uint:
		val, _ := value.(uint)
		return uint64(val), nil
	case uint8:
		val, _ := value.(uint8)
		return uint64(val), nil
	case uint16:
		val, _ := value.(uint16)
		return uint64(val), nil
	case uint32:
		val, _ := value.(uint32)
		return uint64(val), nil
	case uint64:
		val, _ := value.(uint64)
		return val, nil
	default:
		return 0, errors.New("invalid type provided!")
	}
}

func handleFloatTypes(value interface{}) (float64, error) {
	switch value.(type) {
	case uint:
		val, _ := value.(float32)
		return float64(val), nil
	case float64:
		val, _ := value.(float64)
		return val, nil
	default:
		return 0, errors.New("invalid type provided!")
	}
}
