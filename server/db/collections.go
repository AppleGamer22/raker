package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	// Database  *mongo.Database
	Users     mongo.Collection
	Histories mongo.Collection
)

func Connect(URI, databaseName string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		return nil, err
	}

	if err := Client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	database := Client.Database(databaseName, options.Database())
	// Database = database
	Histories = *database.Collection("histories")
	Users = *database.Collection("users")
	return Client, nil
}
