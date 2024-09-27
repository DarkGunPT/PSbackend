package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(ctx context.Context, uri string) (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB")

	return client, nil
}

func TestConnection(ctx context.Context, client mongo.Client) error {
	err := client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}
