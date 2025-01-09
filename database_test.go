package database

import (
	"reflect"
	"testing"
)

func TestEncodeModel(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    M
		wantErr bool
	}{
		{
			name: "Test Empty Struct Tags",
			args: args{
				obj: struct {
					Name string
				}{
					Name: "Jon",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test Invalid Struct Tags",
			args: args{
				obj: struct {
					Name string `tag:"name"`
				}{
					Name: "Jon",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test Valid Struct Tags",
			args: args{
				obj: struct {
					Name string `db:"name,required,unique,index"`
				}{
					Name: "Jon",
				},
			},
			want: M{
				{
					Key:      "name",
					Value:    "Jon",
					Required: true,
					Unique:   true,
					Index:    true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeModel(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeModel(t *testing.T) {
	type user struct {
		Jon string `db:"jon"`
	}
	var res user
	type args struct {
		obj  interface{}
		data M
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test Non Pointer Struct",
			args: args{
				obj:  res,
				data: M{P{Key: "jon", Value: "jon"}}},
			wantErr: true,
		},
		{
			name: "Test Pointer Struct",
			args: args{
				obj:  &res,
				data: M{P{Key: "jon", Value: "jon"}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DecodeModel(tt.args.obj, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("DecodeModel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkTag(t *testing.T) {
	type args struct {
		fieldTags []string
		tag       string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Check Required Tag",
			args: args{
				fieldTags: []string{
					propertyRequired,
					propertyIndex,
					propertyUnique,
				},
				tag: propertyRequired,
			},
			want: true,
		},
		{
			name: "Check Index Tag",
			args: args{
				fieldTags: []string{
					propertyRequired,
					propertyIndex,
					propertyUnique,
				},
				tag: propertyIndex,
			},
			want: true,
		},
		{
			name: "Check Unique Tag",
			args: args{
				fieldTags: []string{
					propertyRequired,
					propertyIndex,
					propertyUnique,
				},
				tag: propertyUnique,
			},
			want: true,
		},
		{
			name: "Check Invalid Tag",
			args: args{
				fieldTags: []string{
					propertyRequired,
					propertyIndex,
					propertyUnique,
				},
				tag: "invalid",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkTag(tt.args.fieldTags, tt.args.tag); got != tt.want {
				t.Errorf("checkTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFieldname(t *testing.T) {
	type args struct {
		fieldTags []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Check Required Field",
			args: args{
				fieldTags: []string{
					propertyRequired,
					propertyIndex,
					propertyUnique,
				},
			},
			want: "",
		},
		{
			name: "Check Index Field",
			args: args{
				fieldTags: []string{
					propertyRequired,
					propertyIndex,
					propertyUnique,
				},
			},
			want: "",
		},
		{
			name: "Check Unique Field",
			args: args{
				fieldTags: []string{
					propertyRequired,
					propertyIndex,
					propertyUnique,
				},
			},
			want: "",
		},
		{
			name: "Check Valid Field",
			args: args{
				fieldTags: []string{
					propertyRequired,
					propertyIndex,
					propertyUnique,
					"name",
				},
			},
			want: "name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFieldname(tt.args.fieldTags); got != tt.want {
				t.Errorf("getFieldname() = %v, want %v", got, tt.want)
			}
		})
	}
}
