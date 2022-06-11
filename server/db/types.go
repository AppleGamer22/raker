package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Instagram = "instagram"
	Highlight = "highlight"
	Story     = "story"
	VSCO      = "vsco"
	TikTok    = "tiktok"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id" json:"-"`
	Username   string             `bson:"username" json:"-"`
	Hash       string             `bson:"hash" json:"-"`
	Key        string             `bson:"key" json:"-"`
	Joined     time.Time          `bson:"joined" json:"-"`
	Network    string             `bson:"network" json:"-"`
	Instagram  bool               `bson:"instagram" json:"-"`
	Categories []string           `bson:"categories" json:"-"`
}

type History struct {
	ID         string    `bson:"_id" json:"-"`
	URLs       []string  `bson:"urls" json:"urls"`
	U_ID       string    `bson:"U_ID" json:"-"`
	Type       string    `bson:"type" json:"type"`
	Owner      string    `bson:"owner" json:"owner"`
	Post       string    `bson:"post" json:"post"`
	Date       time.Time `bson:"date" json:"date"`
	Categories []string  `bson:"categories" json:"categories"`
}

func ValidMediaType(media string) bool {
	return media == Instagram || media == Highlight || media == Story || media == VSCO || media == TikTok
}

func ValidNetworkType(media string) bool {
	return media == Instagram || media == VSCO || media == TikTok
}
