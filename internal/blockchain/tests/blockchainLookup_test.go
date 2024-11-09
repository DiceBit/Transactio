package tests

import (
	"Transactio/internal/blockchain/db"
	mongodb "Transactio/pkg/dbConn/mongo"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"testing"
)

type testObj struct {
	db *mongo.Client
}

var obj *testObj

func TestMain(m *testing.M) {
	tdb, _ := mongodb.New()
	tobj := &testObj{
		db: tdb,
	}
	err := db.CreateIndex(context.Background(), tobj.db)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Info: Index Created")
	obj = tobj

	m.Run()
}
func printRes(client *mongo.Client) {
	val, _ := db.GetInfo(context.Background(), client)
	for _, v := range val {
		fmt.Printf("\n%v {\n \t %v \n}", v.Username, v.Info)
	}
}

func TestInsertInfo(t *testing.T) {
	type args struct {
		client   *mongo.Client
		username string
		filename string
		index    int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "#1",
			args: args{
				client:   obj.db,
				username: "testUser1",
				filename: "testFile1",
				index:    12,
			},
			wantErr: false,
		},
		{
			name: "#2",
			args: args{
				client:   obj.db,
				username: "testUser1",
				filename: "testFile2",
				index:    13,
			},
			wantErr: false,
		},
		{
			name: "#3",
			args: args{
				client:   obj.db,
				username: "testUser2",
				filename: "testFile1",
				index:    12,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.InsertInfo(context.Background(), tt.args.client, tt.args.username, tt.args.filename, tt.args.index); (err != nil) != tt.wantErr {
				t.Errorf("InsertInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		printRes(tt.args.client)

	}

}

func TestRemoveInfo(t *testing.T) {
	type args struct {
		client   *mongo.Client
		username string
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "#1",
			args: args{
				client:   obj.db,
				username: "testUser1",
				filename: "testFile2",
			},
			wantErr: false,
		},
		{
			name: "#2",
			args: args{
				client:   obj.db,
				username: "testUser2",
				filename: "testFile10",
			},
			wantErr: false,
		},
		{
			name: "#3",
			args: args{
				client:   obj.db,
				username: "testUser2",
				filename: "testFile1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.RemoveInfo(context.Background(), tt.args.client, tt.args.username, tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("RemoveInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		printRes(tt.args.client)
	}
}

func TestRemoveUser(t *testing.T) {
	type args struct {
		client   *mongo.Client
		username string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "#1",
			args: args{
				client:   obj.db,
				username: "testUser2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.RemoveUser(context.Background(), tt.args.client, tt.args.username); (err != nil) != tt.wantErr {
				t.Errorf("RemoveUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		printRes(tt.args.client)
	}
}
