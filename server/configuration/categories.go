package configuration

import (
	"context"
	"net/http"

	"github.com/AppleGamer22/raker/server/cleaner"
	db "github.com/AppleGamer22/raker/server/db/mongo"
	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (server *RakerServer) Categories(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(authenticatedUserKey).(db.User)

	if err := request.ParseForm(); err != nil {
		http.Error(writer, "failed to read request form", http.StatusBadRequest)
		return
	}

	for _, category := range user.Categories {
		editedCategory := cleaner.Line(request.Form.Get(category))
		filter := bson.M{
			"$or": bson.A{
				bson.M{
					"_id": user.ID,
				},
				bson.M{
					"U_ID": user.ID.Hex(),
				},
			},
		}
		operations := []mongo.WriteModel{}
		switch editedCategory {
		case "":
			continue
		case http.MethodDelete:
			updateOperation := mongo.NewUpdateManyModel()
			filter["categories"] = category
			updateOperation.SetFilter(filter)
			updateOperation.SetUpdate(bson.M{
				"$pull": bson.M{
					"categories": category,
				},
			})
			operations = append(operations, updateOperation)
		default:
			updateOperation := mongo.NewUpdateOneModel()
			filter["categories"] = category
			updateOperation.SetFilter(filter)
			updateOperation.SetUpdate(bson.M{
				"$set": bson.M{
					"categories.$": editedCategory,
				},
			})
			operations = append(operations, updateOperation)

			sortOperation := mongo.NewUpdateManyModel()
			filter["categories"] = editedCategory
			sortOperation.SetFilter(filter)
			sortOperation.SetUpdate(bson.M{
				"$push": bson.M{
					"$each": bson.A{},
					"$sort": 1,
				},
			})
			operations = append(operations, sortOperation)
		}

		bulkOptions := options.BulkWriteOptions{}
		bulkOptions.SetOrdered(true)

		if _, err := server.Histories.BulkWrite(context.Background(), operations, &bulkOptions); err != nil {
			log.Error(err, "category", category, "edited", editedCategory)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := server.Users.BulkWrite(context.Background(), operations, &bulkOptions); err != nil {
			log.Error(err, "category", category, "edited", editedCategory)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(writer, request, request.Referer(), http.StatusTemporaryRedirect)
}
