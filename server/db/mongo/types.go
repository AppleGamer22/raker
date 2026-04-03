package db

import (
	"html/template"
	"net/url"
	"time"

	"github.com/AppleGamer22/raker/server/authenticator"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/go-webauthn/webauthn/webauthn"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	UpdateOption = options.FindOneAndUpdate().SetReturnDocument(options.After)
	// writeConcern       = writeconcern.New(writeconcern.W(1))
	// readConcern        = readconcern.Snapshot()
	// TransactionOptions = options.Transaction().SetWriteConcern(writeConcern).SetReadConcern(readConcern)
)

type User struct {
	ID          primitive.ObjectID    `bson:"_id" json:"-"`
	Username    string                `bson:"username" json:"username"`
	Hash        string                `bson:"hash" json:"hash"`
	Credentials []webauthn.Credential `bson:"credentials" json:"-"`
	Session     webauthn.SessionData  `bson:"session" json:"-"`
	Instagram   struct {
		FBSR      string `bson:"fbsr" json:"-"`
		SessionID string `bson:"session_id" json:"session_id"`
		UserID    string `bson:"user_id" json:"user_id"`
		// AppID     string `bson:"app_id" json:"-"`
	} `bson:"instagram" json:"Instagram"`
	TikTok struct {
		SessionID      string `bson:"session_id" json:"-"`
		SessionIDGuard string `bson:"session_id_guard" json:"-"`
		ChainToken     string `bson:"chain_token" json:"-"`
	} `bson:"tiktok" json:"-"`
	Joined     time.Time `bson:"joined" json:"-"`
	Network    string    `bson:"network" json:"network"`
	Categories []string  `bson:"categories" json:"categories"`
}

type UserCategoryDisplay struct {
	Username     string
	Categories   []string
	HistoryQuery template.URL
	Version      string
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
	return authenticator.Decrypt(password, user.TikTok.SessionID)
}

func (user User) WebAuthnID() []byte {
	return []byte(user.Username)
}

func (user User) WebAuthnName() string {
	return user.Username
}

func (user User) WebAuthnDisplayName() string {
	return user.Username
}

func (user User) WebAuthnCredentials() []webauthn.Credential {
	return user.Credentials
}

func (user User) WebAuthnIcon() string {
	return ""
}

type HistoryArchive struct {
	ID    string   `bson:"_id" json:"-"`
	U_ID  string   `bson:"U_ID" json:"-"`
	URLs  []string `bson:"urls" json:"urls"`
	Type  string   `bson:"type" json:"type"`
	Owner string   `bson:"owner" json:"owner"`
	Post  string   `bson:"post" json:"post"`
	Date  struct {
		Value time.Time `json:"$date"`
	} `bson:"date" json:"date"`
	Incognito  bool     `bson:"incognito" json:"incognito"`
	Categories []string `bson:"categories" json:"categories"`
}

type HistoryDisplay struct {
	db.History
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
	Histories  [][]db.History
	PostOwner  string
	Types      map[string]bool
	Categories map[string]bool
	Exclusive  bool
	Version    string
	Page       int
	Pages      int
	Count      int
	Error      error
}

func (historiesDisplay HistoriesDisplay) Query(value interface{}) template.URL {
	queryValue := ""
	switch value := value.(type) {
	case string:
		queryValue = value
	case db.PostType:
		queryValue = string(value)
	}

	query := url.Values{}
	if !db.PostType(queryValue).Valid() {
		query.Set("owner", queryValue)
		for mediaType := range historiesDisplay.Types {
			query.Set(mediaType, mediaType)
		}
	} else {
		query.Set(queryValue, queryValue)
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
