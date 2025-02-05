package database

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

const (
	skipFieldTag = "-"
)

var (
	ErrNotStruct = errors.New("error: invalid type received. expected type %s, but got type %s")
)

var pool = sync.Pool{
	New: func() any {
		return &parsedField{
			fieldValue: reflect.Value{},
			fieldName:  "",
			fieldTag:   make([]string, 0),
		}
	},
}

// ParserField is contains information about the fields in the struct that was parsed,
// information pertaining to its value, type, name, structtags etc
type parsedField struct {
	fieldValue reflect.Value
	fieldName  string
	fieldTag   []string
}

func (p *parsedField) Release() {
	p.fieldName = ""
	p.fieldTag = make([]string, 0)
	p.fieldValue = reflect.Value{}
}

func parse(struct_tag string, obj interface{}) ([]parsedField, error) {
	sliceOfParsed := make([]parsedField, 0)
	v := reflect.ValueOf(obj)

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf(ErrNotStruct.Error(), reflect.Struct.String(), v.Kind().String())
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		single, _ := pool.Get().(*parsedField)
		single.fieldName = field.Name
		single.fieldValue = v.Field(i)

		tag := field.Tag
		def, ok := tag.Lookup(struct_tag)

		if !ok {
			return nil, errors.New("struct tag expected, got empty")
		}
		if strings.Contains(def, skipFieldTag) {
			continue
		}
		attr := strings.Split(def, ",")
		single.fieldTag = attr
		sliceOfParsed = append(sliceOfParsed, *single)
		single.Release()
		pool.Put(single)
	}
	return sliceOfParsed, nil
}

func compareParserResponse(a, b []parsedField) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].fieldName != b[i].fieldName ||
			!reflect.DeepEqual(a[i].fieldTag, b[i].fieldTag) ||
			!reflect.DeepEqual(a[i].fieldValue.Interface(), b[i].fieldValue.Interface()) {
			return false
		}
	}
	return true
}
