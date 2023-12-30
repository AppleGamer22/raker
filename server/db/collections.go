package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
// Client *mongo.Client
// Database  mongo.Database
// Users     mongo.Collection
// Histories mongo.Collection
)

func Connect(URI, databaseName string) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		return nil, nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}

	database := client.Database(databaseName, options.Database())
	// Database = database
	// Client = client
	// Histories = *database.Collection("histories")
	// Users = *database.Collection("users")
	return client, database, nil
}
