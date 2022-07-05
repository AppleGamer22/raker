package db

import (
	"html/template"
	"net/url"
	"time"

	"github.com/AppleGamer22/rake/server/authenticator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Instagram = "instagram"
	Highlight = "highlight"
	Story     = "story"
	TikTok    = "tiktok"
	VSCO      = "vsco"
)

var MediaTypes = []string{Highlight, Instagram, Story, TikTok, VSCO}

func ValidMediaType(media string) bool {
	return media == Instagram || media == Highlight || media == Story || media == VSCO || media == TikTok
}

func ValidNetworkType(media string) bool {
	return media == Instagram || media == VSCO || media == TikTok
}

func SelectedMediaTypes(mediaTypes []string) map[string]bool {
	result := make(map[string]bool)
	result[Instagram] = true
	result[Highlight] = true
	result[Story] = true
	result[VSCO] = true
	result[TikTok] = true
	if len(mediaTypes) > 0 {
		for _, mediaType := range mediaTypes {
			if _, ok := result[mediaType]; ok && ValidMediaType(mediaType) {
				result[mediaType] = false
			}
		}
		for mediaType, checked := range result {
			result[mediaType] = !checked
		}
	}
	return result
}

var UpdateOption = options.FindOneAndUpdate().SetReturnDocument(options.After)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"-"`
	Username  string             `bson:"username" json:"-"`
	Hash      string             `bson:"hash" json:"-"`
	Instagram struct {
		FBSR      string `bson:"fbsr" json:"-"`
		SessionID string `bson:"session_id" json:"-"`
		UserID    string `bson:"user_id" json:"-"`
		// AppID     string `bson:"app_id" json:"-"`
	} `bson:"instagram" json:"-"`
	TikTok     string    `bson:"tiktok" json:"-"`
	Joined     time.Time `bson:"joined" json:"-"`
	Network    string    `bson:"network" json:"-"`
	Categories []string  `bson:"categories" json:"-"`
}

func (user *User) OpenInstagram(password string) (fbsr, sessionID, userID string, err error) {
	if err := authenticator.Compare(user.Hash, password); err != nil {
		return "", "", "", err
	}
	fbsr, err = authenticator.Decrypt(password, user.Instagram.FBSR)
	if err != nil {
		return "", "", "", err
	}
	sessionID, err = authenticator.Decrypt(password, user.Instagram.SessionID)
	if err != nil {
		return "", "", "", err
	}
	userID, err = authenticator.Decrypt(password, user.Instagram.UserID)
	if err != nil {
		return "", "", "", err
	}
	return fbsr, sessionID, userID, nil
}

func (user *User) OpenTikTok(password string) (string, error) {
	if err := authenticator.Compare(user.Hash, password); err != nil {
		return "", err
	}
	return authenticator.Decrypt(password, user.TikTok)
}

func (user *User) SelectedCategories(categories []string) map[string]bool {
	result := make(map[string]bool)
	for _, category := range user.Categories {
		result[category] = true
	}
	for _, category := range categories {
		if _, ok := result[category]; ok {
			result[category] = false
		}
	}
	for category, checked := range result {
		result[category] = !checked
	}
	return result
}

type History struct {
	ID         string    `bson:"_id" json:"-"`
	U_ID       string    `bson:"U_ID" json:"-"`
	URLs       []string  `bson:"urls" json:"urls"`
	Type       string    `bson:"type" json:"type"`
	Owner      string    `bson:"owner" json:"owner"`
	Post       string    `bson:"post" json:"post"`
	Date       time.Time `bson:"date" json:"date"`
	Categories []string  `bson:"categories" json:"categories"`
}

type HistoryDisplay struct {
	History
	SelectedCategories map[string]bool
	Errors             []error
	Version            string
}

func (historyDisplay HistoryDisplay) HistoryQuery() template.URL {
	query := url.Values{}
	query.Set("page", "1")
	for category := range historyDisplay.SelectedCategories {
		query.Set(category, category)
	}
	return template.URL(query.Encode())
}

type HistoriesDisplay struct {
	Histories  [][]History
	Owner      string
	Types      map[string]bool
	Categories map[string]bool
	Exclusive  bool
	Version    string
	Page       int
	Pages      int
	Count      int
	Error      error
}

func (historiesDisplay HistoriesDisplay) Query(value string) template.URL {
	query := url.Values{}
	if !ValidMediaType(value) {
		query.Set("owner", value)
		for mediaType := range historiesDisplay.Types {
			query.Set(mediaType, mediaType)
		}
	} else {
		query.Set(value, value)
	}
	if historiesDisplay.Exclusive {
		query.Set("exclusive", "exclusive")
	}
	for category, checked := range historiesDisplay.Categories {
		if checked {
			query.Set(category, category)
		}
	}
	return template.URL(query.Encode())
}
