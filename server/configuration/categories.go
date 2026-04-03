package configuration

import (
	"context"
	"net/http"

	"github.com/AppleGamer22/raker/server/cleaner"
	"github.com/AppleGamer22/raker/server/db"
	"github.com/charmbracelet/log"
)

func (server *RakerServer) Categories(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(authenticatedUserKey).(db.User)

	if err := request.ParseForm(); err != nil {
		http.Error(writer, "failed to read request form", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	for _, category := range user.Categories {
		editedCategory := cleaner.Line(request.Form.Get(category))

		switch editedCategory {
		case "":
			continue
		case http.MethodDelete:
			if err := server.DBClient.UserCategoryRemove(ctx, db.UserCategoryRemoveParams{
				Username: user.Username,
				Category: category,
			}); err != nil {
				log.Error(err, "category", category)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			// Rename: update all user's histories with the old category, then update user's list
			if err := server.DBClient.HistoriesCategoryRename(ctx, db.HistoriesCategoryRenameParams{
				Username:    user.Username,
				OldCategory: category,
				NewCategory: editedCategory,
			}); err != nil {
				log.Error(err, "category", category, "renamed", editedCategory)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := server.DBClient.UserCategoryRemove(ctx, db.UserCategoryRemoveParams{
				Username: user.Username,
				Category: category,
			}); err != nil {
				log.Error(err, "category", category, "removed")
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := server.DBClient.UserCategoryAdd(ctx, db.UserCategoryAddParams{
				Username: user.Username,
				Category: editedCategory,
			}); err != nil {
				log.Error(err, "category", editedCategory, "added")
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	http.Redirect(writer, request, request.Referer(), http.StatusTemporaryRedirect)
}
