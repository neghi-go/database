package mongodb

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/neghi-go/database"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

var test_url string

func TestMain(m *testing.M) {
	client := testcontainers.ContainerRequest{
		Image:        "mongo:8.0",
		ExposedPorts: []string{"27017/tcp"},
	}
	mongoClient, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: client,
		Started:          true,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	test_url, _ = mongoClient.Endpoint(context.Background(), "")
	exitVal := m.Run()
	testcontainers.TerminateContainer(mongoClient)
	os.Exit(exitVal)
}

func TestRegisterModel(t *testing.T) {

	mgd, err := New("mongodb://"+test_url, "test-db")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	type UserModel struct {
		ID    string `db:"mongoid"`
		Email string `db:"email,required,index,unique"`
		Name  string `db:"name,required"`
	}
	type args struct {
		conn  *mongoDatabase
		coll  string
		model UserModel
	}
	tests := []struct {
		name    string
		args    args
		want    database.Model[UserModel]
		wantErr bool
	}{
		{
			name: "Test Register User Model",
			args: args{
				conn:  mgd,
				coll:  "users",
				model: UserModel{},
			},
			want:    &MongoModel[UserModel]{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RegisterModel(tt.args.conn, tt.args.coll, tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.IsType(t, got, tt.want)
		})
	}

}

func TestRegisterUUIDModel(t *testing.T) {

	mgd, err := New("mongodb://"+test_url, "test-db")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	type UserModel struct {
		ID    uuid.UUID `db:"id,index,unique"`
		Email string    `db:"email,required,index,unique"`
		Name  string    `db:"name,required"`
	}
	type args struct {
		conn  *mongoDatabase
		coll  string
		model UserModel
	}
	tests := []struct {
		name    string
		args    args
		want    database.Model[UserModel]
		wantErr bool
	}{
		{
			name: "Test Register User Two Model",
			args: args{
				conn:  mgd,
				coll:  "users-two",
				model: UserModel{},
			},
			want:    &MongoModel[UserModel]{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RegisterModel(tt.args.conn, tt.args.coll, tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.IsType(t, got, tt.want)
		})
	}

}
