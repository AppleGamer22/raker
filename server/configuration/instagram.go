package configuration

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

func (server *RakerServer) InstagramPage(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(authenticatedUserKey).(db.User)

	if err := request.ParseForm(); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	history := db.History{
		Type: types.Instagram,
	}

	post := cleaner.Line(request.Form.Get("post"))
	incognito := cleaner.Line(request.Form.Get("incognito")) == "incognito"
	var errs []error

	if post != "" {
		filter := bson.M{
			"post": post,
			"type": types.Instagram,
		}
		if err := server.Histories.FindOne(context.Background(), filter).Decode(&history); err != nil {
			instagram := shared.NewInstagram(user.Instagram.FBSR, user.Instagram.SessionID, user.Instagram.UserID)
			var (
				username string
				URLs     []string
			)
			log.Debug(request.Form.Get("incognito"))
			if incognito {
				URLs, username, _, err = shared.InstagramIncognito(post)
			} else {
				URLs, username, err = instagram.Post(post)
			}
			if err != nil {
				log.Error(err)
				writer.WriteHeader(http.StatusBadRequest)
				historyHTML(user, history, []error{err}, writer)
				return
			}

			localURLs := make([]string, 0, len(URLs))
			errs = make([]error, 0, len(URLs))
			for _, urlString := range URLs {
				URL, err := url.Parse(urlString)
				if err != nil {
					log.Error(err)
					writer.WriteHeader(http.StatusBadRequest)
					errs = append(errs, err)
					continue
				}
				fileName := fmt.Sprintf("%s_%s", post, path.Base(URL.Path))
				localURLs = append(localURLs, fileName)
			}

			localURLs, saveErrors := StorageHandler.SaveBundle(user, types.Instagram, username, localURLs, URLs, []*http.Cookie{})
			errs = append(errs, saveErrors...)
			for _, err := range saveErrors {
				log.Error(err)
				writer.WriteHeader(http.StatusInternalServerError)
			}

			if len(localURLs) > 0 {
				history = db.History{
					ID:        primitive.NewObjectID().Hex(),
					U_ID:      user.ID.Hex(),
					URLs:      localURLs,
					Type:      types.Instagram,
					Owner:     username,
					Post:      post,
					Incognito: incognito,
					Date:      time.Now(),
				}

				if _, err := server.Histories.InsertOne(context.Background(), history); err != nil {
					log.Error(err)
					writer.WriteHeader(http.StatusInternalServerError)
					errs = append(errs, err)
				}
			}
		}
	}

	historyHTML(user, history, errs, writer)
}
