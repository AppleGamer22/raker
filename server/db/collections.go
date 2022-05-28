package db

import "go.mongodb.org/mongo-driver/mongo"

var (
	Users     mongo.Collection
	Histories mongo.Collection
)
