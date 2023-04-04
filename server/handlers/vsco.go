package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/AppleGamer22/raker/server/cleaner"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/AppleGamer22/raker/shared"
	"github.com/AppleGamer22/raker/shared/types"
	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func VSCOPage(writer http.ResponseWriter, request *http.Request) {
	user, err := Verify(request)
	if err != nil {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		log.Error(err)
		return
	}

	if err := request.ParseForm(); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	history := db.History{
		Type: types.VSCO,
	}
	owner := cleaner.Line(request.Form.Get("owner"))
	post := cleaner.Line(request.Form.Get("post"))
	if post != "" {
		filter := bson.M{
			"post": post,
			"type": types.VSCO,
		}
		if err := db.Histories.FindOne(context.Background(), filter).Decode(&history); err != nil {
			urlString, username, err := shared.VSCO(owner, post)
			if err != nil {
				log.Error(err)
				writer.WriteHeader(http.StatusBadRequest)
				historyHTML(user, history, []error{err}, writer)
				return
			}

			URL, err := url.Parse(urlString)
			if err != nil {
				log.Error(err)
				writer.WriteHeader(http.StatusBadRequest)
				historyHTML(user, history, []error{err}, writer)
				return
			}
			fileName := fmt.Sprintf("%s_%s", post, path.Base(URL.Path))

			if err := StorageHandler.Save(user, types.VSCO, username, fileName, urlString); err != nil {
				log.Error(err)
				writer.WriteHeader(http.StatusInternalServerError)
				historyHTML(user, history, []error{err}, writer)
				return
			}

			history = db.History{
				ID:    primitive.NewObjectID().Hex(),
				U_ID:  user.ID.Hex(),
				URLs:  []string{fileName},
				Type:  types.VSCO,
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
