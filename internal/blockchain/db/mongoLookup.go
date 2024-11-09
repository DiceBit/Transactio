package db

import (
	"Transactio/internal/blockchain/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserInfo struct {
	Username string         `bson:"username"`
	Info     map[string]int `bson:"info"`
}

/*
"username": {
	"filename": index,
    "filename2": index,
...
}
*/

func InsertInfo(ctx context.Context, client *mongo.Client, username, filename string, index int) error {
	collections := client.Database(utils.MongoDbName).Collection(utils.MongoCollections)

	filter := bson.M{"username": username}
	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("info.%s", filename): index,
		},
	}

	_, err := collections.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	//("Added/Updated file '%s' for user '%s' with index %d\n", filename, username, index)
	return nil
}
func RemoveInfo(ctx context.Context, client *mongo.Client, username, filename string) error {
	collections := client.Database(utils.MongoDbName).Collection(utils.MongoCollections)

	filter := bson.M{"username": username}
	update := bson.M{
		"$unset": bson.M{
			fmt.Sprintf("info.%s", filename): "",
		},
	}

	_, err := collections.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func GetInfo(ctx context.Context, client *mongo.Client) ([]UserInfo, error) {
	collections := client.Database(utils.MongoDbName).Collection(utils.MongoCollections)
	filter := bson.D{}

	cursor, err := collections.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var data []UserInfo
	for cursor.Next(ctx) {
		var l UserInfo
		err := cursor.Decode(&l)
		if err != nil {
			return nil, err
		}
		data = append(data, l)
	}

	return data, cursor.Err()
}

// CreateIndex used only when create new table in db
func CreateIndex(ctx context.Context, client *mongo.Client) error {
	collections := client.Database(utils.MongoDbName).Collection(utils.MongoCollections)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := collections.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	return nil
}

// ?
func RemoveUser(ctx context.Context, client *mongo.Client, username string) error {
	collections := client.Database(utils.MongoDbName).Collection(utils.MongoCollections)

	filter := bson.M{"username": username}

	_, err := collections.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
