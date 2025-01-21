package mongodb

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

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

func TestModel(t *testing.T) {
	mgd, err := New("mongodb://"+test_url, "test")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	type UserModel struct {
		ID        uuid.UUID `db:"id,index,unique"`
		Email     string    `db:"email,required,index,unique"`
		Name      string    `db:"name,required"`
		CreatedAt time.Time `db:"created_at"`
		Attempt   int8      `db:"attempt"`
	}

	model, err := RegisterModel(mgd, "users", UserModel{})
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Run("Create User", func(t *testing.T) {
		u := UserModel{
			ID:        uuid.MustParse("e527865d-c83e-4c21-a54b-275f057ecb56"),
			Email:     "jon@doe.com",
			Name:      "Jon Doe",
			CreatedAt: time.Now().UTC(),
			Attempt:   1,
		}
		err := model.WithContext(context.Background()).Save(u)

		require.NoError(t, err)
	})

	t.Run("Find User By Email", func(t *testing.T) {
		u, err := model.WithContext(context.Background()).Query(database.WithFilter("email", "jon@doe.com")).First()
		require.NoError(t, err)
		require.NotEmpty(t, u)
	})

	t.Run("Update By ID", func(t *testing.T) {
		err := model.WithContext(context.Background()).Query(database.WithFilter("id", uuid.MustParse("e527865d-c83e-4c21-a54b-275f057ecb56"))).Update(UserModel{
			Email: "jane@doe.com",
		})
		require.NoError(t, err)
	})
	t.Run("Find User By Updated Email", func(t *testing.T) {
		u, err := model.WithContext(context.Background()).Query(database.WithFilter("email", "jane@doe.com")).First()
		require.NoError(t, err)
		require.NotEmpty(t, u)
	})

	t.Run("Find All Users", func(t *testing.T) {
		u, err := model.WithContext(context.Background()).Query().All()
		require.NoError(t, err)
		require.NotEmpty(t, u)

	})

	t.Run("Delete User By ID", func(t *testing.T) {
		err := model.WithContext(context.Background()).Query(database.WithFilter("id", uuid.MustParse("e527865d-c83e-4c21-a54b-275f057ecb56"))).Delete()
		require.NoError(t, err)
	})
}
