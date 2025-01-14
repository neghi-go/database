package database

import (
	"reflect"
	"testing"
)

func Test_parser_Parse(t *testing.T) {
	testOne := struct {
		Name string `test-tag:"name"`
		Skip string `test-tag:"-"`
	}{
		Name: "Jon Doe",
	}
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []parsedField
		wantErr bool
	}{
		{
			name: "Check parser",
			args: args{
				val: testOne,
			},
			want: []parsedField{
				{
					fieldValue: reflect.ValueOf(testOne).Field(0),
					fieldName:  reflect.ValueOf(testOne).Type().Field(0).Name,
					fieldTag:   []string{"name"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse("test-tag", tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("parser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareParserResponse(got, tt.want) {
				t.Errorf("parser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Benchmark(b *testing.B) {
	type TestStruct struct {
		Name string `test:"name"`
		Age  int    `test:"age"`
	}
	user := TestStruct{Name: "jon", Age: 20}
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			_, err := parse("test", user)
			if err != nil {
				b.Error(err)
			}

		}
	}
}
