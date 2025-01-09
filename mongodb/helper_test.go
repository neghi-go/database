package mongodb

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Test_getIndexes(t *testing.T) {
	type args struct {
		model interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []mongo.IndexModel
		wantErr bool
	}{
		{
			name: "Test With Index Property",
			args: args{
				model: struct {
					Name string `db:"name,index"`
				}{},
			},
			want: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "name", Value: 1}},
					Options: options.Index().SetUnique(false),
				},
			},
			wantErr: false,
		},
		{
			name: "Test With Index and Unique Property",
			args: args{
				model: struct {
					Name string `db:"name,index,unique"`
				}{},
			},
			want: []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "name", Value: 1}},
					Options: options.Index().SetUnique(true),
				},
			},
			wantErr: false,
		},
		{
			name: "Test with Keyed Property",
			args: args{
				model: struct {
					Name string `db:"name"`
				}{},
			},
			want:    []mongo.IndexModel{},
			wantErr: false,
		},
		{
			name: "Test with MongoID Property",
			args: args{
				model: struct {
					ID string `db:"index,mongoid"`
				}{},
			},
			want:    []mongo.IndexModel{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getIndexes(tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("getIndexes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.IsType(t, tt.want, got)
		})
	}
}

func Test_convertToBson(t *testing.T) {
	id, _ := bson.ObjectIDFromHex("677904ef31ac7ccf730d4e39")
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bson.D
		wantErr bool
	}{
		{
			name: "Test With Mongo ID",
			args: args{
				data: struct {
					ID   string `db:"mongoid"`
					Name string `db:"name"`
				}{
					ID:   "677904ef31ac7ccf730d4e39",
					Name: "Jon Doe",
				},
			},
			want: bson.D{
				{Key: "_id", Value: id},
				{Key: "name", Value: "Jon Doe"},
			},
			wantErr: false,
		},
		{
			name: "Test without MongoID",
			args: args{
				data: struct {
					Name string `db:"name"`
				}{
					Name: "Jon Doe",
				},
			},
			want: bson.D{
				{Key: "name", Value: "Jon Doe"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToBson(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToBson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToBson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertFromBson(t *testing.T) {
	id, _ := bson.ObjectIDFromHex("677904ef31ac7ccf730d4e39")
	res := struct {
		ID   string `db:"mongoid"`
		Name string `db:"name"`
	}{}
	type args struct {
		obj interface{}
		doc bson.D
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test One",
			args: args{
				obj: &res,
				doc: bson.D{
					{Key: "_id", Value: id},
					{Key: "name", Value: "Jon Doe"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := convertFromBson(tt.args.obj, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("convertFromBson() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(tt.args.obj)
		})
	}
}
