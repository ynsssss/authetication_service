package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func GetConnection(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetTimeout(time.Second * 5))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

func MustGetConnection(uri string) *mongo.Client {
	client, err := GetConnection(uri)
	if err != nil {
		panic(err)
	}
	return client
}
