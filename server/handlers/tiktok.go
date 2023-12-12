package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AppleGamer22/raker/server/cleaner"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TikTokPage(writer http.ResponseWriter, request *http.Request) {
	user, err := Verify(request)
	if err != nil {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		log.Error(err)
		return
	}

	if err := request.ParseForm(); err != nil {
		log.Error(err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	history := db.History{
		Type: types.TikTok,
	}
	owner := cleaner.Line(request.Form.Get("owner"))
	post := cleaner.Line(request.Form.Get("post"))
	incognito := cleaner.Line(request.Form.Get("incognito")) == "incognito"

	if post != "" {
		filter := bson.M{
			"post": post,
			"type": types.TikTok,
		}
		if err := db.Histories.FindOne(context.Background(), filter).Decode(&history); err != nil {
			tiktok := shared.NewTikTok(user.TikTok.SessionID, user.TikTok.SessionIDGuard, user.TikTok.ChainToken)
			URL, username, cookies, err := tiktok.Post(owner, post, incognito)
			if err != nil {
				log.Error(err)
				writer.WriteHeader(http.StatusBadRequest)
				historyHTML(user, history, []error{err}, writer)
				return
			}

			fileName := fmt.Sprintf("%s.mp4", post)
			if err := StorageHandler.Save(user, types.TikTok, username, fileName, URL, cookies); err != nil {
				log.Error(err)
				writer.WriteHeader(http.StatusInternalServerError)
				historyHTML(user, history, []error{err}, writer)
				return
			}

			history = db.History{
				ID:    primitive.NewObjectID().Hex(),
				U_ID:  user.ID.Hex(),
				URLs:  []string{fileName},
				Type:  types.TikTok,
				Owner: username,
				Post:  post,
				Date:  time.Now(),
			}

			if _, err := db.Histories.InsertOne(context.Background(), history); err != nil {
				log.Error(err)
				writer.WriteHeader(http.StatusInternalServerError)
				historyHTML(user, history, []error{err}, writer)
				return
			}
		}
	}

	historyHTML(user, history, nil, writer)
}
