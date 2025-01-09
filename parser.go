// package Siddon contains the implementation of the parser that marshals/converts a struct to a readable format
// that other pacakges can use

// it expects the struct to have a tag of  `db:"<field-name>"` and or `attr:"<attribute-name>"` in order to
// parse them correctly.

// Example:

// type User struct {
// 	ID      string    `db:"id" attr:"required,mongoid"`
// 	Name    string    `db:"name" attr:"min=2,max=20"`
// 	Email   string    `db:"email" attr:"email,required"`
// 	Age     int       `db:"age" attr:"required"`
// 	DOB     time.Time `db:"date_of_birth" attr:"required,tz"`
// 	Balance float32   `db:"balance"`
// 	Street  Address   `attr:"embed"`
// 	Friends []Friends `db:"friends"`
// }

// type Address struct {
// 	Street  string `db:"street"`
// 	State   string `db:"state"`
// 	ZipCode string `db:"zip_code"`
// 	Country string `db:"country"`
// }

// type Friends struct {
// 	Name  string `db:"name"`
// 	Email string `db:"email" attr:"email"`
// 	Phone uint   `db:"phone" attr:"required"`
// }

package database

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	skipFieldTag = "-"
)

var (
	ErrNotStruct = errors.New("error: invalid type received. expected type %s, but got type %s")
)

// ParserField is contains information about the fields in the struct that was parsed,
// information pertaining to its value, type, name, structtags etc
type parsedField struct {
	fieldValue reflect.Value
	fieldName  string
	fieldTag   []string
}

func parse(struct_tag string, obj interface{}) ([]parsedField, error) {
	var sliceOfParsed []parsedField
	v := reflect.ValueOf(obj)

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf(ErrNotStruct.Error(), reflect.Struct.String(), v.Kind().String())
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		var single parsedField
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
		sliceOfParsed = append(sliceOfParsed, single)
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
